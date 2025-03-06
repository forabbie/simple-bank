package token

import (
	"testing"
	"time"

	"github.com/forabbie/vank-app/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

// func TestInvalidToken(t *testing.T) {
// 	symmetricKey := "12345678901234567890123456789012" // 32-byte key
// 	maker, err := NewPasetoMaker(symmetricKey)
// 	require.NoError(t, err)

// 	// Create a valid token
// 	token, err := maker.CreateToken("testuser", time.Minute)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, token)

// 	// Tamper with the token to make it invalid
// 	invalidToken := token[:len(token)-1] + "x"

// 	// Verify the tampered token
// 	payload, err := maker.VerifyToken(invalidToken)
// 	require.Error(t, err)
// 	require.Nil(t, payload)
// }
