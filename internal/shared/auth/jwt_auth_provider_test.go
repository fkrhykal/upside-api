package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/fkrhykal/upside-api/internal/shared/auth"
	"github.com/fkrhykal/upside-api/internal/shared/helpers"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type JwtAuthProviderTestSuite struct {
	suite.Suite
	jwtAuthProvider auth.AuthProvider
}

func (cs *JwtAuthProviderTestSuite) SetupSuite() {
	logger := log.NewTestLogger(cs.T())
	cs.jwtAuthProvider = auth.NewJwtAuthProvider(logger, &auth.JwtAuthConfig{
		SignedKey: []byte("secret"),
	})
}

func (cs *JwtAuthProviderTestSuite) TestGenerateToken() {
	ctx := context.Background()
	_, err := cs.jwtAuthProvider.GenerateToken(ctx, &auth.UserCredential{ID: uuid.New()}, time.Now())
	cs.Nil(err)
}

func (cs *JwtAuthProviderTestSuite) TestRetrieveCredential() {
	ctx := context.Background()
	id := uuid.New()
	token, err := cs.jwtAuthProvider.GenerateToken(ctx, &auth.UserCredential{ID: id}, time.Now().Add(helpers.WEEK))
	if err != nil {
		cs.FailNow(err.Error())
	}
	userCredential, err := cs.jwtAuthProvider.RetrieveCredential(ctx, token)
	cs.Nil(err)
	cs.NotNil(userCredential)
	cs.EqualValues(id, userCredential.ID)
}

func (cs *JwtAuthProviderTestSuite) TestRetrieveInvalidCredential() {
	ctx := context.Background()
	id := uuid.New()
	_, err := cs.jwtAuthProvider.GenerateToken(ctx, &auth.UserCredential{ID: id}, time.Now().Add(helpers.WEEK))
	if err != nil {
		cs.FailNow(err.Error())
	}
	wrongToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyQ3JlZGVudGlhbCI6eyJpZCI6ImI0YzkzODNkLTM5NjAtNDMxNi1hYTAyLWE4YWEwOTEzNzNjZSJ9fQ.wsVP0lKUlYkfOXOVwzrjoxY3Yx7Uagibm1p3d2Gc3T4"
	userCredential, err := cs.jwtAuthProvider.RetrieveCredential(ctx, auth.Token(wrongToken))
	cs.NotNil(err)
	cs.Nil(userCredential)
}

func (cs *JwtAuthProviderTestSuite) TestRetrieveExpiredCredential() {
	ctx := context.Background()
	id := uuid.New()
	token, err := cs.jwtAuthProvider.GenerateToken(ctx, &auth.UserCredential{ID: id}, time.Now())
	cs.Nil(err)

	userCredential, err := cs.jwtAuthProvider.RetrieveCredential(ctx, token)
	cs.NotNil(err)
	cs.Nil(userCredential)
}

func TestCredentialService(t *testing.T) {
	suite.Run(t, new(JwtAuthProviderTestSuite))
}
