package services

import (
	"pankreatitmed/internal/app/ds"
	"pankreatitmed/internal/app/dto/request"
	"pankreatitmed/internal/app/dto/response"
	"pankreatitmed/internal/app/mapper"
	"strconv"
	"time"
	//"pankreatitmed/internal/app/dto/request"
	//"pankreatitmed/internal/app/singleton"
	//"time"
	"errors"
)

type PankreatitOrdersService interface {
	GetDraft(creatorID uint) (*response.SendCartPankreatitOrder, error)
	List(userID uint, status *string, start, end *time.Time) ([]response.SendPankreatitOrders, error)
	Get(ID uint) (response.SendPankreatitOrder, error)
	Update(ID uint, in *request.UpdatePankreatitOrder) error
	Form(ID uint) error
	CancelOrEnd(ID, moderator uint, status string) error
	Delete(ID uint) error
}

type pankreatitOrdersService struct {
	repo PankreatitOrdersRepoPort
}

func NewPankreatitOrdersService(repo PankreatitOrdersRepoPort) PankreatitOrdersService {
	return &pankreatitOrdersService{repo: repo}
}

// TODO перенести сюда singleton из хэндлера
func (s *pankreatitOrdersService) GetDraft(creatorID uint) (*response.SendCartPankreatitOrder, error) {
	o, err := s.repo.GetOrCreateDraftPankreatitOrder(creatorID)
	if err != nil {
		return nil, err
	}
	amnt, err := s.repo.CountItems(o.ID)
	if err != nil {
		return nil, err
	}
	res := mapper.PankreatitOrderToSendPankreatitOrder(o, uint(amnt))
	return &res, err
}

func (s *pankreatitOrdersService) List(userID uint, status *string, start, end *time.Time) ([]response.SendPankreatitOrders, error) {
	morders, err := s.repo.GetPankreatitOrders(userID, status, start, end)
	if err != nil {
		return nil, err
	}
	res := mapper.PankreatitOrdersToSendPankreatitOrders(morders)
	return res, nil
}

func (s *pankreatitOrdersService) Get(ID uint) (response.SendPankreatitOrder, error) {
	o, items, err := s.repo.GetPankreatitOrderWithItems(ID)
	res := mapper.PankreatitOrderToSendPankreatitOrderWithItems(o, items)
	return res, err
}

func (s *pankreatitOrdersService) Update(ID uint, in *request.UpdatePankreatitOrder) error {
	return s.repo.UpdatePankreatitOrder(ID, in)
}

func (s *pankreatitOrdersService) Form(ID uint) error {
	check, err := s.repo.IsPankreatitOrderDraft(ID)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("MedOrderIsNotDraft")
	}
	return s.repo.FormPankreatitOrder(ID)
}

func (s *pankreatitOrdersService) CancelOrEnd(ID, moderator uint, status string) error {
	check, err := s.repo.IsPankreatitOrderFormed(ID)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("MedOrderIsNotFormed")
	}
	_, criteria, err := s.repo.GetPankreatitOrderWithItems(ID)
	if err != nil {
		return err
	}
	if status == "completed" && CheckReadyToEnd(criteria) {
		rans, rsk, err := s.computeRanson(criteria)
		if err != nil {
			return err
		}
		if err := s.repo.SetRansonAndRisk(ID, rans, rsk); err != nil {
			return err
		}
		if err := s.repo.EndOrCancelPankreatitOrder(ID, moderator, status); err != nil {
			return err
		}
	} else if status == "rejected" {
		if err := s.repo.EndOrCancelPankreatitOrder(ID, moderator, status); err != nil {
			return err
		}
	} else {
		return errors.New("Not all value fields are complete")
	}
	return nil
}

func CheckReadyToEnd(items []ds.PankreatitOrderItem) bool {
	for _, item := range items {
		if item.ValueNum == nil {
			return false
		}
	}
	return true
}

func (s *pankreatitOrdersService) computeRanson(items []ds.PankreatitOrderItem) (int, string, error) {
	var score int
	for _, item := range items {
		crit, err := s.repo.GetCriterionByID(item.CriterionID)
		if err != nil {
			return 0, "", err
		}

		if crit.RefHigh != nil {
			if *item.ValueNum > *crit.RefHigh {
				score++
			}
		} else if crit.RefLow != nil {
			if *item.ValueNum < *crit.RefLow {
				score++
			}
		} else {
			return 0, "", errors.New("One of Ref is null")
		}
	}
	println(score)
	return score, strconv.Itoa(score*100/11) + "%", nil
}

func (s *pankreatitOrdersService) Delete(ID uint) error {
	return s.repo.SoftDeleteOrderSQL(ID)
}
