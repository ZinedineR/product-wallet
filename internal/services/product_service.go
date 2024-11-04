package service

import (
	"context"
	"gorm.io/gorm"
	"product-wallet/internal/model"
	"product-wallet/internal/repository"
	"product-wallet/pkg/exception"
	"product-wallet/pkg/xvalidator"
)

type ProductServiceImpl struct {
	db       *gorm.DB
	repo     repository.ProductRepository
	validate *xvalidator.Validator
}

func NewProductService(
	db *gorm.DB, repo repository.ProductRepository,
	validate *xvalidator.Validator,
) ProductService {
	return &ProductServiceImpl{
		db:       db,
		repo:     repo,
		validate: validate,
	}
}

// CreateExample creates a new campaign
func (s *ProductServiceImpl) Create(
	ctx context.Context, req *model.CreateProductReq,
) (*model.CreateProductRes, *exception.Exception) {
	tx := s.db.Begin()
	defer tx.Rollback()

	duplicateCheck, err := s.repo.FindByFilter(ctx, s.db, model.FilterParams{
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
		return nil, exception.Internal("error finding product", err)
	}
	if duplicateCheck != nil {
		return nil, exception.PermissionDenied("product already exists")
	}

	body := req.ToEntity()

	if err := s.repo.CreateTx(ctx, tx, body); err != nil {
		return nil, exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return &model.CreateProductRes{
		Product: *body,
	}, nil
}

func (s *ProductServiceImpl) Update(
	ctx context.Context, req *model.UpdateProductReq,
) (*model.UpdateProductRes, *exception.Exception) {
	tx := s.db.Begin()
	defer tx.Rollback()
	duplicateCheck, err := s.repo.FindByFilter(ctx, s.db, model.FilterParams{
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
		return nil, exception.Internal("error finding product", err)
	}
	if duplicateCheck != nil && duplicateCheck.Id != req.ID {
		return nil, exception.PermissionDenied("product already exists")
	}
	body := req.ToEntity()
	body.Id = req.ID
	if err := s.repo.UpdateTx(ctx, tx, body); err != nil {
		return nil, exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return &model.UpdateProductRes{
		Product: *body,
	}, nil
}

func (s *ProductServiceImpl) Find(ctx context.Context, req *model.GetAllProductReq) (
	*model.GetAllProductRes, *exception.Exception,
) {
	result, err := s.repo.FindByPagination(ctx, s.db, req.Page, req.Sort, req.Filter)
	if err != nil {
		return nil, exception.Internal("err", err)
	}

	return &model.GetAllProductRes{
		PaginationData: *result,
	}, nil
}

func (s *ProductServiceImpl) Detail(ctx context.Context, req *model.GetProductByIDReq) (
	*model.GetProductByIDRes, *exception.Exception,
) {
	result, err := s.repo.FindByID(ctx, s.db, req.ID)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if result == nil {
		return nil, exception.PermissionDenied("product not found")
	}

	return &model.GetProductByIDRes{
		Product: *result,
	}, nil
}

func (s *ProductServiceImpl) Delete(ctx context.Context, req *model.DeleteProductReq) (
	*model.DeleteProductRes, *exception.Exception,
) {
	tx := s.db.Begin()
	defer tx.Rollback()

	if err := s.repo.DeleteByIDTx(ctx, tx, req.ID); err != nil {
		return nil, exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.Internal("commit transaction", err)
	}
	return &model.DeleteProductRes{
		ID: req.ID,
	}, nil
}
