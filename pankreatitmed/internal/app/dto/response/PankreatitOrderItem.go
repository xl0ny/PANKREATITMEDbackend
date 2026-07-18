package response

import "pankreatitmed/internal/app/ds"

type SendPankreatitOrderItem struct {
	ID             uint         `json:"id"`
	CriterionID    uint         `json:"criterion_id"`
	Criterion      ds.Criterion `json:"criterion"`
	Position       int          `json:"position"`
	ValueNum       *float64     `json:"value_num"`
	ValueIndicator bool         `json:"value_indicator"`
}
