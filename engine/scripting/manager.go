package scripting

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	lua "github.com/yuin/gopher-lua"

	"deepthinking.do/luengo/engine/audio"
	"deepthinking.do/luengo/engine/entity"
	"deepthinking.do/luengo/engine/input"
)

type Manager struct {
	luaState     *lua.LState
	audioManager *audio.Manager
}

func NewManager(audioManager *audio.Manager) *Manager {
	return &Manager{
		luaState:     lua.NewState(),
		audioManager: audioManager,
	}
}

func (sm *Manager) Close() {
	if sm.luaState != nil {
		sm.luaState.Close()
	}
}

func (sm *Manager) GetLuaState() *lua.LState {
	return sm.luaState
}

func (sm *Manager) RegisterGameFunctions(em *entity.Manager, player *entity.Entity) {
	L := sm.luaState

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
		if err := sm.audioManager.PlaySound(path); err != nil {
			fmt.Printf("[Audio Error]: %v\n", err)
		}
		return 0
	}))

	L.SetGlobal("is_key_pressed", L.NewFunction(func(L *lua.LState) int {
		key := L.ToString(1)
		pressed := ebiten.IsKeyPressed(input.KeyFromString(key))
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

func (sm *Manager) LoadScriptsFromFolder(folder string) error {
	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("[Mod Walk Error]: %v\n", err)
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".lua" {
			fmt.Printf("[Mod] Loading: %s\n", path)
			if err := sm.luaState.DoFile(path); err != nil {
				fmt.Printf("[Lua Error]: %v\n", err)
			}
		}
		return nil
	})
}

func (sm *Manager) CallFunction(functionName string) error {
	if fn := sm.luaState.GetGlobal(functionName); fn.Type() == lua.LTFunction {
		if err := sm.luaState.CallByParam(lua.P{Fn: fn, NRet: 0, Protect: true}); err != nil {
			return fmt.Errorf("[Lua Error][%s]: %v", functionName, err)
		}
	}
	return nil
}
