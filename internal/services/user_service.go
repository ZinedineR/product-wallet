package service

import (
	"context"
	"product-wallet/internal/model"
	"product-wallet/internal/repository"
	"product-wallet/pkg/signature"

	//"product-wallet/pkg/exception"
	"gorm.io/gorm"
	"product-wallet/pkg/exception"
	"product-wallet/pkg/xvalidator"
)

type UserServiceImpl struct {
	db         *gorm.DB
	userRepo   repository.UserRepository
	signaturer signature.Signaturer
	validate   *xvalidator.Validator
}

func NewUserService(
	db *gorm.DB, repo repository.UserRepository,
	signaturer signature.Signaturer,
	validate *xvalidator.Validator,
) UserService {
	return &UserServiceImpl{
		db:         db,
		userRepo:   repo,
		signaturer: signaturer,
		validate:   validate,
	}
}

func (s *UserServiceImpl) Register(
	ctx context.Context, req *model.CreateUserReq,
) (*model.CreateUserRes, *exception.Exception) {
	tx := s.db.Begin()
	defer tx.Rollback()
	if errs := s.validate.Struct(req); errs != nil {
		return nil, exception.InvalidArgument(errs)
	}
	duplicateCheck, err := s.userRepo.FindByFilter(ctx, s.db, model.FilterParams{
		{
			Field:    "username",
			Value:    req.Username,
			Operator: "=",
		},
	}, model.OrderParam{
		Order:   "desc",
		OrderBy: "username",
	})
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if duplicateCheck != nil {
		return nil, exception.PermissionDenied("username already exists")
	}
	password, err := s.signaturer.HashBscryptPassword(req.Password)
	if err != nil {
		return nil, exception.Internal("can't create password", err)
	}
	body := req.ToEntity(password)
	if err := s.userRepo.CreateTx(ctx, tx, body); err != nil {
		return nil, exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return &model.CreateUserRes{
		User: *body,
	}, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, req *model.CreateUserReq) (
	*model.LoginUserRes, *exception.Exception,
) {
	if errs := s.validate.Struct(req); errs != nil {
		return nil, exception.InvalidArgument(errs)
	}
	result, err := s.userRepo.FindByFilter(ctx, s.db, model.FilterParams{
		{
			Field:    "username",
			Value:    req.Username,
			Operator: "=",
		},
	}, model.OrderParam{
		Order:   "desc",
		OrderBy: "username",
	})
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if result == nil {
		return nil, exception.NotFound("username not found")
	}
	if ok := s.signaturer.CheckBscryptPasswordHash(req.Password, result.Password); !ok {
		return nil, exception.PermissionDenied("username/password unmatched")
	}
	jwtToken, err := s.signaturer.GenerateJWT(result.Id, result.Username)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	return &model.LoginUserRes{
		Username: result.Username,
		Token:    jwtToken,
	}, nil
}
