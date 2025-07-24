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

func (em *EntityManager) GetAllEntities() map[EntityID]*Entity {
	em.lock.Lock()
	defer em.lock.Unlock()
	// Return a copy to avoid race conditions
	result := make(map[EntityID]*Entity)
	for k, v := range em.entities {
		result[k] = v
	}
	return result
}

type Camera struct {
	X, Y float64  // Camera position
	Zoom float64  // Zoom level (1.0 = normal)
}

func NewCamera() Camera {
	return Camera{
		X: 0,
		Y: 0,
		Zoom: 1.0,
	}
}

type Game struct {
	luaState      *lua.LState
	em            *EntityManager
	player        *Entity
	started       bool
	showDebug     bool
	frame         int
	
	// Editor UI state
	editorMode    bool  // true = editor mode, false = play mode
	selectedEntity *Entity
	inspectorOpen bool
	logMessages   []string
	maxLogMessages int
	
	// Camera and viewport
	camera Camera
	isDragging bool
	dragStart struct{ X, Y int }
	dragEntity *Entity
	lastMouseX, lastMouseY int
	
	// Window state
	screenWidth, screenHeight int
	
	// Input state tracking for proper key detection
	keyStates map[ebiten.Key]bool
	lastKeyStates map[ebiten.Key]bool
}

func registerGameFunctions(L *lua.LState, em *EntityManager, player *Entity) {
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
	default:
		return 0
	}
}

// addLogMessage adds a message to the log with automatic trimming
func (g *Game) addLogMessage(msg string) {
	g.logMessages = append(g.logMessages, fmt.Sprintf("[%d] %s", g.frame, msg))
	if len(g.logMessages) > g.maxLogMessages {
		g.logMessages = g.logMessages[1:]
	}
	fmt.Println(msg) // Also print to console
}

// initializeInputSystem initializes the input state tracking
func (g *Game) initializeInputSystem() {
	g.keyStates = make(map[ebiten.Key]bool)
	g.lastKeyStates = make(map[ebiten.Key]bool)
	
	// Initialize tracking for important keys
	keys := []ebiten.Key{
		ebiten.KeyF1, ebiten.KeyF2, ebiten.KeyF3, ebiten.KeyF11,
		ebiten.KeyArrowUp, ebiten.KeyArrowDown, ebiten.KeyArrowLeft, ebiten.KeyArrowRight,
		ebiten.KeyW, ebiten.KeyA, ebiten.KeyS, ebiten.KeyD,
		ebiten.KeyR, ebiten.KeyEqual, ebiten.KeyMinus,
		ebiten.KeyKPAdd, ebiten.KeyKPSubtract,
	}
	
	for _, key := range keys {
		g.keyStates[key] = false
		g.lastKeyStates[key] = false
	}
}

// updateInputState updates the current key states
func (g *Game) updateInputState() {
	// Save previous states
	for key := range g.keyStates {
		g.lastKeyStates[key] = g.keyStates[key]
		g.keyStates[key] = ebiten.IsKeyPressed(key)
	}
}

// isKeyJustPressed returns true if key was just pressed (not held)
func (g *Game) isKeyJustPressed(key ebiten.Key) bool {
	return g.keyStates[key] && !g.lastKeyStates[key]
}

