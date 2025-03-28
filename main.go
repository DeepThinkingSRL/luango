package main

import (
	"fmt"
	"os"

	lua "github.com/yuin/gopher-lua"
)

// Expose Go function to Lua
func registerGameFunctions(L *lua.LState) {
	L.SetGlobal("log", L.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		fmt.Println("[Lua]:", msg)
		return 0
	}))
}

func runLuaScript(path string, L *lua.LState) {
	if err := L.DoFile(path); err != nil {
		fmt.Println("Error running Lua script:", err)
		os.Exit(1)
	}
}

func main() {
	fmt.Println("[Engine] Starting Game Engine")
	L := lua.NewState()
	defer L.Close()

	registerGameFunctions(L)

	runLuaScript("mod/main.lua", L)
}
