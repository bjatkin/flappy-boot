package pillar

import (
	"math/rand"

	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/math"
)

// BG is a background the scrolls and spawns in pillars for the player to dodge
type BG struct {
	bg          *game.Background
	rand        *rand.Rand
	nextPillar  int
	pillarEvery int
	gapSize     int
	lastPoint   int
	meta        meta
	scrollSpeed math.Fix8

	started bool
}

// NewBG creates a new BG struct
func NewBG(pillarEvery int, bg *game.Background) *BG {
	pillars := &BG{
		bg:          bg,
		gapSize:     7,
		pillarEvery: pillarEvery,
		meta:        meta{},
		scrollSpeed: math.NewFix8(1, 32),
	}

	return pillars
}

// Init sets the background to it's initial state, resetting it's horizontal scroll and removing all active pillars
func (p *BG) Init() {
	p.started = false
	p.bg.HScroll = 0
	p.rand = nil

	for i := range p.meta.pillars {
		p.deletePillar(i)
	}
}

// CheckPoint checks to see if the current math.Rect has passed through a new pillar gap
func (p *BG) CheckPoint(check math.Rect) bool {
	buffer := 4
	for i := range p.meta.pillars {
		if !p.meta.IsSet(i) {
			continue
		}
		pillar := p.meta.pillars[i]

		if pillar.X2 == p.lastPoint {
			continue
		}

		right := pillar.X2 - p.bg.HScroll.Int()
		if right < check.X1+buffer {
			p.lastPoint = pillar.X2
			return true
		}
	}
	return false
}

// Start indicates the that game as started and the background should start spawing pillars
func (p *BG) Start() {
	if p.started {
		return
	}

	p.started = true
	p.lastPoint = p.bg.HScroll.Int() + p.pillarEvery
}

// CollisionCheck checks the provided rect against all the pillars in the current background. If the rect collides
// with any pillar CollisionCheck returns true
func (p *BG) CollisionCheck(check math.Rect) bool {
	for i := range p.meta.pillars {
		if !p.meta.set[i] {
			continue
		}

		left := p.meta.pillars[i].X1 - p.bg.HScroll.Int()
		if left <= 0 {
			continue
		}
		right := p.meta.pillars[i].X2 - p.bg.HScroll.Int()
		top := p.meta.pillars[i].Y1
		bottom := p.meta.pillars[i].Y2

		if check.X2 < left || check.X1 > right {
			continue
		}
		if check.Y1 < top || check.Y2 > bottom {
			return true
		}
	}

	return false
}

// addPillar adds a new pillar to the background
func (p *BG) addPillar(x int) {
	startTile := (x % 512) / 8
	columns := [4]int{startTile, (startTile + 1) % 64, (startTile + 2) % 64, (startTile + 3) % 64}

	gap := p.rand.Intn(15 - p.gapSize)
	for i := 0; i < 18; i++ {
		tiles := [4]int{}
		switch {
		case i == gap:
			tiles = [4]int{13, 22, 11, 10}
		case i == gap+p.gapSize:
			tiles = [4]int{24, 29, 20, 21}
		case i > gap && i < gap+p.gapSize:
			continue
		case i == 16:
			tiles = [4]int{2, 4, 8, 1}
		case i == 17:
			tiles = [4]int{16, 25, 23, 17}
		default:
			tiles = [4]int{14, 30, 28, 15}
		}

		for j := range tiles {
			p.bg.SetTile(columns[j], i, tiles[j])
		}
	}

	hBoxX := (x / 8) * 8
	p.meta.Append(hBoxX+3, gap*8+4, hBoxX+27, (gap+p.gapSize)*8+2)
}

// deletePillar removes the pillar located at math.Rect r
func (p *BG) deletePillar(i int) {
	if !p.meta.IsSet(i) {
		return
	}

	start := (p.meta.pillars[i].X1 % 512) / 8
	columns := [4]int{start, (start + 1) % 64, (start + 2) % 64, (start + 3) % 64}

	for i := 0; i < 18; i++ {
		for j := 0; j < 4; j++ {
			p.bg.SetTile(columns[j], i, 0)
		}
	}

	p.meta.Delete(i)
}

// Update updates the background including scrolling, adding new pillars, and removing old pillars
func (p *BG) Update() {
	p.bg.HScroll += p.scrollSpeed

	if !p.started {
		return
	}
	if p.rand == nil {
		p.rand = rand.New(rand.NewSource(int64(p.bg.HScroll)))

	}

	// add pillars to the right just off screen
	p.nextPillar--
	if p.nextPillar <= 0 {
		x := p.bg.HScroll.Int() + 256
		p.addPillar(x)
		p.nextPillar = p.pillarEvery
	}

	// remove current pillars to the left of the screen
	border := p.bg.HScroll.Int() - 32
	for i := range p.meta.pillars {
		if !p.meta.IsSet(i) {
			continue
		}

		if p.meta.pillars[i].X1 < border {
			p.deletePillar(i)
		}
	}
}

// Show adds the background to the list of active backgrounds
func (p *BG) Show() error {
	return p.bg.Show()
}

// Hide hides the current background
func (p *BG) Hide() {
	p.bg.Hide()
}

// meta hold meta data about the pillars in the background. It can hold metadata for up to 10 pillars at a time
type meta struct {
	pillars [10]math.Rect
	set     [10]bool
	i       int
}

// IsSet returns true only if the pillar at index 'i' is an 'active' onscreen pillar
func (m *meta) IsSet(i int) bool {
	return m.set[i]
}

// Append adds a new pillar to the current list of pillars. meta is a circular buffer so after appending 10 pillars
// Append will wrap around and start adding pillars to the front of the buffer again
func (m *meta) Append(x1, y1, x2, y2 int) {
	m.set[m.i] = true
	m.pillars[m.i].X1 = x1
	m.pillars[m.i].Y1 = y1
	m.pillars[m.i].X2 = x2
	m.pillars[m.i].Y2 = y2
	m.i++
	m.i %= len(m.pillars)
}

// Delete removes a pillar from the list by marking it as unset
func (m *meta) Delete(i int) {
	m.set[i] = false
}
