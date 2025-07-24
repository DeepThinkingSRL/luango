package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Manager struct {
	keyStates map[ebiten.Key]bool
	lastKeyStates map[ebiten.Key]bool
	keyPressTimes map[ebiten.Key]time.Time
	debounceTime time.Duration
}

func NewManager() *Manager {
	return &Manager{
		keyStates: make(map[ebiten.Key]bool),
		lastKeyStates: make(map[ebiten.Key]bool),
		keyPressTimes: make(map[ebiten.Key]time.Time),
		debounceTime: 200 * time.Millisecond,
	}
}

func (m *Manager) Update() {
	// Update key states
	keys := []ebiten.Key{
		ebiten.KeyF1, ebiten.KeyF2, ebiten.KeyF3, ebiten.KeyF11,
		ebiten.KeyArrowUp, ebiten.KeyArrowDown, ebiten.KeyArrowLeft, ebiten.KeyArrowRight,
		ebiten.KeyW, ebiten.KeyA, ebiten.KeyS, ebiten.KeyD,
		ebiten.KeyR, ebiten.KeyEqual, ebiten.KeyMinus,
		ebiten.KeyKPAdd, ebiten.KeyKPSubtract,
		ebiten.KeySpace,
	}
	
	for _, key := range keys {
		m.lastKeyStates[key] = m.keyStates[key]
		m.keyStates[key] = ebiten.IsKeyPressed(key)
	}
}

func (m *Manager) IsKeyPressed(key ebiten.Key) bool {
	return m.keyStates[key]
}

func (m *Manager) IsKeyJustPressed(key ebiten.Key) bool {
	// Check if enough time has passed since last press (debouncing)
	if lastPressTime, exists := m.keyPressTimes[key]; exists {
		if time.Since(lastPressTime) < m.debounceTime {
			return false
		}
	}
	
	// Check if key was just pressed
	if m.keyStates[key] && !m.lastKeyStates[key] {
		m.keyPressTimes[key] = time.Now()
		return true
	}
	
	return false
}

func (m *Manager) IsKeyReleased(key ebiten.Key) bool {
	return !m.keyStates[key] && m.lastKeyStates[key]
}

func (m *Manager) GetMousePosition() (int, int) {
	return ebiten.CursorPosition()
}

func (m *Manager) IsMouseButtonPressed(button ebiten.MouseButton) bool {
	return ebiten.IsMouseButtonPressed(button)
}

func (m *Manager) GetWheelDelta() (float64, float64) {
	return ebiten.Wheel()
}
