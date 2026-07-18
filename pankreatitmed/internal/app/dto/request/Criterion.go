package request

type GetCriterion struct {
	ID   uint   `uri:"id"`
	Code string `uri:"code" binding:"omitempty"`
}

type GetCriteria struct {
	Query string `form:"query"`
}

type CreateCriterion struct {
	Code        *string  `json:"code" binding:"required"`
	Name        *string  `json:"name" binding:"required"`
	Description *string  `json:"description" binding:"required"`
	Duration    *string  `json:"duration" binding:"required"`
	HomeVisit   *bool    `json:"home_visit"`
	Status      *string  `json:"status" binding:"omitempty,oneof=active deleted"`
	Unit        *string  `json:"unit"`
	RefLow      *float64 `json:"ref_low"`
	RefHigh     *float64 `json:"ref_high"`
}

type UpdateCriterion struct {
	Code        *string  `json:"code"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Duration    *string  `json:"duration"`
	ImageURL    *string  `json:"image_url"`
	HomeVisit   *bool    `json:"home_visit"`
	Unit        *string  `json:"unit"`
	RefLow      *float64 `json:"ref_low"`
	RefHigh     *float64 `json:"ref_high"`
}

type CreateCriterionIamage struct {
	ID uint `json:"id"`
}
