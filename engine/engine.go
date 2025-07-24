package engine

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"deepthinking.do/luengo/engine/audio"
	"deepthinking.do/luengo/engine/camera"
	"deepthinking.do/luengo/engine/entity"
	"deepthinking.do/luengo/engine/input"
	"deepthinking.do/luengo/engine/resources"
	"deepthinking.do/luengo/engine/scripting"
	"deepthinking.do/luengo/engine/ui"
)

type Game struct {
	// Core systems
	entityManager   *entity.Manager
	camera          camera.Camera
	inputManager    *input.Manager
	audioManager    *audio.Manager
	scriptManager   *scripting.Manager
	resourceManager *resources.Manager
	ui              *ui.EditorUI

	// Game state
	player     *entity.Entity
	started    bool
	frame      int
	editorMode bool

	// Editor state
	isDragging bool
	dragStart  struct{ X, Y int }
	dragEntity *entity.Entity

	// Window state
	screenWidth, screenHeight int
}

func NewGame() *Game {
	// Initialize managers
	entityManager := entity.NewManager()
	inputManager := input.NewManager()
	audioManager := audio.NewManager()
	resourceManager := resources.NewManager()
	scriptManager := scripting.NewManager(audioManager)
	ui := ui.NewEditorUI()

	return &Game{
		entityManager:   entityManager,
		camera:          camera.NewCamera(),
		inputManager:    inputManager,
		audioManager:    audioManager,
		scriptManager:   scriptManager,
		resourceManager: resourceManager,
		ui:              ui,
		editorMode:      true,
		screenWidth:     1200,
		screenHeight:    800,
	}
}

func (g *Game) Initialize() error {
	// Initialize all systems
	g.inputManager.Initialize()

	if err := g.audioManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize audio: %w", err)
	}

	// Load player sprite
	playerSprite, err := g.resourceManager.LoadSprite("assets/sprites/player.png")
	if err != nil {
		fmt.Printf("Warning: Could not load player sprite: %v\n", err)
		// Create a simple colored rectangle as fallback
		playerSprite = ebiten.NewImage(32, 32)
		playerSprite.Fill(color.RGBA{0, 255, 0, 255})
	}

	// Create player entity
	g.player = g.entityManager.CreateEntity("Player", playerSprite)
	g.player.Position.X = 100
	g.player.Position.Y = 100

	// Create some test entities
	g.createTestEntities(playerSprite)

	// Register Lua functions and load scripts
	g.scriptManager.RegisterGameFunctions(g.entityManager, g.player)
	if err := g.scriptManager.LoadScriptsFromFolder("mod"); err != nil {
		fmt.Printf("Warning: Could not load scripts: %v\n", err)
	}

	// Initial log messages
	g.ui.AddLogMessage("Luengo Engine initialized (Modular)", g.frame)
	g.ui.AddLogMessage("Started in editor mode", g.frame)
	g.ui.AddLogMessage("F1: Toggle Editor/Play mode", g.frame)
	g.ui.AddLogMessage("F2: Toggle Inspector", g.frame)
	g.ui.AddLogMessage("F3: Toggle Debug info", g.frame)
	g.ui.AddLogMessage("F11: Toggle Fullscreen", g.frame)

	return nil
}

func (g *Game) createTestEntities(sprite *ebiten.Image) {
	testEntity1 := g.entityManager.CreateEntity("TestBox1", sprite)
	testEntity1.Position.X = 200
	testEntity1.Position.Y = 150

	testEntity2 := g.entityManager.CreateEntity("TestBox2", sprite)
	testEntity2.Position.X = 300
	testEntity2.Position.Y = 200

	testEntity3 := g.entityManager.CreateEntity("TestBox3", sprite)
	testEntity3.Position.X = 150
	testEntity3.Position.Y = 300
}

func (g *Game) Close() {
	g.scriptManager.Close()
}

func (g *Game) Update() error {
	g.frame++
	g.inputManager.Update()

	// Update screen size
	g.screenWidth, g.screenHeight = ebiten.WindowSize()

	// Handle input
	g.handleInput()

	// Handle mode-specific logic
	if g.editorMode {
		g.handleEditorMode()
	} else {
		g.handlePlayMode()
	}

	return nil
}

func (g *Game) handleInput() {
	// Toggle debug info with F3
	if g.inputManager.IsKeyJustPressed(ebiten.KeyF3) {
		g.ui.ToggleDebug()
		g.ui.AddLogMessage(fmt.Sprintf("Debug mode: %t", g.ui.IsDebugVisible()), g.frame)
	}

	// Toggle editor mode with F1
	if g.inputManager.IsKeyJustPressed(ebiten.KeyF1) {
		g.editorMode = !g.editorMode
		mode := "Editor"
		if !g.editorMode {
			mode = "Play"
		}
		g.ui.AddLogMessage(fmt.Sprintf("Switched to %s mode", mode), g.frame)
	}

	// Toggle fullscreen with F11
	if g.inputManager.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
		g.ui.AddLogMessage("Toggled fullscreen", g.frame)
	}

	// Toggle inspector with F2 (only in editor mode)
	if g.editorMode && g.inputManager.IsKeyJustPressed(ebiten.KeyF2) {
		g.ui.ToggleInspector()
		status := "closed"
		if g.ui.IsInspectorOpen() {
			status = "opened"
		}
		g.ui.AddLogMessage(fmt.Sprintf("Inspector %s", status), g.frame)
	}
}

