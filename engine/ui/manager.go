package ui

import (
	"fmt"
	"image/color"
	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

// DrawGrid draws a grid in the background for editor mode
func (m *Manager) DrawGrid(screen *ebiten.Image, camera interface{}, viewportWidth, viewportHeight int) {
	// Cast camera to get access to its fields
	cam, ok := camera.(*struct{ X, Y, Zoom float64 })
	if !ok {
		return // Skip if camera type is not what we expect
	}
	
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
		screenX := (worldX - cam.X) * cam.Zoom
		
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
		screenY := (worldY - cam.Y) * cam.Zoom
		
		if screenY >= 0 && screenY <= float64(viewportHeight) {
			line := ebiten.NewImage(viewportWidth, 1)
			line.Fill(gridColor)
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(0, screenY)
			screen.DrawImage(line, opts)
		}
	}
}

// DrawSelectionBorder draws a border around the selected entity
func (m *Manager) DrawSelectionBorder(screen *ebiten.Image, entity interface{}, camera interface{}) {
	// This is a simplified version - in a real implementation,
	// we'd properly type the entity and camera parameters
	// For now, we'll skip the actual drawing
}

// DrawInspector draws the inspector panel on the right side
func (m *Manager) DrawInspector(screen *ebiten.Image, selectedEntity interface{}, camera interface{}, dragEntity interface{}, screenWidth, screenHeight int) {
	inspectorWidth := 200
	inspectorX := screenWidth - inspectorWidth
	
	// Background
	inspectorBg := ebiten.NewImage(inspectorWidth, screenHeight)
	inspectorBg.Fill(color.RGBA{40, 40, 40, 200})
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(inspectorX), 0)
	screen.DrawImage(inspectorBg, opts)
	
	// Title
	text.Draw(screen, "INSPECTOR", basicfont.Face7x13, inspectorX+10, 20, color.White)
	
	if selectedEntity != nil {
		y := 50
		text.Draw(screen, "Entity selected", basicfont.Face7x13, inspectorX+10, y, color.White)
		
		// Show dragging status
		if dragEntity != nil {
			y += 30
			text.Draw(screen, "DRAGGING", basicfont.Face7x13, inspectorX+10, y, color.RGBA{255, 255, 0, 255})
		}
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
func (m *Manager) DrawLogPanel(screen *ebiten.Image, logMessages []string, inspectorOpen bool, screenWidth, screenHeight int) {
	logHeight := 120
	logY := screenHeight - logHeight
	logWidth := screenWidth
	if inspectorOpen {
		logWidth -= 200
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
	start := len(logMessages) - maxVisible
	if start < 0 {
		start = 0
	}
	
	for i, msg := range logMessages[start:] {
		y := logY + 40 + i*12
		if y < logY+logHeight-10 {
			text.Draw(screen, msg, basicfont.Face7x13, 10, y, color.RGBA{200, 200, 200, 255})
		}
	}
}

// DrawControls draws control information
func (m *Manager) DrawControls(screen *ebiten.Image, editorMode bool) {
	controlY := 40
	text.Draw(screen, "F1: Mode F2: Inspector F11: Fullscreen", basicfont.Face7x13, 10, controlY, color.RGBA{128, 128, 128, 255})
	if editorMode {
		text.Draw(screen, "WASD/Arrows: Pan  Wheel/+/-: Zoom  R: Reset", basicfont.Face7x13, 10, controlY+15, color.RGBA{128, 128, 128, 255})
		text.Draw(screen, "Drag: Move entity  Middle: Pan camera", basicfont.Face7x13, 10, controlY+30, color.RGBA{128, 128, 128, 255})
	}
}

// DrawDebugInfo draws debug information
func (m *Manager) DrawDebugInfo(screen *ebiten.Image, player interface{}, frame, entityCount, screenWidth, screenHeight int) {
	debugY := 100
	text.Draw(screen, "Luango Engine - DEBUG", basicfont.Face7x13, 10, debugY, color.White)
	text.Draw(screen, fmt.Sprintf("Frame: %d", frame), basicfont.Face7x13, 10, debugY+20, color.White)
	text.Draw(screen, fmt.Sprintf("Entities: %d", entityCount), basicfont.Face7x13, 10, debugY+40, color.White)
	text.Draw(screen, fmt.Sprintf("Screen: %dx%d", screenWidth, screenHeight), basicfont.Face7x13, 10, debugY+60, color.White)
	text.Draw(screen, "F3: Toggle Debug", basicfont.Face7x13, 10, screenHeight-20, color.RGBA{128, 128, 128, 255})
}
