package query

import (
	"app/src/general/database"
	"app/src/general/util/rest"
	"context"
	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5"
)

func GetAllDBPage[T any](getAllQuery string, scanFunc func(pgx.Rows) (*T, error), total int, args ...any) (*rest.Page[T], error) {
	rows, err := database.DB.Query(context.Background(), getAllQuery, args...)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	items := make([]*T, 0)
	for rows.Next() {
		item, err := scanFunc(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return &rest.Page[T]{
		Size:    len(items),
		Total:   total,
		Content: &items,
	}, nil
}

func GetAll[T any](query string, scanFunc func(pgx.Rows) (*T, error), args ...any) ([]*T, error) {
	rows, err := database.DB.Query(context.Background(), query, args...)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	items := make([]*T, 0)
	for rows.Next() {
		item, err := scanFunc(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func CountQuery(query string, args ...any) (total int, err error) {
	var count int
	err = database.DB.QueryRow(context.Background(), query, args...).Scan(&count)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return count, nil
}
