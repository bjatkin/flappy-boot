package fly

import (
	"math/rand"

	"github.com/bjatkin/flappy_boot/gameplay/actor"
	"github.com/bjatkin/flappy_boot/gameplay/gameover"
	"github.com/bjatkin/flappy_boot/gameplay/score"
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/key"
	"github.com/bjatkin/flappy_boot/internal/math"
)

type Scene struct {
	sky         *game.Background
	clouds      *game.Background
	pillars     *pillarBG
	player      *actor.Player
	counter     *score.Counter
	scrollSpeed math.Fix8
	gameOver    bool

	gravity    math.Fix8
	ground     math.Fix8
	jumpHeight math.Fix8
}

func NewScene() *Scene {
	return &Scene{
		scrollSpeed: math.FixOne + math.FixEighth,
		gravity:     math.FixQuarter,
		ground:      math.FixOne * 147,
		jumpHeight:  -math.FixOne * 3,
	}
}

func (s *Scene) Init(e *game.Engine) error {
	s.pillars = newPillarBG(100, e.NewBackground(assets.PillarsTileMap, display.Priority1))
	err := s.pillars.Show()
	if err != nil {
		return err
	}

	s.sky = e.NewBackground(assets.SkyTileMap, display.Priority3)
	err = s.sky.Add()
	if err != nil {
		return err
	}

	s.clouds = e.NewBackground(assets.CloudsTileMap, display.Priority2)
	err = s.clouds.Add()
	if err != nil {
		return err
	}
	s.clouds.VScroll = 20

	s.player = actor.NewPlayer(math.FixOne*40, math.FixOne*62, e.NewSprite(assets.PlayerTileSet))
	err = s.player.Show()
	if err != nil {
		return err
	}

	s.counter = score.NewCounter(97, 24, e)

	return nil
}

func (s *Scene) Update(e *game.Engine, frame int) error {
	var jump math.Fix8
	if key.JustPressed(key.A) {
		s.pillars.start()
		s.player.Start()
		jump = -math.FixOne * 3
	}

	s.player.Update(s.gravity, jump)
	if s.player.Rect().Y2 >= s.ground.Int() {
		s.gameOver = true
	}

	s.sky.HScroll += s.scrollSpeed / 3
	err := s.sky.Add()
	if err != nil {
		return err
	}

	s.clouds.HScroll += s.scrollSpeed / 2
	err = s.clouds.Add()
	if err != nil {
		return err
	}

	s.pillars.Update(s.scrollSpeed)
	err = s.pillars.Show()
	if err != nil {
		return err
	}

	if s.pillars.checkPoint(s.player.Rect()) {
		s.counter.Add()
	}

	s.counter.Draw()

	if s.pillars.collisionCheck(s.player.Rect()) {
		s.gameOver = true
	}

	return nil
}

func (s *Scene) Next() (game.Runable, bool) {
	if s.gameOver {
		return gameover.NewScene(
			s.sky,
			s.clouds,
			s.pillars.bg,
			s.player,
			s.counter,
		), true
	}
	return nil, false
}

type pillarBG struct {
	bg          *game.Background
	rand        *rand.Rand
	nextPillar  int
	pillarEvery int
	pillars     []math.Rect
	gapSize     int
	lastPoint   int

	started bool
}

func newPillarBG(pillarEvery int, bg *game.Background) *pillarBG {
	pillars := &pillarBG{
		bg:          bg,
		gapSize:     7,
		pillarEvery: pillarEvery,
	}

	return pillars
}

func (p *pillarBG) checkPoint(check math.Rect) bool {
	buffer := 4
	for i := range p.pillars {
		if p.pillars[i].X2 == p.lastPoint {
			continue
		}
		right := p.pillars[i].X2 - p.bg.HScroll.Int()
		if right < check.X1+buffer {
			p.lastPoint = p.pillars[i].X2
			return true
		}
	}
	return false
}

func (p *pillarBG) start() {
	if p.started {
		return
	}

	p.started = true
	p.lastPoint = p.bg.HScroll.Int() + p.pillarEvery
}

func (p *pillarBG) collisionCheck(check math.Rect) bool {
	for i := range p.pillars {
		left := p.pillars[i].X1 - p.bg.HScroll.Int()
		if left <= 0 {
			continue
		}
		right := left + 32
		top := p.pillars[i].Y1
		bottom := p.pillars[i].Y2

		if check.X2 < left || check.X1 > right {
			continue
		}
		if check.Y1 < top || check.Y2 > bottom {
			return true
		}
	}

	return false
}

func (p *pillarBG) addPillar(x int) math.Rect {
	start := (x % 512) / 8
	columns := [4]int{start, (start + 1) % 64, (start + 2) % 64, (start + 3) % 64}

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

	return math.Rect{X1: x, Y1: gap*8 + 4, X2: x + 32, Y2: (gap+p.gapSize)*8 + 4}
}

func (p *pillarBG) removePillar(r math.Rect) {
	start := (r.X1 % 512) / 8
	columns := [4]int{start, (start + 1) % 64, (start + 2) % 64, (start + 3) % 64}

	for i := 0; i < 18; i++ {
		for j := 0; j < 4; j++ {
			p.bg.SetTile(columns[j], i, 0)
		}
	}
}

func (p *pillarBG) Update(scrollSpeed math.Fix8) {
	p.bg.HScroll += scrollSpeed

	if !p.started {
		return
	}
	if p.rand == nil {
		p.rand = rand.New(rand.NewSource(int64(p.bg.HScroll)))

	}

	// add pillars to the right just off screen
	p.nextPillar--
	var keep []math.Rect
	if p.nextPillar <= 0 {
		x := p.bg.HScroll.Int() + 256
		r := p.addPillar(x)
		p.nextPillar = p.pillarEvery
		keep = []math.Rect{r}
	}

	// remove current pillars to the left of the screen
	border := p.bg.HScroll.Int() - 32
	for i := range p.pillars {
		if p.pillars[i].X1 < border {
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
