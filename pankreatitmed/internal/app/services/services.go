package services

import "pankreatitmed/internal/app/middleware"

type Services struct {
	Criteria             CriteriaService
	PankreatitOrders     PankreatitOrdersService
	PankreatitOrderItems PankreatitOrderItemsService
	MedUsers             MedUsersService
}

type Reps struct {
	CriteriaRepo             CriteriaRepoPort
	PankreatitOrdersRepo     PankreatitOrdersRepoPort
	PankreatitOrderItemsRepo PankreatitOrderItemsRepoPort
	MedUsersRepo             MedUsersRepoPort
}

type Configs struct {
	JWTConfig    middleware.JWTConfig
	JWTBlackList *middleware.RedisBlacklist
}

func NewServices(d Reps, c Configs) *Services {
	return &Services{
		Criteria:             NewCriteriaService(d.CriteriaRepo),
		PankreatitOrders:     NewPankreatitOrdersService(d.PankreatitOrdersRepo),
		PankreatitOrderItems: NewPankreatitOrderItemsService(d.PankreatitOrderItemsRepo),
		MedUsers:             NewMedUsersService(d.MedUsersRepo, c.JWTConfig, c.JWTBlackList),
	}
}
