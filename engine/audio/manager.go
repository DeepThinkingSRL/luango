package audio

import (
"fmt"
"os"
"time"

"github.com/faiface/beep/speaker"
"github.com/faiface/beep/wav"
)

type Manager struct {
	initialized bool
}

func NewManager() *Manager {
	return &Manager{
		initialized: false,
	}
}

func (am *Manager) Initialize() error {
	am.initialized = true
	return nil
}

func (am *Manager) PlaySound(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open audio file: %w", err)
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode audio file: %w", err)
	}
	defer streamer.Close()

	if !am.initialized {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		am.initialized = true
	}

	speaker.Play(streamer)
	fmt.Printf("[Audio] Playing sound: %s\n", path)
	return nil
}

func (am *Manager) IsInitialized() bool {
	return am.initialized
}
