package request

type GetPankreatitOrderItem struct {
	PankreatitOrderID uint `form:"pankreatit_order_id" binding:"required"`
	CriterionID       uint `form:"criterion_id" binding:"required"`
}

type PankreatitOrderItemDelete struct {
	PankreatitOrderID uint `json:"pankreatit_order_id" binding:"required"`
	CriterionID       uint `json:"criterion_id" binding:"required"`
}

type PankreatitOrderItemUpdate struct {
	Position *uint    `json:"position" binding:"omitempty"`
	ValueNum *float64 `json:"value_num" binding:"omitempty"`
}