// handleCameraControls manages camera movement and zoom
func (g *Game) handleCameraControls() {
	// Camera movement with arrow keys or WASD (continuous for smooth movement)
	moveSpeed := 5.0 / g.camera.Zoom // Slower movement when zoomed in
	
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.camera.X -= moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.camera.X += moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.camera.Y -= moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.camera.Y += moveSpeed
	}
	
	// Zoom with mouse wheel or +/- keys (just pressed for discrete zoom levels)
	_, wheelY := ebiten.Wheel()
	if wheelY > 0 || g.isKeyJustPressed(ebiten.KeyEqual) || g.isKeyJustPressed(ebiten.KeyKPAdd) {
		g.camera.Zoom *= 1.1
		if g.camera.Zoom > 3.0 {
			g.camera.Zoom = 3.0
		}
		g.addLogMessage(fmt.Sprintf("Zoom: %.2fx", g.camera.Zoom))
	}
	if wheelY < 0 || g.isKeyJustPressed(ebiten.KeyMinus) || g.isKeyJustPressed(ebiten.KeyKPSubtract) {
		g.camera.Zoom /= 1.1
		if g.camera.Zoom < 0.5 {
			g.camera.Zoom = 0.5
		}
		g.addLogMessage(fmt.Sprintf("Zoom: %.2fx", g.camera.Zoom))
	}
	
	// Reset camera with R key (just pressed to prevent multiple resets)
	if g.isKeyJustPressed(ebiten.KeyR) {
		g.camera = NewCamera()
		g.addLogMessage("Camera reset")
	}
}

// screenToWorld converts screen coordinates to world coordinates
func (g *Game) screenToWorld(screenX, screenY float64) (float64, float64) {
	worldX := (screenX / g.camera.Zoom) + g.camera.X
	worldY := (screenY / g.camera.Zoom) + g.camera.Y
	return worldX, worldY
}

// worldToScreen converts world coordinates to screen coordinates
func (g *Game) worldToScreen(worldX, worldY float64) (float64, float64) {
	screenX := (worldX - g.camera.X) * g.camera.Zoom
	screenY := (worldY - g.camera.Y) * g.camera.Zoom
	return screenX, screenY
}

// getEntityAt returns the entity at the given world coordinates
func (g *Game) getEntityAt(worldX, worldY float64) *Entity {
	for _, entity := range g.em.entities {
		if entity.Sprite != nil {
			w, h := entity.Sprite.Bounds().Dx(), entity.Sprite.Bounds().Dy()
			if worldX >= entity.Position.X && worldX <= entity.Position.X+float64(w) &&
			   worldY >= entity.Position.Y && worldY <= entity.Position.Y+float64(h) {
				return entity
			}
		}
	}
	return nil
}

// handleMouseInteraction manages mouse-based interactions
func (g *Game) handleMouseInteraction() {
	mouseX, mouseY := ebiten.CursorPosition()
	
	// Convert screen coordinates to world coordinates
	worldX, worldY := g.screenToWorld(float64(mouseX), float64(mouseY))
	
	// Handle mouse button press
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !g.isDragging {
			// Start new interaction
			g.isDragging = true
			g.dragStart.X = mouseX
			g.dragStart.Y = mouseY
			
			// Check if clicking on an entity
			entity := g.getEntityAt(worldX, worldY)
			if entity != nil {
				g.selectedEntity = entity
				g.dragEntity = entity
				g.addLogMessage(fmt.Sprintf("Selected and dragging: %s", entity.Name))
			} else {
				// Start camera panning
				g.selectedEntity = nil
				g.dragEntity = nil
			}
		} else {
			// Continue dragging
			if g.dragEntity != nil {
				// Drag entity
				deltaX := float64(mouseX - g.lastMouseX) / g.camera.Zoom
				deltaY := float64(mouseY - g.lastMouseY) / g.camera.Zoom
				g.dragEntity.Position.X += deltaX
				g.dragEntity.Position.Y += deltaY
			} else {
				// Pan camera
				deltaX := float64(mouseX - g.lastMouseX) / g.camera.Zoom
				deltaY := float64(mouseY - g.lastMouseY) / g.camera.Zoom
				g.camera.X -= deltaX
				g.camera.Y -= deltaY
			}
		}
	} else {
		// Mouse button released
		if g.isDragging {
			g.isDragging = false
			if g.dragEntity != nil {
				g.addLogMessage(fmt.Sprintf("Finished dragging: %s", g.dragEntity.Name))
				g.dragEntity = nil
			}
		}
	}
	
	// Handle middle mouse button for panning
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		if g.lastMouseX != 0 || g.lastMouseY != 0 {
			deltaX := float64(mouseX - g.lastMouseX) / g.camera.Zoom
			deltaY := float64(mouseY - g.lastMouseY) / g.camera.Zoom
			g.camera.X -= deltaX
			g.camera.Y -= deltaY
		}
	}
	
	g.lastMouseX = mouseX
	g.lastMouseY = mouseY
}

