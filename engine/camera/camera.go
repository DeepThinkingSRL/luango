package camera

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

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

func (c *Camera) Reset() {
	c.X = 0
	c.Y = 0
	c.Zoom = 1.0
}

func (c *Camera) Move(deltaX, deltaY float64) {
	c.X += deltaX
	c.Y += deltaY
}

func (c *Camera) SetZoom(zoom float64) {
	if zoom < 0.5 {
		zoom = 0.5
	}
	if zoom > 3.0 {
		zoom = 3.0
	}
	c.Zoom = zoom
}

func (c *Camera) ZoomIn() {
	c.SetZoom(c.Zoom * 1.1)
}

func (c *Camera) ZoomOut() {
	c.SetZoom(c.Zoom / 1.1)
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *Camera) ScreenToWorld(screenX, screenY float64) (float64, float64) {
	worldX := (screenX / c.Zoom) + c.X
	worldY := (screenY / c.Zoom) + c.Y
	return worldX, worldY
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *Camera) WorldToScreen(worldX, worldY float64) (float64, float64) {
	screenX := (worldX - c.X) * c.Zoom
	screenY := (worldY - c.Y) * c.Zoom
	return screenX, screenY
}

// CreateTransformMatrix creates a camera transformation matrix for rendering
func (c *Camera) CreateTransformMatrix() ebiten.GeoM {
	var matrix ebiten.GeoM
	matrix.Scale(c.Zoom, c.Zoom)
	matrix.Translate(-c.X*c.Zoom, -c.Y*c.Zoom)
	return matrix
}

// Update handles camera controls using an input manager
func (c *Camera) Update(inputManager interface{}) {
	// Calculate move speed based on zoom (slower when zoomed in)
	moveSpeed := 5.0 / c.Zoom
	
	// Handle movement using direct key checks for now
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		c.X -= moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		c.X += moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		c.Y -= moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		c.Y += moveSpeed
	}
	
	// Handle zoom
	_, wheelY := ebiten.Wheel()
	if wheelY > 0 || ebiten.IsKeyPressed(ebiten.KeyEqual) || ebiten.IsKeyPressed(ebiten.KeyKPAdd) {
		c.Zoom *= 1.1
		if c.Zoom > 3.0 {
			c.Zoom = 3.0
		}
	}
	if wheelY < 0 || ebiten.IsKeyPressed(ebiten.KeyMinus) || ebiten.IsKeyPressed(ebiten.KeyKPSubtract) {
		c.Zoom /= 1.1
		if c.Zoom < 0.5 {
			c.Zoom = 0.5
		}
	}
	
	// Reset camera
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		c.Reset()
		time.Sleep(200 * time.Millisecond)
	}
}
