package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"

	lua "github.com/yuin/gopher-lua"
	
	// Importar el generador
	"deepthinking.do/luango/engine/generator"
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

func registerGameFunctions(L *lua.LState, em *EntityManager, player *Entity, agent *generator.GenerativeAgent) {
	L.SetGlobal("log", L.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		fmt.Println("[Lua]:", msg)
		return 0
	}))

	L.SetGlobal("debug", L.NewFunction(func(L *lua.LState) int {
		val := L.ToString(1)
		fmt.Println("[DEBUG]:", val)
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
	
	// ============= NUEVAS FUNCIONES DEL AGENTE GENERATIVO =============
	
	L.SetGlobal("generate", L.NewFunction(func(L *lua.LState) int {
		resourceType := L.ToString(1)
		prompt := L.ToString(2)
		
		result, err := agent.GenerateFromPrompt(prompt, generator.ResourceType(resourceType))
		if err != nil {
			fmt.Printf("[Agent Error]: %v\n", err)
			L.Push(lua.LNil)
			return 1
		}
		
		// Retornar ID del resultado para poder aplicarlo luego
		L.Push(lua.LString(result.ID))
		return 1
	}))
	
	L.SetGlobal("apply_generated", L.NewFunction(func(L *lua.LState) int {
		resultID := L.ToString(1)
		
		pending := agent.GetPendingResults()
		result, exists := pending[resultID]
		if !exists {
			fmt.Printf("[Agent Error]: No pending result with ID: %s\n", resultID)
			L.Push(lua.LBool(false))
			return 1
		}
		
		err := agent.ApplyResult(result)
		if err != nil {
			fmt.Printf("[Agent Error]: %v\n", err)
			L.Push(lua.LBool(false))
			return 1
		}
		
		L.Push(lua.LBool(true))
		return 1
	}))
	
	L.SetGlobal("set_agent_mode", L.NewFunction(func(L *lua.LState) int {
		mode := L.ToString(1)
		agent.SetMode(generator.AgentMode(mode))
		return 0
	}))

	// ============= FUNCIONES EXISTENTES =============
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
	case "F3":
		return ebiten.KeyF3
	case "F4":
		return ebiten.KeyF4
	default:
		return 0
	}
}

type Game struct {
	luaState    *lua.LState
	em          *EntityManager
	player      *Entity
	agent       *generator.GenerativeAgent
	agentCLI    *generator.AgentCLI
	started     bool
	showDebug   bool
	showAgent   bool
	frame       int
}

func (g *Game) Update() error {
	g.frame++
	if ebiten.IsKeyPressed(ebiten.KeyF3) {
		g.showDebug = !g.showDebug
		time.Sleep(200 * time.Millisecond)
	}
	
	// Alternar agente con F4
	if ebiten.IsKeyPressed(ebiten.KeyF4) {
		g.showAgent = !g.showAgent
		if g.showAgent {
			fmt.Println("\n🤖 Agente Generativo Activado - Usa la consola para interactuar")
			go g.agentCLI.Start()
		} else {
			g.agentCLI.Stop()
		}
		time.Sleep(200 * time.Millisecond)
	}

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

	if g.showDebug {
		text.Draw(screen, "Luango Generative Engine - Debug Mode", basicfont.Face7x13, 10, 20, color.White)
		text.Draw(screen, fmt.Sprintf("Player Pos: X=%.0f Y=%.0f", g.player.Position.X, g.player.Position.Y), basicfont.Face7x13, 10, 40, color.White)
		text.Draw(screen, fmt.Sprintf("Frame: %d", g.frame), basicfont.Face7x13, 10, 60, color.White)
		
		pending := len(g.agent.GetPendingResults())
		if pending > 0 {
			text.Draw(screen, fmt.Sprintf("🤖 Pending: %d", pending), basicfont.Face7x13, 10, 80, color.RGBA{255, 255, 0, 255})
		}
		
		text.Draw(screen, "F3: Toggle Debug | F4: Toggle Agent", basicfont.Face7x13, 10, 460, color.RGBA{128, 128, 128, 255})
	}
	
	if g.showAgent {
		// Mostrar indicador de que el agente está activo
		text.Draw(screen, "🤖 AGENT ACTIVE", basicfont.Face7x13, 500, 20, color.RGBA{0, 255, 0, 255})
		text.Draw(screen, "Check console", basicfont.Face7x13, 500, 40, color.RGBA{0, 255, 0, 255})
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
	fmt.Println("🚀 [Engine] Starting Luango Generative Engine")
	
	// Inicializar Lua
	L := lua.NewState()
	defer L.Close()

	// Inicializar EntityManager y Player
	em := NewEntityManager()
	playerSprite := loadSprite("assets/sprites/player.png")
	player := em.CreateEntity("Player", playerSprite)
	player.Position.X = 100
	player.Position.Y = 100

	// ============= INICIALIZAR AGENTE GENERATIVO =============
	projectPath, _ := os.Getwd()
	agent := generator.NewGenerativeAgent(projectPath, generator.ModeInteractive)
	
	// Registrar generadores
	luaGenerator := generator.NewLuaScriptGenerator(projectPath)
	agent.RegisterGenerator(generator.ResourceScript, luaGenerator)
	// Aquí se pueden registrar más generadores (sprites, sonidos, etc.)
	
	// Crear CLI del agente
	agentCLI := generator.NewAgentCLI(agent)
	
	// Registrar callbacks para notificaciones
	agent.RegisterCallback(generator.ResourceScript, func(result *generator.GenerationResult) {
		fmt.Printf("🎯 [Callback] Script generado: %s\n", result.Request.Prompt)
	})

	// ============= REGISTRAR FUNCIONES Y EJECUTAR SCRIPTS =============
	registerGameFunctions(L, em, player, agent)
	runLuaScripts("mod", L)

	// ============= INICIALIZAR JUEGO =============
	game := &Game{
		luaState: L, 
		em: em, 
		player: player,
		agent: agent,
		agentCLI: agentCLI,
	}
	
	fmt.Println("🎮 Controles:")
	fmt.Println("   F3: Toggle Debug Info")
	fmt.Println("   F4: Toggle Generative Agent")
	fmt.Println("   WASD/Arrows: Move Player")
	
	// Ejecutar juego
	err := ebiten.RunGame(game)
	if err != nil {
		fmt.Println("[Engine Error]:", err)
	}
}
