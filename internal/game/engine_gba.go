//go:build !standalone

package game

// Harness is the standalone harness that allows the emulator to be used in standalone mode
type Harness struct {
	E *Engine
}

// NewHarness creates a new engine harness
func NewHarness() *Harness {
	return &Harness{
		E: NewEngine(),
	}
}

// Run init's the game and starts running it
func (h *Harness) Run(run Runable) error {
	h.E.Init(run)

	for {
		h.E.Update(run)

		vSyncWait()

		h.E.Draw()
	}
}
