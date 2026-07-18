package mapper

import (
	"pankreatitmed/internal/app/ds"
	"pankreatitmed/internal/app/dto/request"
	"pankreatitmed/internal/app/dto/response"
)

func MedUserRegistrationToMedUser(usr request.MedUserRegistration) ds.MedUser {
	return ds.MedUser{
		Login:    usr.Login,
		Password: usr.Password,
	}
}

func AuthenticateMedUserToMedUser(usr request.AuthenticateMedUser) ds.MedUser {
	return ds.MedUser{
		Login:    usr.Login,
		Password: usr.Password,
	}
}

func MedUserToSendMedUserFields(user *ds.MedUser) response.SendMedUserField {
	return response.SendMedUserField{
		ID:          user.ID,
		Login:       user.Login,
		IsModerator: user.IsModerator,
	}
}
