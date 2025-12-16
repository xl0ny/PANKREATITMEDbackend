package mapper

import (
	"pankreatitmed/internal/app/ds"
	"pankreatitmed/internal/app/dto/request"
	"pankreatitmed/internal/app/dto/response"
)

func PankreatitOrderToSendPankreatitOrder(o *ds.PankreatitOrder, amount uint) response.SendCartPankreatitOrder {
	return response.SendCartPankreatitOrder{
		PankreatitOrderId: o.ID,
		CriteriaAmount:    amount,
	}
}

func PankreatitOrderToSendPankreatitOrders(mo ds.PankreatitOrder) response.SendPankreatitOrders {
	return response.SendPankreatitOrders{
		ID:            mo.ID,
		Status:        mo.Status,
		CreatorID:     mo.CreatorID,
		FormedAt:      mo.FormedAt,
		FinishedAt:    mo.FinishedAt,
		RansonScore:   mo.RansonScore,
		MortalityRisk: mo.MortalityRisk,
	}
}

func PankreatitOrdersToSendPankreatitOrders(mos []ds.PankreatitOrder) []response.SendPankreatitOrders {
	list := make([]response.SendPankreatitOrders, len(mos))
	for i, c := range mos {
		list[i] = PankreatitOrderToSendPankreatitOrders(c)
	}
	return list
}

func PankreatitOrderToSendPankreatitOrderWithItems(mo ds.PankreatitOrder, itms []ds.PankreatitOrderItem) response.SendPankreatitOrder {
	return response.SendPankreatitOrder{
		ID:            mo.ID,
		Status:        mo.Status,
		CreatorID:     mo.CreatorID,
		FormedAt:      mo.FormedAt,
		FinishedAt:    mo.FinishedAt,
		ModeratorID:   mo.ModeratorID,
		RansonScore:   mo.RansonScore,
		MortalityRisk: mo.MortalityRisk,
		Items:         PankreatitOrderItemsToSendPankreatitOrderItems(itms),
	}
}

func PankreatitOrderSetRansonToUpdatePankreatitOrder(order request.PankreatitOrderSetRanson) request.UpdatePankreatitOrder {
	return request.UpdatePankreatitOrder{
		Status:        order.Status,
		FinishedAt:    order.FinishedAt,
		RansonScore:   order.RansonScore,
		MortalityRisk: order.MortalityRisk,
	}
}
