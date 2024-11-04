package repository

import (
	"context"
	"errors"
	"log/slog"
	"product-wallet/internal/model"
	"product-wallet/pkg/pagination"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommonQuery[T any] interface {
	CreateTx(ctx context.Context, tx *gorm.DB, data *T) error
	UpdateAssociationMany2ManyTx(tx *gorm.DB, data *T) error
	UpdateTx(ctx context.Context, tx *gorm.DB, data *T) error
	UpdateTxWithAssociations(ctx context.Context, tx *gorm.DB, data *T) error
	DeleteByIDTx(ctx context.Context, tx *gorm.DB, id string) error
	Find(
		ctx context.Context, tx *gorm.DB, order model.OrderParam, filter model.FilterParams,
	) (*[]T, error)
	FindByPagination(
		ctx context.Context, tx *gorm.DB, page model.PaginationParam, order model.OrderParam,
		filter model.FilterParams,
	) (*model.PaginationData[T], error)
	FindByID(ctx context.Context, tx *gorm.DB, id string) (*T, error)
	FindByFilter(
		ctx context.Context, tx *gorm.DB, filter model.FilterParams, order model.OrderParam,
	) (*T, error)
}

type Repository[T any] struct {
}

func (r *Repository[T]) CreateTx(ctx context.Context, tx *gorm.DB, data *T) error {
	if err := tx.WithContext(ctx).Omit(clause.Associations).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(data).Error; err != nil {
		slog.Error("failed to create", err)
		return err
	}
	return nil
}

func (r *Repository[T]) UpdateAssociationMany2ManyTx(tx *gorm.DB, data *T) error {
	val := reflect.ValueOf(data).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag.Get("gorm")

		if strings.Contains(tag, "many2many") {
			associationName := typeField.Name
			if err := tx.Model(data).Association(associationName).Replace(field.Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Repository[T]) UpdateTx(ctx context.Context, tx *gorm.DB, data *T) error {
	if err := tx.WithContext(ctx).Omit(clause.Associations).Model(data).Select("*").Updates(data).Error; err != nil {
		slog.Error("failed to update", err)
		return err
	}
	return nil
}

func (r *Repository[T]) UpdateTxWithAssociations(ctx context.Context, tx *gorm.DB, data *T) error {
	if err := tx.WithContext(ctx).Model(data).Select("*").Updates(data).Error; err != nil {
		slog.Error("failed to update", err)
		return err
	}
	return nil
}

func (r *Repository[T]) DeleteByIDTx(ctx context.Context, tx *gorm.DB, id string) error {
	if err := tx.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(new(T)).Error; err != nil {
		slog.Error("failed to delete", err)
		return err
	}
	return nil
}

func (r *Repository[T]) FindByPagination(
	ctx context.Context, tx *gorm.DB, page model.PaginationParam, order model.OrderParam,
	filter model.FilterParams,
) (*model.PaginationData[T], error) {
	query := tx.WithContext(ctx).Omit(clause.Associations)
	query = pagination.Where(filter, query)
	query = pagination.Order(order, query)
	result, err := pagination.Paginate[T](page.Page, page.PageSize, query)
	if err != nil {
		return nil, err
	}
	return &model.PaginationData[T]{
		Page:             result.Page,
		PageSize:         result.PageSize,
		TotalPage:        result.TotalPage,
		TotalDataPerPage: result.TotalDataPerPage,
		TotalData:        result.TotalData,
		Data:             result.Data,
	}, nil
}

func (r *Repository[T]) Find(
	ctx context.Context, tx *gorm.DB, order model.OrderParam, filter model.FilterParams,
) (*[]T, error) {
	var data *[]T
	query := tx.WithContext(ctx).Omit(clause.Associations)
	query = pagination.Where(filter, query)
	query = pagination.Order(order, query)
	if err := query.Find(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("failed to find all", err)
		return nil, err
	}
	return data, nil
}

func (r *Repository[T]) FindByID(ctx context.Context, tx *gorm.DB, id string) (*T, error) {
	var data T
	if err := tx.WithContext(ctx).Preload(clause.Associations).Where("id = ?", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("failed to find by id", err)
		return nil, err
	}
	return &data, nil
}

func (r *Repository[T]) FindByFilter(
	ctx context.Context, tx *gorm.DB, filter model.FilterParams, order model.OrderParam,
) (*T, error) {
	var data T
	query := tx.WithContext(ctx).Preload(clause.Associations)
	query = pagination.Where(filter, query)
	query = pagination.Order(order, query)
	if err := query.First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("failed to find by column", err)
		return nil, err
	}
	return &data, nil
}
