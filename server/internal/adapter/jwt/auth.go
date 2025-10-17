package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hr-platform-mosprom/internal/core/application/port"
)

type authClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type jwtService struct {
	clock     port.Clock
	tokenTTL  time.Duration
	secretKey []byte
}

func NewJWTService(clock port.Clock, tokenTTL time.Duration, secretKey string) *jwtService {
	return &jwtService{
		clock:     clock,
		tokenTTL:  tokenTTL,
		secretKey: []byte(secretKey),
	}
}

func (s *jwtService) Generate(payload port.TokenPayload) (string, error) {
	claims := authClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   payload.Sub.String(),
			ExpiresAt: jwt.NewNumericDate(s.clock.Now().Add(s.tokenTTL)),
		},
		Role: string(payload.Role),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secretKey)
}

func (s *jwtService) Validate(tokenStr string) (port.TokenPayload, error) {
	var claims authClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(_ *jwt.Token) (any, error) {
		return s.secretKey, nil
	})
	if err != nil || !token.Valid {
		return port.TokenPayload{}, errors.New("invalid or expired auth token")
	}

	uid, err := uuid.Parse(claims.Subject)
	if err != nil {
		return port.TokenPayload{}, err
	}

	return port.TokenPayload{
		Sub:  uid,
		Role: port.Role(claims.Role),
	}, nil
}
