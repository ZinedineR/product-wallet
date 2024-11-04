package service

import (
	"context"
	"gorm.io/gorm"
	"product-wallet/internal/entity"
	"product-wallet/internal/model"
	"product-wallet/internal/repository"
	"product-wallet/pkg/exception"
	"product-wallet/pkg/utils/converter"
	"product-wallet/pkg/xvalidator"
)

type WalletServiceImpl struct {
	db               *gorm.DB
	userRepository   repository.UserRepository
	walletRepository repository.WalletRepository
	transactionRepo  repository.TransactionRepository
	validate         *xvalidator.Validator
}

func NewWalletService(
	db *gorm.DB, repo repository.WalletRepository,
	userRepository repository.UserRepository,
	transactionRepository repository.TransactionRepository,
	validate *xvalidator.Validator,
) WalletService {
	return &WalletServiceImpl{
		db:               db,
		walletRepository: repo,
		userRepository:   userRepository,
		transactionRepo:  transactionRepository,
		validate:         validate,
	}
}

// CreateExample creates a new campaign
func (s *WalletServiceImpl) Create(
	ctx context.Context, req *model.CreateWalletReq,
) (*model.CreateWalletRes, *exception.Exception) {
	tx := s.db.Begin()
	defer tx.Rollback()
	if errs := s.validate.Struct(req); errs != nil {
		return nil, exception.InvalidArgument(errs)
	}
	userCheck, err := s.userRepository.FindByID(ctx, s.db, req.UserId)
	if err != nil {
		return nil, exception.Internal("error finding user", err)
	}
	if userCheck == nil {
		return nil, exception.PermissionDenied("user does not exists")
	}

	duplicateCheck, err := s.walletRepository.FindByFilter(ctx, s.db, model.FilterParams{
		{
			Field:    "name",
			Value:    req.Name,
			Operator: "=",
		},
	}, model.OrderParam{
		Order:   "desc",
		OrderBy: "name",
	})
	if err != nil {
		return nil, exception.Internal("error finding wallet", err)
	}
	if duplicateCheck != nil && duplicateCheck.User.Id == userCheck.Id {
		return nil, exception.PermissionDenied("wallet already exists")
	}

	body := req.ToEntity()

	if err := s.walletRepository.CreateTx(ctx, tx, body); err != nil {
		return nil, exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return &model.CreateWalletRes{
		Wallet: *body,
	}, nil
}

func (s *WalletServiceImpl) Update(
	ctx context.Context, req *model.UpdateWalletReq,
) (*model.UpdateWalletRes, *exception.Exception) {
	tx := s.db.Begin()
	defer tx.Rollback()
	if errs := s.validate.Struct(req); errs != nil {
		return nil, exception.InvalidArgument(errs)
	}
	userCheck, err := s.userRepository.FindByID(ctx, s.db, req.UserId)
	if err != nil {
		return nil, exception.Internal("error finding user", err)
	}
	if userCheck == nil {
		return nil, exception.PermissionDenied("user does not exists")
	}
	duplicateCheck, err := s.walletRepository.FindByFilter(ctx, s.db, model.FilterParams{
		{
			Field:    "name",
			Value:    req.Name,
			Operator: "=",
		},
	}, model.OrderParam{
		Order:   "desc",
		OrderBy: "name",
	})
	if err != nil {
		return nil, exception.Internal("error finding wallet", err)
	}
	if duplicateCheck != nil && duplicateCheck.User.Id == userCheck.Id && duplicateCheck.Id != req.ID {
		return nil, exception.PermissionDenied("wallet already exists")
	}
	body := req.ToEntity()
	body.Id = req.ID
	if err := s.walletRepository.UpdateTx(ctx, tx, body); err != nil {
		return nil, exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return &model.UpdateWalletRes{
		Wallet: *body,
	}, nil
}

func (s *WalletServiceImpl) DetailWalletTransaction(ctx context.Context, req model.GetWalletByTransactionReq) (
	*model.GetWalletByTransactionRes, *exception.Exception,
) {
	wallet, err := s.walletRepository.FindByID(ctx, s.db, req.ID)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if wallet == nil {
		return nil, exception.PermissionDenied("wallet not found")
	}
	filter := model.FilterParams{
		{
			Field:    "wallet_id",
			Value:    wallet.Id,
			Operator: "=",
		},
		{
			Field:    "transaction_time",
			Value:    converter.ToString(req.From),
			Operator: ">=",
		},
		{
			Field:    "transaction_time",
			Value:    converter.ToString(req.To),
			Operator: "<",
		},
	}
	sort := model.OrderParam{
		Order:   "desc",
		OrderBy: "transaction_time",
	}

	transactions, err := s.transactionRepo.Find(ctx, s.db, sort, filter)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	var walletTransactions []entity.Transaction
	if len(*transactions) > 0 {
		for _, transaction := range *transactions {
			walletTransactions = append(walletTransactions, transaction)
		}
	}

	return model.NewGetWalletByTransactionRes(*wallet, walletTransactions), nil
}

func (s *WalletServiceImpl) Find(ctx context.Context, req *model.GetAllWalletReq) (
	*model.GetAllWalletRes, *exception.Exception,
) {
	result, err := s.walletRepository.FindByPagination(ctx, s.db, req.Page, req.Sort, req.Filter)
	if err != nil {
		return nil, exception.Internal("err", err)
	}

	return &model.GetAllWalletRes{
		PaginationData: *result,
	}, nil
}

func (s *WalletServiceImpl) Detail(ctx context.Context, req *model.GetWalletByIDReq) (
	*model.GetWalletByIDRes, *exception.Exception,
) {
	result, err := s.walletRepository.FindByID(ctx, s.db, req.ID)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if result == nil {
		return nil, exception.PermissionDenied("wallet not found")
	}

	return &model.GetWalletByIDRes{
		Wallet: *result,
	}, nil
}

func (s *WalletServiceImpl) Delete(ctx context.Context, req *model.DeleteWalletReq) (
	*model.DeleteWalletRes, *exception.Exception,
) {
	tx := s.db.Begin()
	defer tx.Rollback()

	if err := s.walletRepository.DeleteByIDTx(ctx, tx, req.ID); err != nil {
		return nil, exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return &model.DeleteWalletRes{
		ID: req.ID,
	}, nil
}
