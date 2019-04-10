package challenge

type match struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Score     float32 `json:"score"`
}

type suggestions struct {
	Suggestions []match `json:"suggestions"`
}