package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"

	"deepthinking.do/luango/engine/entity"
	"deepthinking.do/luango/engine/camera"
)

type EditorUI struct {
	inspectorOpen bool
	logMessages   []string
	maxLogMessages int
	selectedEntity *entity.Entity
	showDebug     bool
}

func NewEditorUI() *EditorUI {
	return &EditorUI{
		inspectorOpen:  true,
		logMessages:    make([]string, 0),
		maxLogMessages: 50,
	}
}

func (ui *EditorUI) AddLogMessage(msg string, frame int) {
	ui.logMessages = append(ui.logMessages, fmt.Sprintf("[%d] %s", frame, msg))
	if len(ui.logMessages) > ui.maxLogMessages {
		ui.logMessages = ui.logMessages[1:]
	}
	fmt.Println(msg) // Also print to console
}

func (ui *EditorUI) ToggleInspector() {
	ui.inspectorOpen = !ui.inspectorOpen
}

func (ui *EditorUI) ToggleDebug() {
	ui.showDebug = !ui.showDebug
}

func (ui *EditorUI) SetSelectedEntity(entity *entity.Entity) {
	ui.selectedEntity = entity
}

func (ui *EditorUI) GetSelectedEntity() *entity.Entity {
	return ui.selectedEntity
}

func (ui *EditorUI) IsInspectorOpen() bool {
	return ui.inspectorOpen
}

func (ui *EditorUI) IsDebugVisible() bool {
	return ui.showDebug
}

// DrawModeIndicator draws the current mode indicator
func (ui *EditorUI) DrawModeIndicator(screen *ebiten.Image, editorMode bool) {
	modeText := "PLAY MODE"
	modeColor := color.RGBA{0, 255, 0, 255}
	if editorMode {
		modeText = "EDITOR MODE"
		modeColor = color.RGBA{255, 165, 0, 255}
	}
	text.Draw(screen, modeText, basicfont.Face7x13, 10, 20, modeColor)
}

// DrawCameraInfo draws camera information in editor mode
func (ui *EditorUI) DrawCameraInfo(screen *ebiten.Image, cam *camera.Camera, editorMode bool) {
	if editorMode {
		cameraInfo := fmt.Sprintf("Cam: %.1f,%.1f Zoom: %.2fx", cam.X, cam.Y, cam.Zoom)
		text.Draw(screen, cameraInfo, basicfont.Face7x13, 150, 20, color.RGBA{200, 200, 200, 255})
	}
}

// DrawControls draws the control help text
func (ui *EditorUI) DrawControls(screen *ebiten.Image, editorMode bool) {
	controlY := 40
	text.Draw(screen, "F1: Mode F2: Inspector F11: Fullscreen", basicfont.Face7x13, 10, controlY, color.RGBA{128, 128, 128, 255})
	if editorMode {
		text.Draw(screen, "WASD/Arrows: Pan  Wheel/+/-: Zoom  R: Reset", basicfont.Face7x13, 10, controlY+15, color.RGBA{128, 128, 128, 255})
		text.Draw(screen, "Drag: Move entity  Middle: Pan camera", basicfont.Face7x13, 10, controlY+30, color.RGBA{128, 128, 128, 255})
	}
}

// DrawDebugInfo draws debug information
func (ui *EditorUI) DrawDebugInfo(screen *ebiten.Image, player *entity.Entity, frame int, entityCount int, screenWidth, screenHeight int) {
	if !ui.showDebug {
		return
	}
	
	debugY := 100
	text.Draw(screen, "Luango Engine - DEBUG", basicfont.Face7x13, 10, debugY, color.White)
	text.Draw(screen, fmt.Sprintf("Player Pos: X=%.0f Y=%.0f", player.Position.X, player.Position.Y), basicfont.Face7x13, 10, debugY+20, color.White)
	text.Draw(screen, fmt.Sprintf("Frame: %d", frame), basicfont.Face7x13, 10, debugY+40, color.White)
	text.Draw(screen, fmt.Sprintf("Entities: %d", entityCount), basicfont.Face7x13, 10, debugY+60, color.White)
	text.Draw(screen, fmt.Sprintf("Screen: %dx%d", screenWidth, screenHeight), basicfont.Face7x13, 10, debugY+80, color.White)
	text.Draw(screen, "F3: Toggle Debug", basicfont.Face7x13, 10, screenHeight-20, color.RGBA{128, 128, 128, 255})
}

