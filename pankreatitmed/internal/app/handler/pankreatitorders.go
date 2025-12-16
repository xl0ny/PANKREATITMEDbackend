package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pankreatitmed/internal/app/authctx"
	"pankreatitmed/internal/app/dto/request"
	"pankreatitmed/internal/app/mapper"
	"time"

	"github.com/gin-gonic/gin"
)

// PankreatitOrderFromCart godoc
// @Summary      Иконка корзины: черновик и количество позиций
// @Tags         orders
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} map[string]any
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      500 {object} map[string]any "internal error"
// @Router       /pankreatitorders/cart [get]
func (h *Handler) PankreatitOrderFromCart(c *gin.Context) {
	usr, check := authctx.Get(c)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "problem with your token"})
		return
	}
	mo, err := h.svcs.PankreatitOrders.GetDraft(usr.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mo)
}

// ListPankreatitOrders godoc
// @Summary      Список заявок (с фильтрацией по статусу и дате формирования)
// @Description  Создатель видит свои заявки; модератор — все.
// @Tags         orders
// @Security     BearerAuth
// @Produce      json
// @Param        status     query string false "Статус (draft|formed|completed|rejected)"
// @Param        from_date  query string false "Дата С (YYYY-MM-DD)"
// @Param        to_date    query string false "Дата ПО (YYYY-MM-DD)"
// @Success      200 {object} map[string]any
// @Failure      400 {object} map[string]any "bad request"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Router       /pankreatitorders [get]
func (h *Handler) ListPankreatitOrders(c *gin.Context) {
	usr, check := authctx.Get(c)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "problem with your token"})
		return
	}
	var filters request.GetPankreatitOrders
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svcs.PankreatitOrders.List(usr.ID, filters.Status, filters.FromDate, filters.ToDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, res)
}

// PankreatitOrderGet godoc
// @Summary      Получить одну заявку (с позициями)
// @Tags         orders
// @Security     BearerAuth
// @Produce      json
// @Param        id   path int true "ID заявки"
// @Success      200 {object} map[string]any
// @Failure      400 {object} map[string]any "bad request"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      404 {object} map[string]any "not found"
// @Router       /pankreatitorders/{id} [get]
func (h *Handler) PankreatitOrderGet(c *gin.Context) {
	var id request.GetPankreatitOrder
	usr, check := authctx.Get(c)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "problem with your token"})
		return
	}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svcs.PankreatitOrders.Get(id.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !usr.IsModerator && res.CreatorID != usr.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "u don't have permission to get this order(not your order)"})
		return
	}
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)

}

// TODO разобраться почему не выводит ошибку, когда не находит ордер по id-ку
// PankreatitOrderUpdate godoc
// @Summary      Обновить поля заявки (модератор)
// @Tags         orders
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path int true "ID заявки"
// @Param        input body request.UpdatePankreatitOrder true "Изменяемые поля"
// @Success      200 {string} string "ok"
// @Failure      400 {object} map[string]any "bad request"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      403 {object} map[string]any "forbidden"
// @Failure      404 {object} map[string]any "not found"
// @Router       /pankreatitorders/{id} [put]
func (h *Handler) PankreatitOrderUpdate(c *gin.Context) {
	usr, check := authctx.Get(c)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "problem with your token"})
		return
	}
	var id request.GetPankreatitOrder
	var mo request.UpdatePankreatitOrder
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&mo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o, err := h.svcs.PankreatitOrders.Get(id.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !usr.IsModerator && o.CreatorID != usr.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "u don't have permission to update this order(not your order)"})
		return
	}
	if err := h.svcs.PankreatitOrders.Update(id.ID, &mo); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.Status(http.StatusOK)
}

