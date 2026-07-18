package response

type SendCriterion struct {
	ID          uint     `json:"id" binding:"required"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Duration    string   `json:"duration"`
	HomeVisit   bool     `json:"home_visit"`
	ImageURL    *string  `json:"image_url"`
	Status      string   `json:"status"`
	Unit        string   `json:"unit"`
	RefLow      *float64 `json:"ref_low"`
	RefHigh     *float64 `json:"ref_high"`
}
