package repository

import (
	"fmt"
	"pankreatitmed/internal/app/ds"
	"pankreatitmed/internal/app/dto/request"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) CountItems(orderID uint) (int64, error) {
	var cnt int64
	return cnt, r.db.Model(&ds.PankreatitOrderItem{}).Where("pankreatit_order_id = ?", orderID).Count(&cnt).Error
}

func (r *Repository) IsPankreatitOrderDeleted(orderID uint) (bool, error) {
	var o ds.PankreatitOrder
	err := r.db.First(&o, "id = ?", orderID).Error
	if err == nil {
		return o.Status == "deleted", nil
	}
	return true, err
}

func (r *Repository) IsPankreatitOrderDraft(orderID uint) (bool, error) {
	var o ds.PankreatitOrder
	err := r.db.First(&o, "id = ?", orderID).Error
	if err == nil {
		return o.Status == "draft", nil
	}
	return true, err
}

func (r *Repository) IsPankreatitOrderFormed(orderID uint) (bool, error) {
	var o ds.PankreatitOrder
	err := r.db.First(&o, "id = ?", orderID).Error
	if err == nil {
		return o.Status == "formed", nil
	}
	return true, err
}

// ----------------------------------------------------------

func (r *Repository) GetOrCreateDraftPankreatitOrder(creatorID uint) (*ds.PankreatitOrder, error) {
	var o ds.PankreatitOrder
	println("GetOrCreateDraftPankreatitOrder")
	if err := r.db.Where("creator_id = ? AND status = 'draft'", creatorID).First(&o).Error; err == nil {
		return &o, nil
	}
	o = ds.PankreatitOrder{Status: "draft", CreatorID: creatorID}
	return &o, r.db.Create(&o).Error
}

func (r *Repository) GetPankreatitOrders(userID uint, status *string, start, end *time.Time) ([]ds.PankreatitOrder, error) {
	var orders []ds.PankreatitOrder
	usr, err := r.GetMedUserByID(userID)
	if err != nil {
		return nil, err
	}

	q := r.db.Model(&ds.PankreatitOrder{})

	if !usr.IsModerator {
		q = q.Where("creator_id = ?", userID)
	}

	if status != nil && *status != "" {
		q = q.Where("status = ?", *status)
	}

	var endInclusive *time.Time

	if end != nil {
		tmp := end.AddDate(0, 0, 1)
		endInclusive = &tmp
	} else {
		endInclusive = nil
	}
	switch {
	case start != nil && end != nil:
		q = q.Where("created_at >= ? AND created_at < ?", *start, endInclusive)
	case start != nil:
		q = q.Where("created_at >= ?", *start)
	case end != nil:
		q = q.Where("created_at < ?", endInclusive)
	}

	if err := q.Order("created_at ASC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *Repository) GetPankreatitOrderWithItems(orderID uint) (ds.PankreatitOrder, []ds.PankreatitOrderItem, error) {
	var o ds.PankreatitOrder
	if err := r.db.First(&o, orderID).Error; err != nil {
		return ds.PankreatitOrder{}, nil, err
	}
	var items []ds.PankreatitOrderItem
	if err := r.db.Preload("Criterion").Where("pankreatit_order_id = ?", orderID).Order("id").Find(&items).Error; err != nil {
		return ds.PankreatitOrder{}, nil, err
	}
	fmt.Println(items)
	return o, items, nil
}

func (r *Repository) UpdatePankreatitOrder(id uint, order *request.UpdatePankreatitOrder) error {
	tx := r.db.Model(&ds.PankreatitOrder{}).Where("id = ?", id).Updates(order)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrInvalidData
	}
	return nil
}

func (r *Repository) FormPankreatitOrder(id uint) error {
	tx := r.db.Model(&ds.PankreatitOrder{}).
		Where("id = ?", id).
		UpdateColumns(map[string]any{
			"status":    "formed",
			"formed_at": time.Now(),
		})

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) EndOrCancelPankreatitOrder(id, moderator uint, status string) error {
	tx := r.db.Model(&ds.PankreatitOrder{}).Where("id = ?", id).UpdateColumns(map[string]any{
		"status":       status,
		"finished_at":  time.Now(),
		"moderator_id": moderator,
	})

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) SoftDeleteOrderSQL(orderID uint) error {
	sql := `UPDATE pankreatitorders SET status='deleted', formed_at = NOW() WHERE id=$1 AND status = 'draft'`
	tx := r.db.Exec(sql, orderID)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return fmt.Errorf("pankreatitorder %d not updated because not draft or not exists", orderID)
	}
	return nil
}

func (r *Repository) SetRansonAndRisk(orderID uint, score int, risk string) error {
	return r.db.Model(&ds.PankreatitOrder{}).
		Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"ranson_score":   score,
			"mortality_risk": risk,
		}).Error
}
