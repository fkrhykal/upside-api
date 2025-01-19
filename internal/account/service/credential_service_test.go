package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/fkrhykal/upside-api/internal/account/dto"
	"github.com/fkrhykal/upside-api/internal/account/service"
	"github.com/fkrhykal/upside-api/internal/shared/helpers"
	"github.com/fkrhykal/upside-api/internal/shared/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type CredentialServiceTestSuite struct {
	suite.Suite
	credentialService service.CredentialService
}

func (cs *CredentialServiceTestSuite) SetupSuite() {
	logger := log.NewTestLogger(cs.T())
	cs.credentialService = service.NewJwtCredentialService(logger, &service.JwtCredentialConfig{
		SignedKey: []byte("secret"),
	})
}

func (cs *CredentialServiceTestSuite) TestGenerateToken() {
	ctx := context.Background()
	_, err := cs.credentialService.GenerateToken(ctx, &dto.UserCredential{ID: uuid.New()}, time.Now())
	cs.Nil(err)
}

func (cs *CredentialServiceTestSuite) TestRetrieveCredential() {
	ctx := context.Background()
	id := uuid.New()
	token, err := cs.credentialService.GenerateToken(ctx, &dto.UserCredential{ID: id}, time.Now().Add(helpers.WEEK))
	if err != nil {
		cs.FailNow(err.Error())
	}
	userCredential, err := cs.credentialService.RetrieveUserCredential(ctx, token)
	cs.Nil(err)
	cs.NotNil(userCredential)
	cs.EqualValues(id, userCredential.ID)
}

func (cs *CredentialServiceTestSuite) TestRetrieveInvalidCredential() {
	ctx := context.Background()
	id := uuid.New()
	_, err := cs.credentialService.GenerateToken(ctx, &dto.UserCredential{ID: id}, time.Now().Add(helpers.WEEK))
	if err != nil {
		cs.FailNow(err.Error())
	}
	wrongToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyQ3JlZGVudGlhbCI6eyJpZCI6ImI0YzkzODNkLTM5NjAtNDMxNi1hYTAyLWE4YWEwOTEzNzNjZSJ9fQ.wsVP0lKUlYkfOXOVwzrjoxY3Yx7Uagibm1p3d2Gc3T4"
	userCredential, err := cs.credentialService.RetrieveUserCredential(ctx, dto.CredentialToken(wrongToken))
	cs.NotNil(err)
	cs.Nil(userCredential)
}

func (cs *CredentialServiceTestSuite) TestRetrieveExpiredCredential() {
	ctx := context.Background()
	id := uuid.New()
	token, err := cs.credentialService.GenerateToken(ctx, &dto.UserCredential{ID: id}, time.Now())
	cs.Nil(err)

	userCredential, err := cs.credentialService.RetrieveUserCredential(ctx, token)
	cs.NotNil(err)
	cs.Nil(userCredential)
}

func TestCredentialService(t *testing.T) {
	suite.Run(t, new(CredentialServiceTestSuite))
}
