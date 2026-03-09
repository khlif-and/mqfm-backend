package response

type StatsRecapResponse struct {
	WeeklyMinutes  int                  `json:"weekly_minutes"`
	MonthlyMinutes int                  `json:"monthly_minutes"`
	TopCategories  []CategoryStatResponse `json:"top_categories"`
	TopArtists     []ArtistStatResponse   `json:"top_artists"`
	DailyStats     []DailyStatResponse    `json:"daily_stats"`
}

type CategoryStatResponse struct {
	CategoryID uint   `json:"category_id"`
	Name       string `json:"name"`
	TotalTime  int    `json:"total_time_seconds"`
}

type ArtistStatResponse struct {
	Artist    string `json:"artist"`
	TotalTime int    `json:"total_time_seconds"`
}

type DailyStatResponse struct {
	Date      string `json:"date"`
	TotalTime int    `json:"total_time_seconds"`
}
