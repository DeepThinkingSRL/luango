package main

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/hajimehoshi/ebiten/v2"

	lua "github.com/yuin/gopher-lua"
)

type EntityID int

type Entity struct {
	ID       EntityID
	Name     string
	Position struct {
		X, Y float64
	}
	Sprite *ebiten.Image
}

type EntityManager struct {
	entities map[EntityID]*Entity
	nextID   EntityID
	lock     sync.Mutex
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		entities: make(map[EntityID]*Entity),
		nextID:   1,
	}
}

func (em *EntityManager) CreateEntity(name string, sprite *ebiten.Image) *Entity {
	em.lock.Lock()
	defer em.lock.Unlock()

	e := &Entity{
		ID:     em.nextID,
		Name:   name,
		Sprite: sprite,
	}
	em.entities[e.ID] = e
	em.nextID++
	return e
}

func (em *EntityManager) GetEntity(id EntityID) (*Entity, bool) {
	em.lock.Lock()
	defer em.lock.Unlock()
	e, ok := em.entities[id]
	return e, ok
}

func registerGameFunctions(L *lua.LState, em *EntityManager, player *Entity) {
	L.SetGlobal("log", L.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		fmt.Println("[Lua]:", msg)
		return 0
	}))

	L.SetGlobal("emit", L.NewFunction(func(L *lua.LState) int {
		event := L.ToString(1)
		payload := L.ToString(2)
		fmt.Printf("[Event] %s -> %s\n", event, payload)
		return 0
	}))

	L.SetGlobal("play_sound", L.NewFunction(func(L *lua.LState) int {
		path := L.ToString(1)
		f, err := os.Open(path)
		if err != nil {
			fmt.Println("[Audio Error]:", err)
			return 0
		}
		defer f.Close()

		streamer, format, err := wav.Decode(f)
		if err != nil {
			fmt.Println("[Audio Error]:", err)
			return 0
		}
		defer streamer.Close()

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		speaker.Play(streamer)
		fmt.Printf("[Audio] Playing sound: %s\n", path)
		return 0
	}))

	L.SetGlobal("is_key_pressed", L.NewFunction(func(L *lua.LState) int {
		key := L.ToString(1)
		pressed := ebiten.IsKeyPressed(keyFromString(key))
		L.Push(lua.LBool(pressed))
		return 1
	}))

	L.SetGlobal("move_player", L.NewFunction(func(L *lua.LState) int {
		dx := L.ToNumber(1)
		dy := L.ToNumber(2)
		player.Position.X += float64(dx)
		player.Position.Y += float64(dy)
		return 0
	}))
}

func keyFromString(k string) ebiten.Key {
	switch k {
	case "ArrowUp":
		return ebiten.KeyArrowUp
	case "ArrowDown":
		return ebiten.KeyArrowDown
	case "ArrowLeft":
		return ebiten.KeyArrowLeft
	case "ArrowRight":
		return ebiten.KeyArrowRight
	case "Space":
		return ebiten.KeySpace
	case "W":
		return ebiten.KeyW
	case "A":
		return ebiten.KeyA
	case "S":
		return ebiten.KeyS
	case "D":
		return ebiten.KeyD
	default:
		return 0
	}
}

type Game struct {
	luaState *lua.LState
	em       *EntityManager
	player   *Entity
	started  bool
}

func (g *Game) Update() error {
	if !g.started {
		g.started = true
		if fn := g.luaState.GetGlobal("on_start"); fn.Type() == lua.LTFunction {
			if err := g.luaState.CallByParam(lua.P{Fn: fn, NRet: 0, Protect: true}); err != nil {
				fmt.Println("[Lua Error][on_start]:", err)
			}
		}
	}

	if fn := g.luaState.GetGlobal("on_update"); fn.Type() == lua.LTFunction {
		if err := g.luaState.CallByParam(lua.P{Fn: fn, NRet: 0, Protect: true}); err != nil {
			fmt.Println("[Lua Error][on_update]:", err)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, e := range g.em.entities {
		if e.Sprite != nil {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(e.Position.X, e.Position.Y)
			screen.DrawImage(e.Sprite, opts)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func loadSprite(path string) *ebiten.Image {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("[Sprite Error]:", err)
		return nil
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		fmt.Println("[Decode Error]:", err)
		return nil
	}

	ebitenImg := ebiten.NewImageFromImage(img)
	return ebitenImg
}

func runLuaScripts(folder string, L *lua.LState) {
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("[Mod Walk Error]:", err)
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".lua" {
			fmt.Printf("[Mod] Loading: %s\n", path)
			if err := L.DoFile(path); err != nil {
				fmt.Printf("[Lua Error]: %v\n", err)
			}
		}
		return nil
	})
}

func main() {
	fmt.Println("[Engine] Starting Luango Engine")
	L := lua.NewState()
	defer L.Close()

	em := NewEntityManager()
	playerSprite := loadSprite("assets/sprites/player.png")
	player := em.CreateEntity("Player", playerSprite)
	player.Position.X = 100
	player.Position.Y = 100

	registerGameFunctions(L, em, player)
	runLuaScripts("mod", L)

	game := &Game{luaState: L, em: em, player: player}
	egOpts := ebiten.RunGame(game)
	if egOpts != nil {
		fmt.Println("[Engine Error]:", egOpts)
	}
}
