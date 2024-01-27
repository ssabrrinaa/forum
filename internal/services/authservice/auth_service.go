package authservice

import (
	"errors"
	"forum/internal/models"
	"forum/internal/repositories/authrepo"
	"forum/internal/schemas"
	"forum/pkg/hashbcrypt"
	"forum/pkg/validator"
	"log"
	"time"

	"github.com/gofrs/uuid"
)

type AuthService struct {
	AuthRepo authrepo.AuthRepoI
}

func NewAuthService(authRepo authrepo.AuthRepoI) *AuthService {
	return &AuthService{
		AuthRepo: authRepo,
	}
}

type AuthServiceI interface {
	CreateUser(user schemas.CreateUser) error
	CreateSession(user schemas.AuthUser) (models.Session, error)

	CheckUserPassword(email, password string) error
	GetSession(token string) (models.Session, error)

	DeleteSession(token string) error
}

func (as *AuthService) CreateUser(user schemas.CreateUser) error {
	if err := validator.ValidateRegisterInput(user); err != nil {
		return err
	}

	hashedPassword, err := hashbcrypt.GenerateHashedPassword(user.Password)
	if err != nil {
		return err
	}
	// как чекать если email уже сущетсвует в БД

	userFromDb, _ := as.AuthRepo.GetUserByEmail(user.Email)
	if userFromDb.Email == user.Email {
		return errors.New("User already exists") // handle the error properly
	}
	user_id, err := uuid.NewV4()
	if err != nil {
		log.Fatal("error while generating uuid")
	}
	user_model := models.User{
		ID:             user_id,
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: hashedPassword,
	}
	if err := as.AuthRepo.CreateUser(user_model); err != nil {
		return err // handle the error properly
	}

	return nil
}

func (as *AuthService) CreateSession(user schemas.AuthUser) (models.Session, error) {
	session := models.Session{}
	if err := validator.ValidateSignInInput(user); err != nil {
		return session, err
	}

	userDB, err := as.AuthRepo.GetUserByEmail(user.Email)
	if err != nil {
		return session, err
	}

	err = as.CheckUserPassword(userDB.HashedPassword, user.Password)
	if err != nil {
		return session, err
	}

	if err := as.AuthRepo.DeleteSession(userDB.ID); err != nil {
		return session, err
	}

	token, err := uuid.NewV4()
	session_id, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	session = models.Session{
		ID:         session_id,
		UserID:     userDB.ID,
		Token:      token.String(),
		ExpireTime: time.Now().Add(time.Hour * 2),
	}

	err = as.AuthRepo.CreateSession(session)
	if err != nil {
		return models.Session{}, err
	}

	return session, nil
}

func (as *AuthService) CheckUserPassword(passwordDB, password string) error {
	// user, err := as.AuthRepo.GetUserByEmail(email)
	// if err != nil {
	// 	return err
	// }

	if !hashbcrypt.CheckHashedPassword(password, passwordDB) {
		return errors.New("Wrong password") // где хранить все ошибки?
	}

	return nil
}

func (as *AuthService) DeleteSession(token string) error {
	/*
		get session
		exists -> check the time
		delete
	*/
	session, err := as.AuthRepo.GetSession(token)
	if err != nil {
		log.Fatal("No session found")
	}

	err = as.AuthRepo.DeleteSession(session.ID)
	if err != nil {
		log.Fatal("session delete issue")
	}
	return nil
}

func (as *AuthService) GetSession(token string) (models.Session, error) {
	session, err := as.AuthRepo.GetSession(token)
	if err != nil {
		return models.Session{}, err
	}
	return session, nil
}

// идея продлить сессию когда остается 5 минут или типо того
// func (as *AuthService) ExtendSession(token string) (error) {
// 	session, _ := as.GetSession(token)

// }
