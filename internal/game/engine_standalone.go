//go:build standalone

package game

import (
	"fmt"
	"log"

	"github.com/bjatkin/flappy_boot/internal/emu/ppu"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/key"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	"github.com/bjatkin/flappy_boot/internal/hardware/save"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const saveFile = "flappy_boot_stand.sav"

// Harness is the standalone harness that allows the emulator to be used in standalone mode
type Harness struct {
	E        *Engine
	R        Runable
	PPU      *ppu.PPU
	saveData [save.DataLen]byte
	frame    int
}

// NewHarness creates a new engine harness
func NewHarness() *Harness {
	save.LoadData(saveFile)

	harness := &Harness{
		E:   NewEngine(),
		PPU: ppu.New(),
	}

	harness.PPU.Backgrounds[0].SkipGFXUpdate = true
	harness.PPU.Backgrounds[1].SkipGFXUpdate = true

	return harness
}

// Update runs the GBA update/ draw code at 60TPS
func (h *Harness) Update() error {
	h.frame++
	h.E.Draw()

	// update the input register (note we only need to update keys that the game actually uses)
	keyReg := memmap.Input(0xFFFF)
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		keyReg &= ^key.StartMask
	}
	if ebiten.IsKeyPressed(ebiten.KeyC) {
		keyReg &= ^key.AMask
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		keyReg &= ^key.UpMask
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		keyReg &= ^key.DownMask
	}
	*key.Input = keyReg

	h.E.Update(h.R)
	if h.frame%10 == 0 {
		// only check the save buffer every 10 frame to help improve performance
		h.updateSaveData(saveFile)
	}
	return nil
}

// updateSaveData updates the save data in the save file
func (h *Harness) updateSaveData(path string) {
	var delta bool
	for i := 0; i < len(h.saveData); i++ {
		if h.saveData[i] != byte(save.SRAM[i]) {
			delta = true
		}
		// update save data so we can tell if something changes
		h.saveData[i] = byte(save.SRAM[i])
	}

	if !delta {
		// save data has not changed, no reason to do an update
		return
	}

	save.SaveData(saveFile, h.saveData[:])

	return
}

// Draw takes the data form inside the simulated GBA memory and draws it onto the screen
// this can happy more than 60 times a second which is why the actuall GBA draw call needs to be in the Update function
func (h *Harness) Draw(screen *ebiten.Image) {
	h.PPU.Update()

	// i is the priority, we draw from lowest prioirity(3) to highest priority(0) so higher priorities
	// overwrite pixels from lower priorities
	for i := 3; i >= 0; i-- {

		screen.DrawImage(ebiten.NewImageFromImage(h.PPU.Final), &ebiten.DrawImageOptions{})

		for _, spr := range h.PPU.Sprites {
			if !spr.Enabled {
				continue
			}
			if spr.Priority != i {
				continue
			}

			y := ((spr.Pos.Y + spr.Size.Y) % 0xFF) - spr.Size.Y
			x := ((spr.Pos.X + spr.Size.X) % 0x1FF) - spr.Size.X
			if x > 240 || y > 160 {
				continue
			}

			transform := ebiten.GeoM{}
			switch {
			case spr.VFlip && spr.HFlip:
				transform.Scale(-1, -1)
				// Transform must take scale into account since all the sprites are 64x64 by default
				transform.Translate(float64(x+spr.Size.X), float64((y + spr.Size.Y)))
			case spr.VFlip:
				transform.Scale(1, -1)
				// Transform must take scale into account since all the sprites are 64x64 by default
				transform.Translate(float64(x), float64((y + spr.Size.Y)))
			case spr.HFlip:
				transform.Scale(-1, 1)
				// Transform must take scale into account since all the sprites are 64x64 by default
				transform.Translate(float64(x+spr.Size.X), float64(y))
			default:
				transform.Translate(float64(x), float64(y))
			}
			screen.DrawImage(ebiten.NewImageFromImage(spr.Image), &ebiten.DrawImageOptions{GeoM: transform})
		}
	}

	// performance
	msg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

// Layout just returns the resolution of the GBA
func (h *Harness) Layout(outsideWidth, outsizeHeight int) (int, int) {
	return display.Width, display.Height
}

// Run init's the game and starts running it
func (h *Harness) Run(run Runable) {
	ebiten.SetWindowSize(display.Width*4, display.Height*4)
	ebiten.SetWindowTitle("Flappy Boot Advance")
	ebiten.SetTPS(60) // match the refresh rate of the GBA (more or less)

	h.E.Init(run)
	h.R = run

	if err := ebiten.RunGame(h); err != nil {
		log.Fatal(err)
	}
}