func (g *Game) Update() error {
	g.frame++
	
	// Update input state first
	g.updateInputState()
	
	// Update screen size for fullscreen support
	g.screenWidth, g.screenHeight = ebiten.WindowSize()
	
	// Toggle debug info with F3
	if g.isKeyJustPressed(ebiten.KeyF3) {
		g.showDebug = !g.showDebug
		g.addLogMessage(fmt.Sprintf("Debug mode: %t", g.showDebug))
	}
	
	// Toggle editor mode with F1
	if g.isKeyJustPressed(ebiten.KeyF1) {
		g.editorMode = !g.editorMode
		g.addLogMessage(fmt.Sprintf("Switched to %s mode", map[bool]string{true: "Editor", false: "Play"}[g.editorMode]))
	}
	
	// Toggle fullscreen with F11
	if g.isKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
		g.addLogMessage("Toggled fullscreen")
	}
	
	// Toggle inspector with F2 (only in editor mode)
	if g.editorMode && g.isKeyJustPressed(ebiten.KeyF2) {
		g.inspectorOpen = !g.inspectorOpen
		g.addLogMessage(fmt.Sprintf("Inspector %s", map[bool]string{true: "opened", false: "closed"}[g.inspectorOpen]))
	}
	
	// Camera controls in editor mode
	if g.editorMode {
		g.handleCameraControls()
		g.handleMouseInteraction()
	}

	// Only run lua scripts in play mode
	if !g.editorMode {
		if !g.started {
			g.started = true
			if fn := g.luaState.GetGlobal("on_start"); fn.Type() == lua.LTFunction {
				if err := g.luaState.CallByParam(lua.P{Fn: fn, NRet: 0, Protect: true}); err != nil {
					g.addLogMessage(fmt.Sprintf("[Lua Error][on_start]: %v", err))
				}
			}
		}

		if fn := g.luaState.GetGlobal("on_update"); fn.Type() == lua.LTFunction {
			if err := g.luaState.CallByParam(lua.P{Fn: fn, NRet: 0, Protect: true}); err != nil {
				g.addLogMessage(fmt.Sprintf("[Lua Error][on_update]: %v", err))
			}
		}
	}
	
	return nil
}