func (g *Game) handleEditorMode() {
	g.handleCameraControls()
	g.handleMouseInteraction()
}

func (g *Game) handlePlayMode() {
	g.handlePlayerMovement()

	// Run Lua scripts
	if !g.started {
		g.started = true
		if err := g.scriptManager.CallFunction("on_start"); err != nil {
			g.ui.AddLogMessage(err.Error(), g.frame)
		}
	}

	if err := g.scriptManager.CallFunction("on_update"); err != nil {
		g.ui.AddLogMessage(err.Error(), g.frame)
	}
}

func (g *Game) handleCameraControls() {
	moveSpeed := 5.0 / g.camera.Zoom

	if g.inputManager.IsKeyPressed(ebiten.KeyArrowLeft) || g.inputManager.IsKeyPressed(ebiten.KeyA) {
		g.camera.Move(-moveSpeed, 0)
	}
	if g.inputManager.IsKeyPressed(ebiten.KeyArrowRight) || g.inputManager.IsKeyPressed(ebiten.KeyD) {
		g.camera.Move(moveSpeed, 0)
	}
	if g.inputManager.IsKeyPressed(ebiten.KeyArrowUp) || g.inputManager.IsKeyPressed(ebiten.KeyW) {
		g.camera.Move(0, -moveSpeed)
	}
	if g.inputManager.IsKeyPressed(ebiten.KeyArrowDown) || g.inputManager.IsKeyPressed(ebiten.KeyS) {
		g.camera.Move(0, moveSpeed)
	}

	// Zoom controls
	_, wheelY := g.inputManager.GetWheelDelta()
	if wheelY > 0 || g.inputManager.IsKeyJustPressed(ebiten.KeyEqual) || g.inputManager.IsKeyJustPressed(ebiten.KeyKPAdd) {
		g.camera.ZoomBy(1.1)
		g.ui.AddLogMessage(fmt.Sprintf("Zoom: %.2fx", g.camera.Zoom), g.frame)
	}
	if wheelY < 0 || g.inputManager.IsKeyJustPressed(ebiten.KeyMinus) || g.inputManager.IsKeyJustPressed(ebiten.KeyKPSubtract) {
		g.camera.ZoomBy(1.0 / 1.1)
		g.ui.AddLogMessage(fmt.Sprintf("Zoom: %.2fx", g.camera.Zoom), g.frame)
	}

	// Reset camera
	if g.inputManager.IsKeyJustPressed(ebiten.KeyR) {
		g.camera.Reset()
		g.ui.AddLogMessage("Camera reset", g.frame)
	}
}

func (g *Game) handlePlayerMovement() {
	if g.player == nil {
		return
	}

	moveSpeed := 3.0
	moved := false

	if g.inputManager.IsKeyPressed(ebiten.KeyArrowLeft) || g.inputManager.IsKeyPressed(ebiten.KeyA) {
		g.player.Position.X -= moveSpeed
		moved = true
	}
	if g.inputManager.IsKeyPressed(ebiten.KeyArrowRight) || g.inputManager.IsKeyPressed(ebiten.KeyD) {
		g.player.Position.X += moveSpeed
		moved = true
	}
	if g.inputManager.IsKeyPressed(ebiten.KeyArrowUp) || g.inputManager.IsKeyPressed(ebiten.KeyW) {
		g.player.Position.Y -= moveSpeed
		moved = true
	}
	if g.inputManager.IsKeyPressed(ebiten.KeyArrowDown) || g.inputManager.IsKeyPressed(ebiten.KeyS) {
		g.player.Position.Y += moveSpeed
		moved = true
	}

	// Make camera follow player
	if moved {
		g.camera.FollowTarget(g.player.Position.X, g.player.Position.Y, g.screenWidth, g.screenHeight, 0.1)
	}
}

