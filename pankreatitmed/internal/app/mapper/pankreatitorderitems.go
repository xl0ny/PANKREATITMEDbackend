package mapper

import (
	"pankreatitmed/internal/app/ds"
	"pankreatitmed/internal/app/dto/response"
)

func PankreatitOrderItemToSendPankreatitOrderItem(item ds.PankreatitOrderItem) response.SendPankreatitOrderItem {
	return response.SendPankreatitOrderItem{
		ID:             item.ID,
		CriterionID:    item.CriterionID,
		Criterion:      item.Criterion,
		Position:       item.Position,
		ValueNum:       item.ValueNum,
		ValueIndicator: item.ValueIndicator,
	}
}

func PankreatitOrderItemsToSendPankreatitOrderItems(items []ds.PankreatitOrderItem) []response.SendPankreatitOrderItem {
	list := make([]response.SendPankreatitOrderItem, len(items))
	for i, c := range items {
		list[i] = PankreatitOrderItemToSendPankreatitOrderItem(c)
	}
	return list
}
