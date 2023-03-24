package fly

import (
	"math/rand"

	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/fix"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/key"
)

type Stage struct {
	sky         *game.Background
	pillarBG    *pillarBG
	player      *player
	gravity     fix.P8
	ground      fix.P8
	scrollSpeed fix.P8
}

func NewStage() *Stage {
	return &Stage{
		gravity:     fix.Quarter,
		ground:      fix.One * 131,
		scrollSpeed: fix.One + fix.Eighth,
	}
}

func (s *Stage) Init(e *game.Engine) error {
	s.pillarBG = newPillarBG(e.NewBackground(assets.PillarsTileMap, display.Priority2),
		8, 4, 5,
	)
	err := s.pillarBG.Show()
	if err != nil {
		return err
	}

	s.sky = e.NewBackground(assets.SkyTileMap, display.Priority3)
	err = s.sky.Add()
	if err != nil {
		return err
	}

	s.player = newPlayer(fix.One*40, fix.One*10, e.NewSprite(assets.PlayerTileSet))
	err = s.player.Show()
	if err != nil {
		return err
	}

	return nil
}

func (s *Stage) Update(e *game.Engine, frame int) error {
	var jump fix.P8
	if key.JustPressed(key.A) {
		jump = -fix.One * 3
	}

	s.player.Update(s.gravity, jump, s.ground)
	s.pillarBG.Update(s.scrollSpeed)

	return nil
}

func (t *Stage) Next() (game.Runable, bool) {
	return nil, false
}

type player struct {
	sprite *game.Sprite
	dy     fix.P8
	maxDy  fix.P8
}

func newPlayer(x, y fix.P8, sprite *game.Sprite) *player {
	p := &player{
		sprite: sprite,
		maxDy:  fix.One * 5,
	}

	p.sprite.X = x
	p.sprite.Y = y
	return p
}

func (p *player) Show() error {
	err := p.sprite.Add()
	if err != nil {
		return err
	}

	return nil
}

func (p *player) Hide() {
	p.sprite.Remove()
}

func (p *player) Update(gravity, jump, ground fix.P8) {
	p.dy += gravity
	if p.dy > p.maxDy {
		p.dy = p.maxDy
	}

	if jump != 0 {
		p.dy = jump
	}

	p.sprite.Y += p.dy
	if p.sprite.Y > ground {
		p.sprite.Y = ground
	}

	if p.sprite.Y < 0 {
		p.sprite.Y = 0
		p.dy = 0
	}
}

type pillarBG struct {
	bg      *game.Background
	scrollX fix.P8
}

func newPillarBG(bg *game.Background, pillarSpace, pillarShift, gapSize int) *pillarBG {
	pillars := &pillarBG{
		bg: bg,
	}

	// create all the initial pillars
	for i := 0; i < 8; i++ {
		pillarX := (i * pillarSpace) + pillarShift
		gap := rand.Intn(15 - gapSize)
		for j := 0; j < 16; j++ {
			switch {
			case j == gap:
				pillars.bg.SetTile(pillarX, j, 13)
				pillars.bg.SetTile(pillarX+1, j, 22)
				pillars.bg.SetTile(pillarX+2, j, 11)
				pillars.bg.SetTile(pillarX+3, j, 10)
			case j == gap+gapSize:
				pillars.bg.SetTile(pillarX, j, 24)
				pillars.bg.SetTile(pillarX+1, j, 29)
				pillars.bg.SetTile(pillarX+2, j, 20)
				pillars.bg.SetTile(pillarX+3, j, 21)
			case j > gap && j < gap+gapSize:
				continue
			default:
				pillars.bg.SetTile(pillarX, j, 14)
				pillars.bg.SetTile(pillarX+1, j, 30)
				pillars.bg.SetTile(pillarX+2, j, 28)
				pillars.bg.SetTile(pillarX+3, j, 15)
			}

		}
	}

	return pillars
}

func (p *pillarBG) Update(scrollSpeed fix.P8) {
	p.scrollX += scrollSpeed
	p.bg.SetScroll(p.scrollX.Int(), 0)
}

func (p *pillarBG) Show() error {
	err := p.bg.Add()
	if err != nil {
		return err
	}
	return nil
}

func (p *pillarBG) Hide() {
	p.bg.Remove()
}
