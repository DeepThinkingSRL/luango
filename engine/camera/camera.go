package camera

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	X, Y float64  // Camera position
	Zoom float64  // Zoom level (1.0 = normal)
}

func New() *Camera {
	return &Camera{
		X:    0,
		Y:    0,
		Zoom: 1.0,
	}
}

func NewCamera() Camera {
	return Camera{
		X:    0,
		Y:    0,
		Zoom: 1.0,
	}
}

func (c *Camera) Reset() {
	c.X = 0
	c.Y = 0
	c.Zoom = 1.0
}

func (c *Camera) Move(dx, dy float64) {
	c.X += dx
	c.Y += dy
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

func (c *Camera) ZoomBy(factor float64) {
	c.SetZoom(c.Zoom * factor)
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

// GetTransformMatrix returns the camera transform matrix for rendering
func (c *Camera) GetTransformMatrix() ebiten.GeoM {
	var matrix ebiten.GeoM
	matrix.Scale(c.Zoom, c.Zoom)
	matrix.Translate(-c.X*c.Zoom, -c.Y*c.Zoom)
	return matrix
}

// FollowTarget smoothly moves camera to follow a target position
func (c *Camera) FollowTarget(targetX, targetY float64, screenWidth, screenHeight int, lerpFactor float64) {
	// Calculate desired camera position (center target on screen)
	desiredX := targetX - float64(screenWidth)/(2*c.Zoom)
	desiredY := targetY - float64(screenHeight)/(2*c.Zoom)
	
	// Smooth interpolation
	c.X += (desiredX - c.X) * lerpFactor
	c.Y += (desiredY - c.Y) * lerpFactor
}

func (c *Camera) String() string {
	return fmt.Sprintf("Camera(%.1f, %.1f, %.2fx)", c.X, c.Y, c.Zoom)
}
