package state

import "github.com/bjatkin/flappy_boot/internal/math"

const (
	A = State(1)
	B = State(1 << 1)
	C = State(1 << 2)
	D = State(1 << 3)
	E = State(1 << 4)
	F = State(1 << 5)
)

// State represents one of the valid Tracker states
type State int

// Tracker tracks high level state in a scene
type Tracker struct {
	SceneFrames map[State]int
	state       State
	frame       int
}

// Init initializes the basic values of the Tracker, it is safe to call multiple times
func (t *Tracker) Init() {
	t.state = A
	t.frame = 0
}

// Update should be called each frame as it updates the Trackers internal state
func (t *Tracker) Update() {
	t.frame++
	sceneFrames := t.SceneFrames[t.state]
	if sceneFrames == 0 {
		return
	}

	if t.frame > sceneFrames {
		t.state <<= 1
		t.frame = 0
	}
}

// Current returns the current state of the Tracker
func (t *Tracker) Current() State {
	return t.state
}

// Frame returns the current frame of the specific state
func (t *Tracker) Frame() int {
	return t.frame
}

// Frac returns the percentage the tracke is into it's current state. It's a number between 0 and 1
func (t *Tracker) Frac() math.Fix8 {
	sceneFrames := t.SceneFrames[t.state]
	if sceneFrames == 0 {
		return 0
	}

	return math.Fix8((t.frame << 8) / sceneFrames)
}

// Next moves the tracker from it's current state into the next state
func (t *Tracker) Next() {
	t.state <<= 1
}

// Is returns true the provided state reflects the current state of the of the tracker.
// states can be compositis of multiple states (e.g. A | D)
func (t *Tracker) Is(state State) bool {
	return t.state&state > 0
}
