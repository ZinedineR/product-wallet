package service

import (
	"context"
	"gorm.io/gorm"
	"product-wallet/internal/model"
	"product-wallet/internal/repository"
	"product-wallet/pkg/utils/converter"
	"product-wallet/pkg/xvalidator"

	//"product-wallet/pkg/exception"
	"product-wallet/pkg/exception"
)

type TransactionServiceImpl struct {
	db                    *gorm.DB
	transactionRepository repository.TransactionRepository
	productRepository     repository.ProductRepository
	walletRepository      repository.WalletRepository
	validate              *xvalidator.Validator
}

func NewTransactionService(
	db *gorm.DB,
	repo repository.TransactionRepository,
	productRepository repository.ProductRepository,
	walletRepository repository.WalletRepository,
	validate *xvalidator.Validator,
) TransactionService {
	return &TransactionServiceImpl{
		db:                    db,
		transactionRepository: repo,
		productRepository:     productRepository,
		walletRepository:      walletRepository,
		validate:              validate,
	}
}

// CreateExample creates a new campaign
func (s *TransactionServiceImpl) Create(
	ctx context.Context, req *model.CreateTransactionReq,
) (*model.CreateTransactionRes, *exception.Exception) {
	tx := s.db.Begin()
	defer tx.Rollback()
	if errs := s.validate.Struct(req); errs != nil {
		return nil, exception.InvalidArgument(errs)
	}
	wallet, err := s.walletRepository.FindByID(ctx, s.db, req.WalletId)
	if err != nil {
		return nil, exception.Internal("failed getting wallet detail", err)
	}
	if wallet == nil {
		return nil, exception.NotFound("wallet detail not found")
	}
	body := req.ToEntity()
	product, err := s.productRepository.FindByID(ctx, s.db, *req.ProductId)
	if err != nil {
		return nil, exception.Internal("error in finding product", err)
	}
	if product == nil {
		return nil, exception.PermissionDenied("product does not exists")
	}
	if !product.Available || product.Quantity < 1 || *req.ProductQuantity > product.Quantity {
		return nil, exception.PermissionDenied("product does not have enough quantity/unavailable")
	}
	totalprice, err := converter.ToFloat64(*req.ProductQuantity)
	if err != nil {
		return nil, exception.Internal("error parsing total price", err)
	}
	totalprice = totalprice * product.Price
	if wallet.Balance < totalprice {
		return nil, exception.PermissionDenied("wallet does not have enough balance to buy this product, balance: " + converter.ToString(wallet.Balance))
	}

	product = req.ToProductEntity(product)
	if err := s.productRepository.UpdateTx(ctx, tx, product); err != nil {
		return nil, exception.Internal("failed updating wallet", err)
	}

	wallet.Decrease(totalprice)
	if err := s.walletRepository.UpdateTx(ctx, tx, wallet); err != nil {
		return nil, exception.Internal("failed updating wallet", err)
	}

	body.Description = "Buying " + product.Name + ", quantity: " + converter.ToString(*req.ProductQuantity) + " for " + converter.ToString(totalprice)
	body.Amount = totalprice
	if err := s.transactionRepository.CreateTx(ctx, tx, body); err != nil {
		return nil, exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return &model.CreateTransactionRes{
		Transaction: *body,
	}, nil
}

func (s *TransactionServiceImpl) Detail(
	ctx context.Context, req *model.GetTransactionByIDReq,
) (*model.GetTransactionByIDRes, *exception.Exception) {
	result, err := s.transactionRepository.FindByID(ctx, s.db, req.ID)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	return &model.GetTransactionByIDRes{
		Transaction: *result,
	}, nil
}

func (s *TransactionServiceImpl) Find(ctx context.Context, req *model.GetAllTransactionReq) (
	*model.GetAllTransactionRes, *exception.Exception,
) {
	result, err := s.transactionRepository.FindByPagination(ctx, s.db, req.Page, req.Sort, req.Filter)
	if err != nil {
		return nil, exception.Internal("err", err)
	}

	return &model.GetAllTransactionRes{
		PaginationData: *result,
	}, nil
}

func (s *TransactionServiceImpl) Credit(
	ctx context.Context, req *model.CreditTransactionReq,
) (*model.CreditTransactionRes, *exception.Exception) {
	tx := s.db.Begin()
	defer tx.Rollback()
	if errs := s.validate.Struct(req); errs != nil {
		return nil, exception.InvalidArgument(errs)
	}
	if req.Amount < 1 {
		return nil, exception.PermissionDenied("Input of amount must be greater than zero")
	}
	wallet, err := s.walletRepository.FindByID(ctx, s.db, req.WalletId)
	if err != nil {
		return nil, exception.Internal("failed getting wallet detail", err)
	}
	if wallet == nil {
		return nil, exception.NotFound("wallet detail not found")
	}

	//category, err := s.categoryRepository.FindByID(ctx, s.db, categoryid)
	//if err != nil {
	//	return exception.Internal("failed getting category detail", err)
	//}
	//if category == nil {
	//	return exception.PermissionDenied("category does not exists")
	//}

	userTransaction := req.ToEntity()
	if err := s.transactionRepository.CreateTx(ctx, tx, userTransaction); err != nil {
		return nil, exception.Internal("failed creating transaction", err)
	}

	wallet.Increase(req.Amount)
	if err := s.walletRepository.UpdateTx(ctx, tx, wallet); err != nil {
		return nil, exception.Internal("failed updating wallet", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return &model.CreditTransactionRes{
		Transaction: *userTransaction,
	}, nil
}

func (s *TransactionServiceImpl) Transfer(
	ctx context.Context, req *model.TransferTransactionReq,
) (*model.TransferTransactionRes, *exception.Exception) {
	tx := s.db.Begin()
	defer tx.Rollback()
	if errs := s.validate.Struct(req); errs != nil {
		return nil, exception.InvalidArgument(errs)
	}
	if req.Amount < 1 {
		return nil, exception.PermissionDenied("Input of amount must be greater than zero")
	}
	sender, err := s.walletRepository.FindByID(ctx, s.db, req.SenderId)
	if err != nil {
		return nil, exception.Internal("failed getting sender detail", err)
	}
	if sender == nil {
		return nil, exception.NotFound("sender wallet detail not found")
	}
	receiver, err := s.walletRepository.FindByID(ctx, s.db, req.ReceiverId)
	if err != nil {
		return nil, exception.Internal("failed getting receiver detail", err)
	}
	if receiver == nil {
		return nil, exception.NotFound("receiver wallet detail not found")
	}
	if sender.Balance < req.Amount {
		return nil, exception.PermissionDenied(sender.Name + " does not have enough balance. Balance: " + converter.ToString(sender.Balance))
	}
	//category, err := s.categoryRepository.FindByName(ctx, s.db, "name", "Transfer")
	//if err != nil {
	//	return exception.Internal("failed getting category detail", err)
	//}
	//if category == nil {
	//	return exception.PermissionDenied("category does not exists")
	//}
	senderTransaction := req.ToSenderEntity(receiver.Name, sender.Id)
	if err := s.transactionRepository.CreateTx(ctx, tx, senderTransaction); err != nil {
		return nil, exception.Internal("failed creating transaction", err)
	}
	sender.Decrease(req.Amount)
	if err := s.walletRepository.UpdateTx(ctx, tx, sender); err != nil {
		return nil, exception.Internal("failed updating wallet", err)
	}
	receiverTransaction := req.ToReceiverEntity(sender.Name, receiver.Id)
	if err := s.transactionRepository.CreateTx(ctx, tx, receiverTransaction); err != nil {
		return nil, exception.Internal("failed creating transaction", err)
	}
	receiver.Increase(req.Amount)
	if err := s.walletRepository.UpdateTx(ctx, tx, receiver); err != nil {
		return nil, exception.Internal("failed updating wallet", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return &model.TransferTransactionRes{
		SenderTransaction:   *senderTransaction,
		ReceiverTransaction: *receiverTransaction,
	}, nil
}

func (s *TransactionServiceImpl) Delete(
	ctx context.Context, req *model.DeleteTransactionReq,
) (*model.DeleteTransactionRes, *exception.Exception) {
	tx := s.db.Begin()
	defer tx.Rollback()

	if err := s.transactionRepository.DeleteByIDTx(ctx, tx, req.ID); err != nil {
		return nil, exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return &model.DeleteTransactionRes{
		ID: req.ID,
	}, nil
}
