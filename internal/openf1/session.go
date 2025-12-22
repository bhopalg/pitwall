package openf1

// OpenF1 returns times as strings in JSON; we'll decode to time.Time in Lesson 1.4.
// For now keep them as string so it’s easy.
type Session struct {
	SessionKey  int    `json:"session_key"`
	SessionName string `json:"session_name"`
	DateStart   string `json:"date_start"`
	DateEnd     string `json:"date_end"`
	Location    string `json:"location"`
	CountryName string `json:"country_name"`
	CircuitName string `json:"circuit_short_name"`
	MeetingKey  int    `json:"meeting_key"`
	Year        int    `json:"year"`
}
