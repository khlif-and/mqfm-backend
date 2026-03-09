package response

type PreferenceResponse struct {
	PlaybackSpeed     float64 `json:"playback_speed"`
	SleepTimerMinutes int     `json:"sleep_timer_minutes"`
	AutoDownloadWifi  bool    `json:"auto_download_wifi"`
}
