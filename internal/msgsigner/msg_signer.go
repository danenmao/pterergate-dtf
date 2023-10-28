package msgsigner

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CommonClaims struct {
	Message string `json:"message"`
	jwt.RegisteredClaims
}

type MsgSigner struct {
	keeper     *KeyKeeper
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewMsgSigner() *MsgSigner {
	signer := &MsgSigner{}
	signer.Init()
	return signer
}

func (s *MsgSigner) Init() {
	s.keeper = &KeyKeeper{}
	err := s.keeper.Load()
	if err != nil {
		return
	}

	s.privateKey, err = s.keeper.GetPrivateKey()
	if err != nil {
		return
	}

	s.publicKey, err = s.keeper.GetPublicKey()
	if err != nil {
		return
	}
}

func (s *MsgSigner) Sign(
	issuer string, subject string, audience []string, msg string, expire time.Duration,
) (string, error) {
	now := time.Now()
	claims := CommonClaims{
		Message: msg,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expire)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    issuer,
			ID:        uuid.NewString(),
			Audience:  audience,
			Subject:   subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

func (s *MsgSigner) Verify(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CommonClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.publicKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid claims")
	}

	claims, ok := token.Claims.(*CommonClaims)
	if !ok {
		return "", errors.New("invalid claim type")
	}

	return claims.Message, nil
}
