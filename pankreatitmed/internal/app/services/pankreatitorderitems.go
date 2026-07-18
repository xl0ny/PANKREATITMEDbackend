package services

type PankreatitOrderItemsService interface {
	Delete(pankreatitorder, criterion uint) error
	Update(pankreatitorder, criterion uint, position *uint, val *float64) error
}

type pankreatitOrderItemsService struct {
	repo PankreatitOrderItemsRepoPort
}

func NewPankreatitOrderItemsService(repo PankreatitOrderItemsRepoPort) PankreatitOrderItemsService {
	return &pankreatitOrderItemsService{repo: repo}
}
func (s *pankreatitOrderItemsService) Delete(pankreatitorder, criterion uint) error {
	return s.repo.DeleteFromPankreatitOrder(pankreatitorder, criterion)
}

func (s *pankreatitOrderItemsService) Update(pankreatitorder, criterion uint, position *uint, val *float64) error {
	return s.repo.UpdatePankreatitOrderItem(pankreatitorder, criterion, position, val)
}
