package request

type MedUserRegistration struct {
	Login    string `form:"login" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type GetMedUser struct {
	Login string `json:"login" binding:"required"`
}

type UpdateMedUser struct {
	Login    *string `json:"login"`
	Password *string `json:"password"`
}

type AuthenticateMedUser struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
