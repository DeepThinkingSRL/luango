package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"deepthinking.do/luango/engine"
)

func main() {
	fmt.Println("🚀 [Engine] Starting Luango Engine - Modular Architecture")

	// Set window properties
	ebiten.SetWindowSize(1200, 800)
	ebiten.SetWindowTitle("Luango Engine - Modular Editor")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Create and initialize game
	game := engine.NewGame()
	if err := game.Initialize(); err != nil {
		fmt.Printf("Failed to initialize game: %v\n", err)
		return
	}
	defer game.Close()

	fmt.Println("🎮 Controls:")
	fmt.Println("   F1: Toggle Editor/Play Mode")
	fmt.Println("   F2: Toggle Inspector (Editor mode only)")
	fmt.Println("   F3: Toggle Debug Info")
	fmt.Println("   F11: Toggle Fullscreen")
	fmt.Println("   WASD/Arrows: Pan Camera (Editor mode) / Move Player (Play mode)")
	fmt.Println("   Mouse Wheel/+/-: Zoom")
	fmt.Println("   R: Reset Camera")
	fmt.Println("   Drag: Move Entities (Editor mode)")
	fmt.Println("   Middle Mouse: Pan Camera")

	// Run game
	if err := ebiten.RunGame(game); err != nil {
		fmt.Printf("[Engine Error]: %v\n", err)
	}
}
