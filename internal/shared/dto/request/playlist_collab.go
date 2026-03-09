package request

type AddCollaboratorRequest struct {
	PlaylistID     uint `json:"playlist_id" binding:"required"`
	CollaboratorID uint `json:"collaborator_id" binding:"required"`
}

type ContributeAudioRequest struct {
	PlaylistID uint `json:"playlist_id" binding:"required"`
	AudioID    uint `json:"audio_id" binding:"required"`
}

type JoinPlaylistRequest struct {
	ShareToken string `json:"share_token" binding:"required"`
}
