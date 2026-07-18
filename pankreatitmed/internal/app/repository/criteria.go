package repository

import (
	"errors"
	"pankreatitmed/internal/app/ds"
	"pankreatitmed/internal/app/dto/request"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *Repository) GetCriteria(q string) ([]ds.Criterion, error) {
	var list []ds.Criterion
	db := r.db.Model(&ds.Criterion{}).Where("status = 'active'").Order("id ASC")
	if q != "" {
		db = db.Where("name ILIKE ?", "%"+q+"%")
	}
	return list, db.Find(&list).Error
}

func (r *Repository) GetCriterionByID(id uint) (*ds.Criterion, error) {
	var c ds.Criterion
	err := r.db.First(&c, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	return &c, err
}

func (r *Repository) CreateCriterion(c *ds.Criterion) error {
	return r.db.Create(c).Error
}

func (r *Repository) UpdateCriterion(id uint, c *request.UpdateCriterion) error {
	tx := r.db.Model(&ds.Criterion{}).Where("id = ?", id).Updates(c)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) DeleteCriterion(id uint) error {
	tx := r.db.Delete(&ds.Criterion{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// TODO нормальную ошибку кидать, а не пойми что
func (r *Repository) AddItem(orderID, criterionID uint) error {
	var lastOI ds.PankreatitOrderItem
	r.db.Last(&lastOI, "med_order_id = ?", orderID)
	item := ds.PankreatitOrderItem{PankreatitOrderID: orderID, CriterionID: criterionID, Position: lastOI.Position + 1}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&item).Error
}

func (r *Repository) GetSeq() (uint, error) {
	var nextID int64
	if err := r.db.Raw(`SELECT nextval('criteria_id_seq')`).Scan(&nextID).Error; err != nil {
		return 999999, err
	}
	return uint(nextID), nil
}

func (r *Repository) ResetCriterionSequence() error {
	sql := `
        SELECT setval(
            pg_get_serial_sequence('criteria', 'id'),
            COALESCE((SELECT MAX(id) FROM criteria), 0)
        )
    `
	return r.db.Exec(sql).Error
}

func (r *Repository) GetImageName(critid uint) (string, error) {
	var item ds.Criterion
	if err := r.db.Where("id = ?", critid).First(&item).Error; err != nil {
		return "", err
	}
	if item.ImageURL == nil {
		return "", nil
	} else {
		return *item.ImageURL, nil
	}

}