func (h *Handler) PankreatitOrderSetRanson(c *gin.Context) {
	var ranson request.PankreatitOrderSetRanson
	if err := c.ShouldBindJSON(&ranson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ranson.Key != "A9F3C47E2B8D1C6A" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "problem with access key"})
		return
	}
	order := mapper.PankreatitOrderSetRansonToUpdatePankreatitOrder(ranson)
	h.svcs.PankreatitOrders.Update(ranson.ID, &order)
}

// PankreatitOrderForm godoc
// @Summary      Сформировать заявку (создатель)
// @Description  Проверяет владельца; валидирует обязательные поля; устанавливает дату формирования
// @Tags         orders
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "ID заявки"
// @Success      200 {string} string "ok"
// @Failure      400 {object} map[string]any "bad request"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      403 {object} map[string]any "forbidden (not your order)"
// @Failure      404 {object} map[string]any "not found"
// @Failure      409 {object} map[string]any "MedOrderIsNotDraft"
// @Router       /pankreatitorders/{id}/form [put]
func (h *Handler) PankreatitOrderForm(c *gin.Context) {
	usr, check := authctx.Get(c)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "problem with your token"})
		return
	}
	var id request.GetPankreatitOrder
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o, err := h.svcs.PankreatitOrders.Get(id.ID)
	//fmt.Println(o.CreatorID, usr.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !usr.IsModerator && o.CreatorID != usr.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "u don't have permission to form this order(not your order)"})
		return
	}
	if err := h.svcs.PankreatitOrders.Form(id.ID); err != nil {
		if err.Error() == "MedOrderIsNotDraft" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// PankreatitOrderComplete godoc
// @Summary      Завершить/отклонить заявку (модератор)
// @Description  Меняет статус на complete/reject, рассчитывает поля, ставит moderator_id/finished_at
// @Tags         orders
// @Security     BearerAuth
// @Produce      json
// @Param        id     path int    true "ID заявки"
// @Param        status path string true "complete|reject"
// @Success      200 {string} string "ok"
// @Failure      400 {object} map[string]any "bad request"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      403 {object} map[string]any "forbidden"
// @Failure      409 {object} map[string]any "MedOrderIsNotFormed"
// @Router       /pankreatitorders/{id}/set/{status} [put]
func (h *Handler) PankreatitOrderComplete(c *gin.Context) {
	var idstatus request.EndOrCancelPankreatitOrder
	usr, check := authctx.Get(c)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "problem with your token"})
		return
	}
	if err := c.ShouldBindUri(&idstatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if idstatus.Status == "completed" {
		data, _ := json.Marshal(map[string]int{"id": int(idstatus.ID)})

		client := &http.Client{Timeout: 5 * time.Second}

		req, _ := http.NewRequest(
			http.MethodPost,
			"http://localhost:8000/",
			bytes.NewBuffer(data),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := client.Do(req)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		fmt.Println(resp.StatusCode)
		fmt.Println(string(body))
	} else if err := h.svcs.PankreatitOrders.CancelOrEnd(idstatus.ID, usr.ID, idstatus.Status); err != nil {
		if err.Error() == "MedOrderIsNotFormed" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.Status(http.StatusOK)
}

// PankreatitOrderDelete godoc
// @Summary      Удалить черновую заявку (создатель)
// @Description  Soft-delete: переводит заявку в статус deleted (только для draft)
// @Tags         orders
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "ID заявки"
// @Success      200 {string} string "ok"
// @Failure      400 {object} map[string]any "bad request"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      409 {object} map[string]any "not draft / conflict"
// @Router       /pankreatitorders/{id} [delete]
func (h *Handler) PankreatitOrderDelete(c *gin.Context) {
	var id request.GetPankreatitOrder
	usr, check := authctx.Get(c)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "problem with your token"})
		return
	}
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o, err := h.svcs.PankreatitOrders.Get(id.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !usr.IsModerator && o.CreatorID != usr.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "u don't have permission to delete this order(not your order)"})
		return
	}
	if err := h.svcs.PankreatitOrders.Delete(id.ID); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
