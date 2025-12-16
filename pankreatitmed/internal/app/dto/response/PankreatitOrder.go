package response

import (
	"time"
)

type SendCartPankreatitOrder struct {
	PankreatitOrderId uint `json:"pankreatit_order_id"`
	CriteriaAmount    uint `json:"criteria_amount"`
}

type SendPankreatitOrders struct {
	ID            uint       `json:"id"`
	Status        string     `json:"status"`
	CreatorID     uint       `json:"creator_id"`
	FormedAt      *time.Time `json:"formed_at"`
	FinishedAt    *time.Time `json:"finished_at"`
	RansonScore   *int       `json:"ranson_score"`
	MortalityRisk *string    `json:"mortality_risk"`
}

type SendPankreatitOrder struct {
	ID            uint                      `json:"id"`
	Status        string                    `json:"status"`
	CreatorID     uint                      `json:"creator_id"`
	FormedAt      *time.Time                `json:"formed_at"`
	FinishedAt    *time.Time                `json:"finished_at"`
	ModeratorID   *uint                     `json:"moderator_id"`
	RansonScore   *int                      `json:"ranson_score"`
	MortalityRisk *string                   `json:"mortality_risk"`
	Items         []SendPankreatitOrderItem `json:"criteria"`
}
