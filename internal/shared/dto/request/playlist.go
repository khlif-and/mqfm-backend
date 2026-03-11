package request

type CreatePlaylistRequest struct {
	Name string `form:"name" binding:"required"`
}

type CreatePlaylistFromAudioRequest struct {
	Name    string `json:"name" binding:"required"`
	AudioID uint   `json:"audio_id" binding:"required"`
}

type UpdatePlaylistRequest struct {
	Name string `form:"name"`
}

type PlaylistAudioRequest struct {
	PlaylistID uint `json:"playlist_id" binding:"required"`
	AudioID    uint `json:"audio_id" binding:"required"`
}
