package token

import (
	"github.com/Bakhram74/small_bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)
	duration := time.Minute
	username := util.RandomOwner()
	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)
	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.UserName)

	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, _ := NewPasetoMaker(util.RandomString(32))

	username := util.RandomOwner()

	token, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	payload, err := maker.VerifyToken(token)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

//func TestInvalidPasetoToken(t *testing.T) {
//	payload, err := NewPayload("Alex", time.Minute)
//	require.NoError(t, err)
//	require.NotEmpty(t, payload)
//	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
//	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
//	require.NoError(t, err)
//
//	maker, _ := NewJWTMaker(util.RandomString(32))
//	payload, err = maker.VerifyToken(token)
//
//	require.EqualError(t, err, ErrInvalidToken.Error())
//	require.Nil(t, payload)
//}
