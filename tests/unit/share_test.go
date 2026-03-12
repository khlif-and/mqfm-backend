package service_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"mqfm-backend/internal/app/service"
)

// ──────────────────────── GenerateAudioShareLink ────────────────────────

func TestShare_GenerateAudioShareLink(t *testing.T) {
	svc := service.NewShareService("https://mqfm.app")
	link := svc.GenerateAudioShareLink(42)
	assert.Equal(t, "https://mqfm.app/audio/42", link)
}

func TestShare_GenerateAudioShareLink_DifferentBase(t *testing.T) {
	svc := service.NewShareService("https://example.com")
	link := svc.GenerateAudioShareLink(1)
	assert.Equal(t, "https://example.com/audio/1", link)
}

// ──────────────────────── GenerateClipShareLink ────────────────────────

func TestShare_GenerateClipShareLink(t *testing.T) {
	svc := service.NewShareService("https://mqfm.app")
	link := svc.GenerateClipShareLink("abc123token")
	assert.Equal(t, "https://mqfm.app/clip/abc123token", link)
}

func TestShare_GenerateClipShareLink_EmptyToken(t *testing.T) {
	svc := service.NewShareService("https://mqfm.app")
	link := svc.GenerateClipShareLink("")
	assert.Equal(t, "https://mqfm.app/clip/", link)
}

// ──────────────────────── GeneratePlaylistShareLink ────────────────────────

func TestShare_GeneratePlaylistShareLink(t *testing.T) {
	svc := service.NewShareService("https://mqfm.app")
	link := svc.GeneratePlaylistShareLink("playlist-xyz")
	assert.Equal(t, "https://mqfm.app/playlist/playlist-xyz", link)
}

// ──────────────────────── Format consistency ────────────────────────

func TestShare_AllLinks_ConsistentFormat(t *testing.T) {
	base := "https://mqfm.app"
	svc := service.NewShareService(base)

	audioLink := svc.GenerateAudioShareLink(10)
	clipLink := svc.GenerateClipShareLink("token")
	playlistLink := svc.GeneratePlaylistShareLink("ptoken")

	assert.Equal(t, fmt.Sprintf("%s/audio/10", base), audioLink)
	assert.Equal(t, fmt.Sprintf("%s/clip/token", base), clipLink)
	assert.Equal(t, fmt.Sprintf("%s/playlist/ptoken", base), playlistLink)
}
