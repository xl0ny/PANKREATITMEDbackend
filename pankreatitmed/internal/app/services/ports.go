package services

import (
	"pankreatitmed/internal/app/ds"
	"pankreatitmed/internal/app/dto/request"
	"time"
)

type CriteriaRepoPort interface {
	GetCriteria(q string) ([]ds.Criterion, error)
	GetCriterionByID(id uint) (*ds.Criterion, error)
	CreateCriterion(c *ds.Criterion) error
	UpdateCriterion(id uint, in *request.UpdateCriterion) error
	DeleteCriterion(id uint) error
	AddItem(orderID, criterionID uint) error
	GetSeq() (uint, error)
	ResetCriterionSequence() error
	GetOrCreateDraftPankreatitOrder(creatorID uint) (*ds.PankreatitOrder, error)
	GetImageName(critid uint) (string, error)
}

type PankreatitOrdersRepoPort interface {
	CountItems(orderID uint) (int64, error)
	IsPankreatitOrderDeleted(orderID uint) (bool, error)
	IsPankreatitOrderDraft(orderID uint) (bool, error)
	IsPankreatitOrderFormed(orderID uint) (bool, error)

	GetOrCreateDraftPankreatitOrder(creatorID uint) (*ds.PankreatitOrder, error)
	GetPankreatitOrders(userID uint, status *string, start, end *time.Time) ([]ds.PankreatitOrder, error)
	GetPankreatitOrderWithItems(orderID uint) (ds.PankreatitOrder, []ds.PankreatitOrderItem, error)

	UpdatePankreatitOrder(id uint, order *request.UpdatePankreatitOrder) error
	FormPankreatitOrder(id uint) error
	EndOrCancelPankreatitOrder(id, moderator uint, status string) error
	SoftDeleteOrderSQL(orderID uint) error
	SetRansonAndRisk(orderID uint, score int, risk string) error
	GetCriterionByID(id uint) (*ds.Criterion, error)
}

type PankreatitOrderItemsRepoPort interface {
	DeleteFromPankreatitOrder(medorder, criterion uint) error
	UpdatePankreatitOrderItem(medorder, criterion uint, position *uint, val *float64) error
}

type MedUsersRepoPort interface {
	CreateMedUser(user *ds.MedUser) error
	GetMedUserByLogin(login string) (*ds.MedUser, error)
	ChangeMedUser(id uint, user *request.UpdateMedUser) error
	GetMedUserByID(id uint) (*ds.MedUser, error)
}