// drawGrid draws a grid in the background for editor mode
func (g *Game) drawGrid(screen *ebiten.Image, viewportWidth, viewportHeight int) {
	gridSize := 50.0 // Grid cell size in world units
	gridColor := color.RGBA{60, 60, 70, 255}
	
	// Calculate grid lines to draw based on camera position and zoom
	startX := int((g.camera.X / gridSize)) - 1
	endX := int(((g.camera.X + float64(viewportWidth)/g.camera.Zoom) / gridSize)) + 1
	startY := int((g.camera.Y / gridSize)) - 1
	endY := int(((g.camera.Y + float64(viewportHeight)/g.camera.Zoom) / gridSize)) + 1
	
	// Draw vertical lines
	for x := startX; x <= endX; x++ {
		worldX := float64(x) * gridSize
		screenX, _ := g.worldToScreen(worldX, 0)
		
		if screenX >= 0 && screenX <= float64(viewportWidth) {
			line := ebiten.NewImage(1, viewportHeight)
			line.Fill(gridColor)
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(screenX, 0)
			screen.DrawImage(line, opts)
		}
	}
	
	// Draw horizontal lines
	for y := startY; y <= endY; y++ {
		worldY := float64(y) * gridSize
		_, screenY := g.worldToScreen(0, worldY)
		
		if screenY >= 0 && screenY <= float64(viewportHeight) {
			line := ebiten.NewImage(viewportWidth, 1)
			line.Fill(gridColor)
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(0, screenY)
			screen.DrawImage(line, opts)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{30, 30, 35, 255})
	
	// Calculate viewport dimensions
	viewportWidth := g.screenWidth
	if g.inspectorOpen {
		viewportWidth -= 200 // Leave space for inspector
	}
	logPanelHeight := 120
	viewportHeight := g.screenHeight - logPanelHeight
	
	// Create camera transform matrix
	var cameraMatrix ebiten.GeoM
	cameraMatrix.Scale(g.camera.Zoom, g.camera.Zoom)
	cameraMatrix.Translate(-g.camera.X*g.camera.Zoom, -g.camera.Y*g.camera.Zoom)
	
	// Draw grid in editor mode
	if g.editorMode {
		g.drawGrid(screen, viewportWidth, viewportHeight)
	}
	
	// Draw entities with camera transform
	for _, e := range g.em.entities {
		if e.Sprite != nil {
			// Calculate screen position
			screenX, screenY := g.worldToScreen(e.Position.X, e.Position.Y)
			
			// Cull entities outside viewport
			w, h := e.Sprite.Bounds().Dx(), e.Sprite.Bounds().Dy()
			scaledW, scaledH := float64(w)*g.camera.Zoom, float64(h)*g.camera.Zoom
			
			if screenX+scaledW >= 0 && screenX <= float64(viewportWidth) &&
			   screenY+scaledH >= 0 && screenY <= float64(viewportHeight) {
				
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Concat(cameraMatrix)
				opts.GeoM.Translate(e.Position.X, e.Position.Y)
				
				// Highlight selected entity
				if g.selectedEntity == e {
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
	
	// Draw UI panels (these are not affected by camera)
	if g.inspectorOpen {
		g.drawInspector(screen)
	}
	
	g.drawLogPanel(screen)
	
	// Draw mode indicator and controls
	modeText := "PLAY MODE"
	modeColor := color.RGBA{0, 255, 0, 255}
	if g.editorMode {
		modeText = "EDITOR MODE"
		modeColor = color.RGBA{255, 165, 0, 255}
	}
	text.Draw(screen, modeText, basicfont.Face7x13, 10, 20, modeColor)
	
	// Draw camera info in editor mode
	if g.editorMode {
		cameraInfo := fmt.Sprintf("Cam: %.1f,%.1f Zoom: %.2fx", g.camera.X, g.camera.Y, g.camera.Zoom)
		text.Draw(screen, cameraInfo, basicfont.Face7x13, 150, 20, color.RGBA{200, 200, 200, 255})
	}
	
	// Draw controls
	controlY := 40
	text.Draw(screen, "F1: Mode F2: Inspector F11: Fullscreen", basicfont.Face7x13, 10, controlY, color.RGBA{128, 128, 128, 255})
	if g.editorMode {
		text.Draw(screen, "WASD/Arrows: Pan  Wheel/+/-: Zoom  R: Reset", basicfont.Face7x13, 10, controlY+15, color.RGBA{128, 128, 128, 255})
		text.Draw(screen, "Drag: Move entity  Middle: Pan camera", basicfont.Face7x13, 10, controlY+30, color.RGBA{128, 128, 128, 255})
	}
	
	// Debug info
	if g.showDebug {
		debugY := 100
		text.Draw(screen, "Luango Engine - DEBUG", basicfont.Face7x13, 10, debugY, color.White)
		text.Draw(screen, fmt.Sprintf("Player Pos: X=%.0f Y=%.0f", g.player.Position.X, g.player.Position.Y), basicfont.Face7x13, 10, debugY+20, color.White)
		text.Draw(screen, fmt.Sprintf("Frame: %d", g.frame), basicfont.Face7x13, 10, debugY+40, color.White)
		text.Draw(screen, fmt.Sprintf("Entities: %d", len(g.em.entities)), basicfont.Face7x13, 10, debugY+60, color.White)
		text.Draw(screen, fmt.Sprintf("Screen: %dx%d", g.screenWidth, g.screenHeight), basicfont.Face7x13, 10, debugY+80, color.White)
		text.Draw(screen, "F3: Toggle Debug", basicfont.Face7x13, 10, g.screenHeight-20, color.RGBA{128, 128, 128, 255})
	}
}

// drawInspector draws the inspector panel on the right side
func (g *Game) drawInspector(screen *ebiten.Image) {
	inspectorWidth := 200
	inspectorX := g.screenWidth - inspectorWidth // Right side panel
	
	// Background
	inspectorBg := ebiten.NewImage(inspectorWidth, g.screenHeight)
	inspectorBg.Fill(color.RGBA{40, 40, 40, 200})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(inspectorX), 0)
	screen.DrawImage(inspectorBg, opts)
	
	// Title
	text.Draw(screen, "INSPECTOR", basicfont.Face7x13, inspectorX+10, 20, color.White)
	
	if g.selectedEntity != nil {
		y := 50
		text.Draw(screen, fmt.Sprintf("Name: %s", g.selectedEntity.Name), basicfont.Face7x13, inspectorX+10, y, color.White)
		y += 20
		text.Draw(screen, fmt.Sprintf("ID: %d", g.selectedEntity.ID), basicfont.Face7x13, inspectorX+10, y, color.White)
		y += 20
		text.Draw(screen, fmt.Sprintf("X: %.1f", g.selectedEntity.Position.X), basicfont.Face7x13, inspectorX+10, y, color.White)
		y += 20
		text.Draw(screen, fmt.Sprintf("Y: %.1f", g.selectedEntity.Position.Y), basicfont.Face7x13, inspectorX+10, y, color.White)
		
		if g.selectedEntity.Sprite != nil {
			y += 20
			w, h := g.selectedEntity.Sprite.Bounds().Dx(), g.selectedEntity.Sprite.Bounds().Dy()
			text.Draw(screen, fmt.Sprintf("Size: %dx%d", w, h), basicfont.Face7x13, inspectorX+10, y, color.White)
		}
		
		// Show if entity is being dragged
		if g.dragEntity == g.selectedEntity {
			y += 30
			text.Draw(screen, "DRAGGING", basicfont.Face7x13, inspectorX+10, y, color.RGBA{255, 255, 0, 255})
		}
		
		// Camera-relative position
		y += 30
		text.Draw(screen, "-- Camera View --", basicfont.Face7x13, inspectorX+10, y, color.RGBA{150, 150, 150, 255})
		y += 20
		screenX, screenY := g.worldToScreen(g.selectedEntity.Position.X, g.selectedEntity.Position.Y)
		text.Draw(screen, fmt.Sprintf("Screen: %.1f,%.1f", screenX, screenY), basicfont.Face7x13, inspectorX+10, y, color.RGBA{150, 150, 150, 255})
		
	} else {
		text.Draw(screen, "No entity selected", basicfont.Face7x13, inspectorX+10, 50, color.RGBA{128, 128, 128, 255})
		text.Draw(screen, "Click and drag", basicfont.Face7x13, inspectorX+10, 70, color.RGBA{128, 128, 128, 255})
		text.Draw(screen, "to select and move", basicfont.Face7x13, inspectorX+10, 90, color.RGBA{128, 128, 128, 255})
	}
	
	// Camera controls help
	y := g.screenHeight - 120
	text.Draw(screen, "-- Camera --", basicfont.Face7x13, inspectorX+10, y, color.RGBA{150, 150, 150, 255})
	y += 20
	text.Draw(screen, "WASD: Pan", basicfont.Face7x13, inspectorX+10, y, color.RGBA{120, 120, 120, 255})
	y += 15
	text.Draw(screen, "Wheel: Zoom", basicfont.Face7x13, inspectorX+10, y, color.RGBA{120, 120, 120, 255})
	y += 15
	text.Draw(screen, "R: Reset", basicfont.Face7x13, inspectorX+10, y, color.RGBA{120, 120, 120, 255})
}

func (g *Game) drawLogPanel(screen *ebiten.Image) {
	logHeight := 120
	logY := g.screenHeight - logHeight
	logWidth := g.screenWidth
	if g.inspectorOpen {
		logWidth -= 200 // Leave space for inspector
	}
	
	// Background
	logBg := ebiten.NewImage(logWidth, logHeight)
	logBg.Fill(color.RGBA{20, 20, 20, 200})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(0, float64(logY))
	screen.DrawImage(logBg, opts)
	
	// Title
	text.Draw(screen, "EXECUTION LOG", basicfont.Face7x13, 10, logY+20, color.White)
	
	// Log messages (show last few)
	maxVisible := 7
	start := len(g.logMessages) - maxVisible
	if start < 0 {
		start = 0
	}
	
	for i, msg := range g.logMessages[start:] {
		y := logY + 40 + i*12
		if y < logY+logHeight-10 {
			text.Draw(screen, msg, basicfont.Face7x13, 10, y, color.RGBA{200, 200, 200, 255})
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// Update screen dimensions
	g.screenWidth = outsideWidth
	g.screenHeight = outsideHeight
	return outsideWidth, outsideHeight
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
	fmt.Println("🚀 [Engine] Starting Luango Engine - Clean Architecture")
	
	// Set window properties
	ebiten.SetWindowSize(1200, 800)
	ebiten.SetWindowTitle("Luango Engine - Clean Editor")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	
	// Initialize Lua
	L := lua.NewState()
	defer L.Close()

	// Initialize EntityManager and Player
	em := NewEntityManager()
	playerSprite := loadSprite("assets/sprites/player.png")
	player := em.CreateEntity("Player", playerSprite)
	player.Position.X = 100
	player.Position.Y = 100

	// Create some additional entities for testing the editor
	testEntity1 := em.CreateEntity("TestBox1", playerSprite)
	testEntity1.Position.X = 200
	testEntity1.Position.Y = 150
	
	testEntity2 := em.CreateEntity("TestBox2", playerSprite)
	testEntity2.Position.X = 300
	testEntity2.Position.Y = 200
	
	testEntity3 := em.CreateEntity("TestBox3", playerSprite)
	testEntity3.Position.X = 150
	testEntity3.Position.Y = 300

	// Register game functions
	registerGameFunctions(L, em, player)
	runLuaScripts("mod", L)

	// Initialize game with editor features
	game := &Game{
		luaState:       L, 
		em:             em, 
		player:         player,
		editorMode:     true,  // Start in editor mode
		inspectorOpen:  true,  // Start with inspector open
		logMessages:    make([]string, 0),
		maxLogMessages: 50,
		camera:         NewCamera(),
		screenWidth:    1200,
		screenHeight:   800,
		keyStates:      make(map[ebiten.Key]bool),
		lastKeyStates:  make(map[ebiten.Key]bool),
	}
	
	// Initialize input system
	game.initializeInputSystem()
	
	// Add initial log messages
	game.addLogMessage("Luango Engine initialized (Clean)")
	game.addLogMessage("Started in editor mode")
	game.addLogMessage("F1: Toggle Editor/Play mode")
	game.addLogMessage("F2: Toggle Inspector")
	game.addLogMessage("F3: Toggle Debug info")
	game.addLogMessage("F11: Toggle Fullscreen")
	
	fmt.Println("🎮 Controls:")
	fmt.Println("   F1: Toggle Editor/Play Mode")
	fmt.Println("   F2: Toggle Inspector (Editor mode only)")
	fmt.Println("   F3: Toggle Debug Info")
	fmt.Println("   F11: Toggle Fullscreen")
	fmt.Println("   WASD/Arrows: Pan Camera (Editor mode)")
	fmt.Println("   Mouse Wheel/+/-: Zoom")
	fmt.Println("   R: Reset Camera")
	fmt.Println("   Drag: Move Entities (Editor mode)")
	fmt.Println("   Middle Mouse: Pan Camera")
	
	// Run game
	err := ebiten.RunGame(game)
	if err != nil {
		fmt.Println("[Engine Error]:", err)
	}
}
