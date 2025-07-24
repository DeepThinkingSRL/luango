package input

import (
"github.com/hajimehoshi/ebiten/v2"
)

type Manager struct {
	keyStates     map[ebiten.Key]bool
	lastKeyStates map[ebiten.Key]bool
	mouseX        int
	mouseY        int
	lastMouseX    int
	lastMouseY    int
}

func NewManager() *Manager {
	return &Manager{
		keyStates:     make(map[ebiten.Key]bool),
		lastKeyStates: make(map[ebiten.Key]bool),
	}
}

func (im *Manager) Initialize() {
	keys := []ebiten.Key{
		ebiten.KeyF1, ebiten.KeyF2, ebiten.KeyF3, ebiten.KeyF11,
		ebiten.KeyArrowUp, ebiten.KeyArrowDown, ebiten.KeyArrowLeft, ebiten.KeyArrowRight,
		ebiten.KeyW, ebiten.KeyA, ebiten.KeyS, ebiten.KeyD,
		ebiten.KeyR, ebiten.KeyEqual, ebiten.KeyMinus,
		ebiten.KeyKPAdd, ebiten.KeyKPSubtract,
		ebiten.KeySpace,
	}
	
	for _, key := range keys {
		im.keyStates[key] = false
		im.lastKeyStates[key] = false
	}
}

func (im *Manager) Update() {
	for key := range im.keyStates {
		im.lastKeyStates[key] = im.keyStates[key]
		im.keyStates[key] = ebiten.IsKeyPressed(key)
	}
	
	im.lastMouseX = im.mouseX
	im.lastMouseY = im.mouseY
	im.mouseX, im.mouseY = ebiten.CursorPosition()
}

func (im *Manager) IsKeyJustPressed(key ebiten.Key) bool {
	return im.keyStates[key] && !im.lastKeyStates[key]
}

func (im *Manager) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

func (im *Manager) GetMousePosition() (int, int) {
	return im.mouseX, im.mouseY
}

func (im *Manager) GetMouseDelta() (int, int) {
	return im.mouseX - im.lastMouseX, im.mouseY - im.lastMouseY
}

func (im *Manager) IsMouseButtonPressed(button ebiten.MouseButton) bool {
	return ebiten.IsMouseButtonPressed(button)
}

func (im *Manager) GetWheelDelta() (float64, float64) {
	return ebiten.Wheel()
}

func KeyFromString(k string) ebiten.Key {
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
	case "F1":
		return ebiten.KeyF1
	case "F2":
		return ebiten.KeyF2
	case "F3":
		return ebiten.KeyF3
	case "F11":
		return ebiten.KeyF11
	case "R":
		return ebiten.KeyR
	default:
		return 0
	}
}