func (g *Game) handleMouseInteraction() {
	mouseX, mouseY := g.inputManager.GetMousePosition()
	worldX, worldY := g.camera.ScreenToWorld(float64(mouseX), float64(mouseY))

	if g.inputManager.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.isDragging {
			g.isDragging = true
			g.dragStart.X = mouseX
			g.dragStart.Y = mouseY

			// Check if clicking on an entity
			entity := g.getEntityAt(worldX, worldY)
			if entity != nil {
				g.ui.SetSelectedEntity(entity)
				g.dragEntity = entity
				g.ui.AddLogMessage(fmt.Sprintf("Selected and dragging: %s", entity.Name), g.frame)
			} else {
				g.ui.SetSelectedEntity(nil)
				g.dragEntity = nil
			}
		} else {
			// Continue dragging
			if g.dragEntity != nil {
				deltaX, deltaY := g.inputManager.GetMouseDelta()
				g.dragEntity.Position.X += float64(deltaX) / g.camera.Zoom
				g.dragEntity.Position.Y += float64(deltaY) / g.camera.Zoom
			} else {
				// Pan camera
				deltaX, deltaY := g.inputManager.GetMouseDelta()
				g.camera.Move(-float64(deltaX)/g.camera.Zoom, -float64(deltaY)/g.camera.Zoom)
			}
		}
	} else {
		// Mouse button released
		if g.isDragging {
			g.isDragging = false
			if g.dragEntity != nil {
				g.ui.AddLogMessage(fmt.Sprintf("Finished dragging: %s", g.dragEntity.Name), g.frame)
				g.dragEntity = nil
			}
		}
	}

	// Middle mouse button for panning
	if g.inputManager.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		deltaX, deltaY := g.inputManager.GetMouseDelta()
		g.camera.Move(-float64(deltaX)/g.camera.Zoom, -float64(deltaY)/g.camera.Zoom)
	}
}

func (g *Game) getEntityAt(worldX, worldY float64) *entity.Entity {
	for _, e := range g.entityManager.GetAllEntities() {
		if e.Sprite != nil {
			w, h := e.Sprite.Bounds().Dx(), e.Sprite.Bounds().Dy()
			if worldX >= e.Position.X && worldX <= e.Position.X+float64(w) &&
				worldY >= e.Position.Y && worldY <= e.Position.Y+float64(h) {
				return e
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{30, 30, 35, 255})

	// Calculate viewport dimensions
	viewportWidth := g.screenWidth
	if g.ui.IsInspectorOpen() {
		viewportWidth -= 200
	}
	logPanelHeight := 120
	viewportHeight := g.screenHeight - logPanelHeight

	// Draw grid in editor mode
	if g.editorMode {
		g.ui.DrawGrid(screen, &g.camera, viewportWidth, viewportHeight)
	}

	// Draw entities
	g.drawEntities(screen, viewportWidth, viewportHeight)

	// Draw UI
	g.ui.DrawModeIndicator(screen, g.editorMode)
	g.ui.DrawCameraInfo(screen, &g.camera, g.editorMode)
	g.ui.DrawControls(screen, g.editorMode)
	g.ui.DrawDebugInfo(screen, g.player, g.frame, g.entityManager.Count(), g.screenWidth, g.screenHeight)

	if g.ui.IsInspectorOpen() {
		g.ui.DrawInspector(screen, &g.camera, g.screenWidth, g.screenHeight)
	}

	g.ui.DrawLogPanel(screen, g.screenWidth, g.screenHeight)
}

func (g *Game) drawEntities(screen *ebiten.Image, viewportWidth, viewportHeight int) {
	cameraMatrix := g.camera.GetTransformMatrix()

	for _, e := range g.entityManager.GetAllEntities() {
		if e.Sprite != nil {
			// Calculate screen position for culling
			screenX, screenY := g.camera.WorldToScreen(e.Position.X, e.Position.Y)

			// Cull entities outside viewport
			w, h := e.Sprite.Bounds().Dx(), e.Sprite.Bounds().Dy()
			scaledW, scaledH := float64(w)*g.camera.Zoom, float64(h)*g.camera.Zoom

			if screenX+scaledW >= 0 && screenX <= float64(viewportWidth) &&
				screenY+scaledH >= 0 && screenY <= float64(viewportHeight) {

				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Concat(cameraMatrix)
				opts.GeoM.Translate(e.Position.X, e.Position.Y)

				// Highlight selected entity
				if g.ui.GetSelectedEntity() == e {
					// Draw selection border
					borderOpts := &ebiten.DrawImageOptions{}
					borderOpts.GeoM.Concat(cameraMatrix)
					borderOpts.GeoM.Translate(e.Position.X-2/g.camera.Zoom, e.Position.Y-2/g.camera.Zoom)

					borderW := int(float64(w) + 4/g.camera.Zoom)
					borderH := int(float64(h) + 4/g.camera.Zoom)
					if borderW > 0 && borderH > 0 {
						border := ebiten.NewImage(borderW, borderH)
						border.Fill(color.RGBA{255, 255, 0, 255})
						screen.DrawImage(border, borderOpts)
					}
				}

				screen.DrawImage(e.Sprite, opts)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screenWidth = outsideWidth
	g.screenHeight = outsideHeight
	return outsideWidth, outsideHeight
}
