package request

import (
	"time"
)

type GetPankreatitOrders struct {
	Status   *string    `form:"status" binding:"omitempty,oneof=formed completed rejected draft deleted"`
	FromDate *time.Time `form:"from_date" time_format:"2006-01-02T15:04:05"`
	ToDate   *time.Time `form:"to_date" time_format:"2006-01-02T15:04:05"`
}

type GetPankreatitOrder struct {
	ID uint `uri:"id" binding:"required"`
}

type UpdatePankreatitOrder struct {
	Status        string     `json:"status"`
	FinishedAt    *time.Time `json:"finished_at" time_format:"2006-01-02T15:04:05"`
	RansonScore   int        `json:"ranson_score"`
	MortalityRisk string     `json:"mortality_risk"`
}

type FormPankreatitOrder struct {
	ID uint `json:"id" binding:"required"`
	//CreatorID uint   `json:"creator_id" binding:"required"`
	//Password  string `json:"password" binding:"required"`
}

type EndOrCancelPankreatitOrder struct {
	ID     uint   `uri:"id" binding:"required"`
	Status string `uri:"status" binding:"required,oneof=completed rejected"`
}

type GetModerator struct {
	ModeratorID uint   `form:"moderator_id" binding:"required"`
	Password    string `form:"password" binding:"required"`
}

type DeletePankreatitOrder struct {
	ID uint `json:"id" binding:"required"`
}

type PankreatitOrderSetRanson struct {
	ID            uint       `json:"id" binding:"required"`
	Status        string     `json:"status"`
	FinishedAt    *time.Time `json:"finished_at" time_format:"2006-01-02T15:04:05"`
	RansonScore   int        `json:"ranson_score" binding:"omitempty"`
	MortalityRisk string     `json:"mortality_risk" binding:"omitempty"`
	Key           string     `json:"key" binding:"required"`
}
