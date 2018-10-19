package auth

import (
	"errors"
	"github.com/openaccounting/oa-server/core/model/db"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
)

var Instance Interface

type AuthService struct {
	db     db.Datastore
	bcrypt util.Bcrypt
}

type Interface interface {
	Authenticate(string, string) (*types.User, error)
	AuthenticateUser(email string, password string) (*types.User, error)
	AuthenticateSession(string) (*types.User, error)
	AuthenticateApiKey(string) (*types.User, error)
}

func NewAuthService(db db.Datastore, bcrypt util.Bcrypt) *AuthService {
	authService := &AuthService{db: db, bcrypt: bcrypt}
	Instance = authService
	return authService
}

func (auth *AuthService) Authenticate(emailOrKey string, password string) (*types.User, error) {
	// authenticate via session, apikey or user
	user, err := auth.AuthenticateSession(emailOrKey)

	if err == nil {
		return user, nil
	}

	user, err = auth.AuthenticateApiKey(emailOrKey)

	if err == nil {
		return user, nil
	}

	user, err = auth.AuthenticateUser(emailOrKey, password)

	if err == nil {
		return user, nil
	}

	return nil, errors.New("Unauthorized")
}

func (auth *AuthService) AuthenticateUser(email string, password string) (*types.User, error) {
	u, err := auth.db.GetVerifiedUserByEmail(email)

	if err != nil {
		return nil, errors.New("Invalid email or password")
	}

	err = auth.bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	if err != nil {
		return nil, errors.New("Invalid email or password")
	}

	return u, nil
}

func (auth *AuthService) AuthenticateSession(id string) (*types.User, error) {
	u, err := auth.db.GetUserByActiveSession(id)

	if err != nil {
		return nil, errors.New("Invalid session")
	}

	auth.db.UpdateSessionActivity(id)

	return u, nil
}

func (auth *AuthService) AuthenticateApiKey(id string) (*types.User, error) {
	u, err := auth.db.GetUserByApiKey(id)

	if err != nil {
		return nil, errors.New("Access denied")
	}

	auth.db.UpdateApiKeyActivity(id)

	return u, nil
}
