package tinybird

import (
	"testing"

	"github.com/jonashiltl/openchangelog/internal/analytics"
)

func TestNew(t *testing.T) {
	bird := New(TinybirdOptions{}).(*bird)

	if len(bird.buffer) != 0 {
		t.Errorf("Expected buffer to be empty, got %d", len(bird.buffer))
	}
}

func TestEmit(t *testing.T) {
	bird := New(TinybirdOptions{}).(*bird)
	bird.Emit(analytics.Event{})

	if len(bird.buffer) != 1 {
		t.Errorf("Expected buffer length to be %d, got %d", 1, len(bird.buffer))
	}
}

func TestFlush(t *testing.T) {
	bird := New(TinybirdOptions{}).(*bird)
	bird.Emit(analytics.Event{})
	bird.flush()

	if len(bird.buffer) != 0 {
		t.Errorf("Expected buffer to be empty, got %d", len(bird.buffer))
	}
}
