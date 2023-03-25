package fly

import (
	"errors"
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
	skyScroll   fix.P8
}

func NewStage() *Stage {
	return &Stage{
		gravity:     fix.Quarter,
		ground:      fix.One * 131,
		scrollSpeed: fix.One + fix.Eighth,
	}
}

func (s *Stage) Init(e *game.Engine) error {
	s.pillarBG = newPillarBG(100, e.NewBackground(assets.PillarsTileMap, display.Priority2))
	err := s.pillarBG.Show()
	if err != nil {
		return err
	}

	s.sky = e.NewBackground(assets.SkyTileMap, display.Priority3)
	err = s.sky.Add()
	if err != nil {
		return err
	}

	s.player = newPlayer(fix.One*40, fix.One*62, e.NewSprite(assets.PlayerTileSet))
	err = s.player.Show()
	if err != nil {
		return err
	}

	return nil
}

func (s *Stage) Update(e *game.Engine, frame int) error {
	var jump fix.P8
	if key.JustPressed(key.A) {
		s.player.started = true
		s.pillarBG.started = true
		jump = -fix.One * 3
	}

	s.player.Update(s.gravity, jump, s.ground)

	s.skyScroll += s.scrollSpeed / 2
	s.sky.SetScroll(s.skyScroll.Int(), 0)
	err := s.sky.Add()
	if err != nil {
		return err
	}

	s.pillarBG.Update(s.scrollSpeed)
	err = s.pillarBG.Show()
	if err != nil {
		return err
	}

	if s.pillarBG.collisionCheck(s.player.Rect()) {
		return errors.New("game over")
	}

	return nil
}

func (t *Stage) Next() (game.Runable, bool) {
	return nil, false
}

type player struct {
	sprite *game.Sprite
	dy     fix.P8
	maxDy  fix.P8

	started bool
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

func (p *player) Rect() rect {
	return rect{
		x1: p.sprite.X.Int() + 2,
		y1: p.sprite.Y.Int() + 2,
		x2: p.sprite.X.Int() + 12,
		y2: p.sprite.Y.Int() + 12,
	}
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
	if !p.started {
		// don't update physics if the game has not started yet
		return
	}

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

type rect struct {
	x1, y1 int
	x2, y2 int
}

type pillarBG struct {
	bg      *game.Background
	scrollX fix.P8

	rand        *rand.Rand
	nextPillar  int
	pillarEvery int
	pillars     []rect
	gapSize     int

	started bool
}

func newPillarBG(pillarEvery int, bg *game.Background) *pillarBG {
	pillars := &pillarBG{
		bg:          bg,
		gapSize:     6,
		pillarEvery: pillarEvery,
	}

	return pillars
}

func (p *pillarBG) collisionCheck(check rect) bool {
	for i := range p.pillars {
		left := p.pillars[i].x1 - p.scrollX.Int()
		if left <= 0 {
			continue
		}
		right := left + 32
		top := p.pillars[i].y1
		bottom := p.pillars[i].y2

		if check.x2 < left || check.x1 > right {
			continue
		}
		if check.y1 < top || check.y2 > bottom {
			return true
		}
	}

	return false
}

func (p *pillarBG) addPillar(x, gapSize int) rect {
	start := (x % 512) / 8
	columns := [4]int{start, (start + 1) % 64, (start + 2) % 64, (start + 3) % 64}

	gap := p.rand.Intn(15 - gapSize)
	for i := 0; i < 18; i++ {
		tiles := [4]int{}
		switch {
		case i == gap:
			tiles = [4]int{13, 22, 11, 10}
		case i == gap+gapSize:
			tiles = [4]int{24, 29, 20, 21}
		case i > gap && i < gap+gapSize:
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

	return rect{x1: x, y1: gap*8 + 4, x2: x + 32, y2: (gap+gapSize)*8 + 4}
}

func (p *pillarBG) removePillar(r rect) {
	start := (r.x1 % 512) / 8
	columns := [4]int{start, (start + 1) % 64, (start + 2) % 64, (start + 3) % 64}

	for i := 0; i < 18; i++ {
		for j := 0; j < 4; j++ {
			p.bg.SetTile(columns[j], i, 0)
		}
	}
}

func (p *pillarBG) Update(scrollSpeed fix.P8) {
	p.scrollX += scrollSpeed
	p.bg.SetScroll(p.scrollX.Int(), 0)

	if !p.started {
		return
	}
	if p.rand == nil {
		p.rand = rand.New(rand.NewSource(int64(p.scrollX)))

	}

	// add pillars to the right just off screen
	p.nextPillar--
	var keep []rect
	if p.nextPillar <= 0 {
		x := p.scrollX.Int() + 256
		r := p.addPillar(x, 5)
		p.nextPillar = p.pillarEvery
		keep = []rect{r}
	}

	// remove current pillars to the left of the screen
	border := p.scrollX.Int() - 32
	for i := range p.pillars {
		if p.pillars[i].x1 < border {
			p.removePillar(p.pillars[i])
		} else {
			keep = append(keep, p.pillars[i])
		}
	}

	p.pillars = keep
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
