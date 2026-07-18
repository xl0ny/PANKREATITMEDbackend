package services

import (
	"errors"
	"fmt"
	"pankreatitmed/internal/app/ds"
	"pankreatitmed/internal/app/dto/request"
	"pankreatitmed/internal/app/dto/response"
	"pankreatitmed/internal/app/mapper"
	"pankreatitmed/internal/app/middleware"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrLoginTaken         = errors.New("login already taken")
	ErrWeakPassword       = errors.New("password must be at least 6 chars")
	ErrLoginIsRequired    = errors.New("login is required")
	ErrPasswordIsRequired = errors.New("password is required")
)

type MedUsersService interface {
	Login(au request.AuthenticateMedUser) (string, error)
	Register(usr request.MedUserRegistration) (ds.MedUser, string, error)
	Logout(token string) error
	GetMyField(id uint) (*response.SendMedUserField, error)
	UpdateField(id uint, user *request.UpdateMedUser) error
	GetConfig() middleware.JWTConfig
}

type medUsersService struct {
	repo         MedUsersRepoPort
	jwtConfig    middleware.JWTConfig
	jwtBlackList *middleware.RedisBlacklist
}

func NewMedUsersService(repo MedUsersRepoPort, jwtconfig middleware.JWTConfig, jwtblacklist *middleware.RedisBlacklist) MedUsersService {
	return &medUsersService{repo: repo, jwtConfig: jwtconfig, jwtBlackList: jwtblacklist}
}

func (s *medUsersService) Login(au request.AuthenticateMedUser) (string, error) {
	u, err := s.repo.GetMedUserByLogin(au.Login)
	if err != nil {
		return "", errors.New("MedUser not found")
	}
	println()
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(au.Password))
	if err != nil {
		return "", errors.New("Invalid credentials")
	}
	return s.IssueJWT(u, s.jwtConfig.TTL)
}

func (s *medUsersService) Register(usr request.MedUserRegistration) (ds.MedUser, string, error) {
	if len(usr.Password) < 6 {
		return ds.MedUser{}, "", ErrWeakPassword
	}
	if usr.Login == "" {
		return ds.MedUser{}, "", ErrLoginIsRequired
	}

	if _, err := s.repo.GetMedUserByLogin(usr.Login); err == nil {
		return ds.MedUser{}, "", ErrLoginTaken
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return ds.MedUser{}, "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(usr.Password), bcrypt.DefaultCost)
	if err != nil {
		return ds.MedUser{}, "", err
	}
	user := mapper.MedUserRegistrationToMedUser(usr)
	user.IsModerator = false
	user.Password = string(hash)

	if err := s.repo.CreateMedUser(&user); err != nil {
		if isUniqueViolation(err, "users_login_key") {
			return ds.MedUser{}, "", ErrLoginTaken
		}
		return ds.MedUser{}, "", err
	}
	token, err := s.Login(request.AuthenticateMedUser{Login: user.Login, Password: usr.Password})
	if err != nil {
		return ds.MedUser{}, "", err
	}
	return user, token, nil
}

func (s *medUsersService) Logout(token string) error {
	return s.jwtBlackList.Add(token, time.Now().Add(30*24*time.Hour)) //баним на месяц
}

func (s *medUsersService) GetMyField(id uint) (*response.SendMedUserField, error) {
	user, err := s.repo.GetMedUserByID(id)
	if err != nil {
		return nil, err
	}
	res := mapper.MedUserToSendMedUserFields(user)
	return &res, nil
}

func (s *medUsersService) UpdateField(id uint, user *request.UpdateMedUser) error {
	if user.Password != nil {
		if len(*user.Password) < 6 {
			return ErrWeakPassword
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		pwd := string(hash)
		if err != nil {
			return err
		}
		user.Password = &pwd
	}

	return s.repo.ChangeMedUser(id, user)
}

func (s *medUsersService) IssueJWT(u *ds.MedUser, ttl time.Duration) (string, error) {
	now := time.Now()

	claims := middleware.Claims{
		Sub:         u.ID,
		Login:       u.Login,
		IsModerator: u.IsModerator,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "pankreatitmed",
			Subject:   fmt.Sprint(u.ID),
			// Audience:  []string{"spa"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token == nil {
		return "", errors.New("token is nil")
	}
	return token.SignedString([]byte(s.jwtConfig.Secret))
}

func (s *medUsersService) GetConfig() middleware.JWTConfig {
	return s.jwtConfig
}

// TODO добавить правило из бд
func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return constraint == "" || pgErr.ConstraintName == constraint
	}
	return false
}