// DrawInspector draws the inspector panel on the right side
func (ui *EditorUI) DrawInspector(screen *ebiten.Image, cam *camera.Camera, screenWidth, screenHeight int) {
	if !ui.inspectorOpen {
		return
	}
	
	inspectorWidth := 200
	inspectorX := screenWidth - inspectorWidth // Right side panel
	
	// Background
	inspectorBg := ebiten.NewImage(inspectorWidth, screenHeight)
	inspectorBg.Fill(color.RGBA{40, 40, 40, 200})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(inspectorX), 0)
	screen.DrawImage(inspectorBg, opts)
	
	// Title
	text.Draw(screen, "INSPECTOR", basicfont.Face7x13, inspectorX+10, 20, color.White)
	
	if ui.selectedEntity != nil {
		y := 50
		text.Draw(screen, fmt.Sprintf("Name: %s", ui.selectedEntity.Name), basicfont.Face7x13, inspectorX+10, y, color.White)
		y += 20
		text.Draw(screen, fmt.Sprintf("ID: %d", ui.selectedEntity.ID), basicfont.Face7x13, inspectorX+10, y, color.White)
		y += 20
		text.Draw(screen, fmt.Sprintf("X: %.1f", ui.selectedEntity.Position.X), basicfont.Face7x13, inspectorX+10, y, color.White)
		y += 20
		text.Draw(screen, fmt.Sprintf("Y: %.1f", ui.selectedEntity.Position.Y), basicfont.Face7x13, inspectorX+10, y, color.White)
		
		if ui.selectedEntity.Sprite != nil {
			y += 20
			w, h := ui.selectedEntity.Sprite.Bounds().Dx(), ui.selectedEntity.Sprite.Bounds().Dy()
			text.Draw(screen, fmt.Sprintf("Size: %dx%d", w, h), basicfont.Face7x13, inspectorX+10, y, color.White)
		}
		
		// Camera-relative position
		y += 30
		text.Draw(screen, "-- Camera View --", basicfont.Face7x13, inspectorX+10, y, color.RGBA{150, 150, 150, 255})
		y += 20
		screenX, screenY := cam.WorldToScreen(ui.selectedEntity.Position.X, ui.selectedEntity.Position.Y)
		text.Draw(screen, fmt.Sprintf("Screen: %.1f,%.1f", screenX, screenY), basicfont.Face7x13, inspectorX+10, y, color.RGBA{150, 150, 150, 255})
		
	} else {
		text.Draw(screen, "No entity selected", basicfont.Face7x13, inspectorX+10, 50, color.RGBA{128, 128, 128, 255})
		text.Draw(screen, "Click and drag", basicfont.Face7x13, inspectorX+10, 70, color.RGBA{128, 128, 128, 255})
		text.Draw(screen, "to select and move", basicfont.Face7x13, inspectorX+10, 90, color.RGBA{128, 128, 128, 255})
	}
	
	// Camera controls help
	y := screenHeight - 120
	text.Draw(screen, "-- Camera --", basicfont.Face7x13, inspectorX+10, y, color.RGBA{150, 150, 150, 255})
	y += 20
	text.Draw(screen, "WASD: Pan", basicfont.Face7x13, inspectorX+10, y, color.RGBA{120, 120, 120, 255})
	y += 15
	text.Draw(screen, "Wheel: Zoom", basicfont.Face7x13, inspectorX+10, y, color.RGBA{120, 120, 120, 255})
	y += 15
	text.Draw(screen, "R: Reset", basicfont.Face7x13, inspectorX+10, y, color.RGBA{120, 120, 120, 255})
}

// DrawLogPanel draws the log panel at the bottom
func (ui *EditorUI) DrawLogPanel(screen *ebiten.Image, screenWidth, screenHeight int) {
	logHeight := 120
	logY := screenHeight - logHeight
	logWidth := screenWidth
	if ui.inspectorOpen {
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
	start := len(ui.logMessages) - maxVisible
	if start < 0 {
		start = 0
	}
	
	for i, msg := range ui.logMessages[start:] {
		y := logY + 40 + i*12
		if y < logY+logHeight-10 {
			text.Draw(screen, msg, basicfont.Face7x13, 10, y, color.RGBA{200, 200, 200, 255})
		}
	}
}

// DrawGrid draws a grid in the background for editor mode
func (ui *EditorUI) DrawGrid(screen *ebiten.Image, cam *camera.Camera, viewportWidth, viewportHeight int) {
	gridSize := 50.0 // Grid cell size in world units
	gridColor := color.RGBA{60, 60, 70, 255}
	
	// Calculate grid lines to draw based on camera position and zoom
	startX := int((cam.X / gridSize)) - 1
	endX := int(((cam.X + float64(viewportWidth)/cam.Zoom) / gridSize)) + 1
	startY := int((cam.Y / gridSize)) - 1
	endY := int(((cam.Y + float64(viewportHeight)/cam.Zoom) / gridSize)) + 1
	
	// Draw vertical lines
	for x := startX; x <= endX; x++ {
		worldX := float64(x) * gridSize
		screenX, _ := cam.WorldToScreen(worldX, 0)
		
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
		_, screenY := cam.WorldToScreen(0, worldY)
		
		if screenY >= 0 && screenY <= float64(viewportHeight) {
			line := ebiten.NewImage(viewportWidth, 1)
			line.Fill(gridColor)
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(0, screenY)
			screen.DrawImage(line, opts)
		}
	}
}
