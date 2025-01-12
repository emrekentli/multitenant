package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/emrekentli/multitenant-boilerplate/src/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ProductRepository struct {
	DB *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) FindAll(ctx context.Context, tenantSchema string) ([]model.Product, error) {
	query := fmt.Sprintf("SELECT id, name, description, price, created_at FROM %s.products", tenantSchema)
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		var description sql.NullString

		if err := rows.Scan(&product.ID, &product.Name, &description, &product.Price, &product.CreatedAt); err != nil {
			return nil, err
		}

		if description.Valid {
			product.Description = &description.String
		} else {
			product.Description = nil
		}

		products = append(products, product)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return products, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, tenantSchema string, id int) (*model.Product, error) {
	query := fmt.Sprintf("SELECT id, name, description, price, created_at FROM %s.products WHERE id = $1", tenantSchema)
	row := r.DB.QueryRow(ctx, query, id)

	var product model.Product
	var description sql.NullString

	if err := row.Scan(&product.ID, &product.Name, &description, &product.Price, &product.CreatedAt); err != nil {
		return nil, err
	}

	if description.Valid {
		product.Description = &description.String
	} else {
		product.Description = nil
	}

	return &product, nil
}

func (r *ProductRepository) Create(ctx context.Context, tenantSchema string, product *model.Product) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.products (name, description, price)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, tenantSchema)
	err := r.DB.QueryRow(ctx, query, product.Name, product.Description, product.Price).Scan(&product.ID, &product.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) Update(ctx context.Context, tenantSchema string, product *model.Product) error {
	query := fmt.Sprintf(`
		UPDATE %s.products
		SET name = $1, description = $2, price = $3
		WHERE id = $4
	`, tenantSchema)
	_, err := r.DB.Exec(ctx, query, product.Name, product.Description, product.Price, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, tenantSchema string, id int) error {
	query := fmt.Sprintf("DELETE FROM %s.products WHERE id = $1", tenantSchema)
	_, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
