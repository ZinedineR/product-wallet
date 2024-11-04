package service_test

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"product-wallet/internal/entity"
	"product-wallet/internal/mocks"
	"product-wallet/internal/model"
	service "product-wallet/internal/services"
	mocksSignature "product-wallet/pkg/mocks"
	"product-wallet/pkg/xvalidator"
	"testing"
)

func setupSQLMock(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	// Setup SQL mock
	db, mockSql, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	// Setup GORM with the mock DB
	gormDB, gormDBErr := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if gormDBErr != nil {
		t.Fatalf("failed to open GORM connection: %v", gormDBErr)
	}
	return mockSql, gormDB
}
func TestRegisterUser(t *testing.T) {
	mockAppCtx := context.Background()

	t.Run("RegisterUser Success", func(t *testing.T) {
		// Set up input
		request := &model.CreateUserReq{
			BaseUserReq: model.BaseUserReq{
				Username: "john_doe",
				Password: "SecurePass123!",
			},
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("FindByFilter", mockAppCtx, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		mockRepository.On("CreateTx", mockAppCtx, mock.Anything, mock.Anything).Return(nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("HashBscryptPassword", request.Password).Return("$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu", nil)

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectCommit()
		_, errService := mockService.Register(mockAppCtx, request)

		// Assert the result
		assert.Nil(t, errService)
	})

	t.Run("RegisterUser Username Exists", func(t *testing.T) {
		// Set up input
		request := &model.CreateUserReq{
			BaseUserReq: model.BaseUserReq{
				Username: "john_doe",
				Password: "SecurePass123!",
			},
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		existingUser := &entity.User{
			Id:       "123e4567-e89b-12d3-a456-426614174000",
			Username: "john_doe",
		}
		mockRepository.On("FindByFilter", mockAppCtx, mock.Anything, mock.Anything, mock.Anything).Return(existingUser, nil)

		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		_, errService := mockService.Register(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
	})

	t.Run("RegisterUser Validation Failed", func(t *testing.T) {
		// Set up input (invalid data - missing password)
		request := &model.CreateUserReq{
			BaseUserReq: model.BaseUserReq{
				Username: "john_doe",
			},
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		_, errService := mockService.Register(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
	})

	t.Run("RegisterUser HashPassword Failed", func(t *testing.T) {
		// Set up input
		request := &model.CreateUserReq{
			BaseUserReq: model.BaseUserReq{
				Username: "john_doe",
				Password: "SecurePass123!",
			},
		}

		// Mocks
		mockSql, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("FindByFilter", mockAppCtx, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("HashBscryptPassword", request.Password).Return("", errors.New("hash error"))

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		mockSql.ExpectBegin()
		mockSql.ExpectRollback()
		_, errService := mockService.Register(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
	})
}

func TestLoginUser(t *testing.T) {
	mockAppCtx := context.Background()

	t.Run("LoginUser Success", func(t *testing.T) {
		// Set up input
		request := &model.CreateUserReq{
			BaseUserReq: model.BaseUserReq{
				Username: "john_doe",
				Password: "SecurePass123!",
			},
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		existingUser := &entity.User{
			Id:       "123e4567-e89b-12d3-a456-426614174000",
			Username: "john_doe",
			Password: "$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu", // Hashed password
		}
		mockRepository.On("FindByFilter", mockAppCtx, mock.Anything, mock.Anything, mock.Anything).Return(existingUser, nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("CheckBscryptPasswordHash", request.Password, existingUser.Password).Return(true)
		mockSignaturer.On("GenerateJWT", existingUser.Username).Return("jwt_token", nil)

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.Login(mockAppCtx, request)

		// Assert the result
		assert.Nil(t, errService)
		assert.NotNil(t, result)
	})

	t.Run("LoginUser Username Not Found", func(t *testing.T) {
		// Set up input
		request := &model.CreateUserReq{
			BaseUserReq: model.BaseUserReq{
				Username: "john_doe",
				Password: "SecurePass123!",
			},
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		mockRepository.On("FindByFilter", mockAppCtx, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		validate, _ := xvalidator.NewValidator()
		mockSignaturer := new(mocksSignature.Signaturer)
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.Login(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})

	t.Run("LoginUser Invalid Password", func(t *testing.T) {
		// Set up input
		request := &model.CreateUserReq{
			BaseUserReq: model.BaseUserReq{
				Username: "john_doe",
				Password: "WrongPass123!",
			},
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		existingUser := &entity.User{
			Id:       "123e4567-e89b-12d3-a456-426614174000",
			Username: "john_doe",
			Password: "$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu", // Hashed password
		}
		mockRepository.On("FindByFilter", mockAppCtx, mock.Anything, mock.Anything, mock.Anything).Return(existingUser, nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("CheckBscryptPasswordHash", request.Password, existingUser.Password).Return(false)

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.Login(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})

	t.Run("LoginUser JWT Generation Failed", func(t *testing.T) {
		// Set up input
		request := &model.CreateUserReq{
			BaseUserReq: model.BaseUserReq{
				Username: "john_doe",
				Password: "SecurePass123!",
			},
		}

		// Mocks
		_, gormDB := setupSQLMock(t)
		mockRepository := new(mocks.UserRepository)
		existingUser := &entity.User{
			Id:       "123e4567-e89b-12d3-a456-426614174000",
			Username: "john_doe",
			Password: "$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu", // Hashed password
		}
		mockRepository.On("FindByFilter", mockAppCtx, mock.Anything, mock.Anything, mock.Anything).Return(existingUser, nil)
		mockSignaturer := new(mocksSignature.Signaturer)
		mockSignaturer.On("CheckBscryptPasswordHash", request.Password, existingUser.Password).Return(true)
		mockSignaturer.On("GenerateJWT", existingUser.Username).Return("", errors.New("jwt error"))

		validate, _ := xvalidator.NewValidator()
		mockService := service.NewUserService(gormDB, mockRepository, mockSignaturer, validate)

		// Call the function under test
		result, errService := mockService.Login(mockAppCtx, request)

		// Assert the result
		assert.NotNil(t, errService)
		assert.Nil(t, result)
	})
}
