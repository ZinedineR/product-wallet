package service

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/repository"
	"boiler-plate-clean/pkg/exception"
	"boiler-plate-clean/pkg/xvalidator"
	"context"
	"gorm.io/gorm"
)

type ExampleServiceImpl struct {
	db           *gorm.DB
	campaignRepo repository.ExampleRepository
	validate     *xvalidator.Validator
}

func NewExampleService(
	db *gorm.DB, repo repository.ExampleRepository,
	validate *xvalidator.Validator,
) ExampleService {
	return &ExampleServiceImpl{
		db:           db,
		campaignRepo: repo,
		validate:     validate,
	}
}

// CreateExample creates a new campaign
func (s *ExampleServiceImpl) CreateExample(
	ctx context.Context, model *entity.Example,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()

	//txRead := s.db

	//result, err := s.campaignRepo.FindByName(ctx, txRead, "year", strconv.Itoa(model.Year))
	//if err != nil {
	//	return exception.Internal("err", err)
	//}

	//if result != nil {
	//	return exception.AlreadyExists("example already exists")
	//}

	if err := s.campaignRepo.CreateTx(ctx, tx, model); err != nil {
		return exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}
