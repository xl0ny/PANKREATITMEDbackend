package handler

import (
	"net/http"
	"pankreatitmed/internal/app/authctx"
	"pankreatitmed/internal/app/dto/request"

	"github.com/gin-gonic/gin"
)

// DeletePankreatitOrderItem godoc
// @Summary      Удалить услугу из заявки (м-м) без PK м-м
// @Tags         order-items
// @Security     BearerAuth
// @Produce      json
// @Param        pankreatit_order_id  query int true "ID заявки"
// @Param        criterion_id         query int true "ID услуги"
// @Success      200 {object} map[string]any "status: ok"
// @Failure      400 {object} map[string]any "bad request"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      404 {object} map[string]any "not found"
// @Router       /pankreatitorders/items [delete]
func (h *Handler) DeletePankreatitOrderItem(c *gin.Context) {
	var item request.GetPankreatitOrderItem
	usr, check := authctx.Get(c)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "problem with your token"})
		return
	}
	if err := c.ShouldBindQuery(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	o, err := h.svcs.PankreatitOrders.Get(item.PankreatitOrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !usr.IsModerator && o.CreatorID != usr.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "u don't have permission to delete this order(not your order)"})
		return
	}
	if err := h.svcs.PankreatitOrderItems.Delete(item.PankreatitOrderID, item.CriterionID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// UpdatePankreatitOrderItem godoc
// @Summary      Изменить поля м-м (кол-во/порядок/значение) без PK м-м
// @Tags         order-items
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        pankreatit_order_id  query int  true  "ID заявки"
// @Param        criterion_id         query int  true  "ID услуги"
// @Param        input                body  request.PankreatitOrderItemUpdate true "Поле(я) для обновления"
// @Success      200 {object} map[string]any "status: ok"
// @Failure      400 {object} map[string]any "bad request"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      403 {object} map[string]any "forbidden"
// @Failure      404 {object} map[string]any "not found"
// @Router       /pankreatitorders/items [put]
func (h *Handler) UpdatePankreatitOrderItem(c *gin.Context) {
	var item request.GetPankreatitOrderItem
	var fields request.PankreatitOrderItemUpdate
	usr, check := authctx.Get(c)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "problem with your token"})
		return
	}
	if err := c.ShouldBindQuery(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&fields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o, err := h.svcs.PankreatitOrders.Get(item.PankreatitOrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !usr.IsModerator && o.CreatorID != usr.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "u don't have permission to delete this order(not your order)"})
		return
	}
	if err := h.svcs.PankreatitOrderItems.Update(item.PankreatitOrderID, item.CriterionID, fields.Position, fields.ValueNum); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
