package resources

import (
	"fmt"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Manager struct {
	sprites map[string]*ebiten.Image
}

func NewManager() *Manager {
	return &Manager{
		sprites: make(map[string]*ebiten.Image),
	}
}

func (rm *Manager) LoadSprite(path string) (*ebiten.Image, error) {
	// Check if already loaded
	if sprite, exists := rm.sprites[path]; exists {
		return sprite, nil
	}

	// Load the sprite
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sprite file %s: %w", path, err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode sprite file %s: %w", path, err)
	}

	sprite := ebiten.NewImageFromImage(img)
	rm.sprites[path] = sprite

	fmt.Printf("[Resources] Loaded sprite: %s\n", path)
	return sprite, nil
}

func (rm *Manager) GetSprite(path string) (*ebiten.Image, bool) {
	sprite, exists := rm.sprites[path]
	return sprite, exists
}

func (rm *Manager) UnloadSprite(path string) {
	delete(rm.sprites, path)
}

func (rm *Manager) UnloadAll() {
	rm.sprites = make(map[string]*ebiten.Image)
}

func (rm *Manager) GetLoadedSprites() []string {
	paths := make([]string, 0, len(rm.sprites))
	for path := range rm.sprites {
		paths = append(paths, path)
	}
	return paths
}
