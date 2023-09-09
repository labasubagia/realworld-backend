package token

import (
	"time"

	"github.com/google/uuid"
	"github.com/labasubagia/realworld-backend/internal/core/util/exception"
)

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, exception.New(exception.TypeInternal, "failed generate token id", err)
	}
	payload := &Payload{
		ID:        tokenID,
		Email:     username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return exception.New(exception.TypeTokenExpired, "token expired", nil)
	}
	return nil
}
