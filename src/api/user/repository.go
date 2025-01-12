package user

import (
	"context"
	"github.com/emrekentli/multitenant-boilerplate/config/database"
	"github.com/emrekentli/multitenant-boilerplate/src/util/query"
	"github.com/emrekentli/multitenant-boilerplate/src/util/rest"
	"github.com/jackc/pgx/v5"
	"strings"
)

const getAllQuery = `SELECT id, created, modified, email, password FROM schemaName.usr ORDER BY id LIMIT $1 OFFSET $2`

const getByIdQuery = `SELECT id, created, modified, email, password FROM schemaName.usr WHERE id = $1`

const createQuery = `
    INSERT INTO schemaName.usr(email, password) 
    VALUES($1, $2) RETURNING id, created, modified, email, password
`

const updateByIdQuery = `UPDATE schemaName.usr SET email = $1, password = $2 WHERE id = $3 RETURNING id, created, modified, email, password`

const countQuery = `SELECT COUNT(*) FROM schemaName.usr`

const findByEmailAndPasswordQuery = `SELECT id, created, modified, email, password FROM schemaName.usr WHERE email = $1 AND password = $2`

const deleteByIdsQuery = `DELETE FROM schemaName.usr WHERE id = ANY($1)`

func GetAllDB(limit, offset int) (*rest.Page[Modal], error) {
	total, err := query.CountQuery(countQuery)
	if err != nil {
		return nil, err
	}

	res, err := query.GetAllDBPage[Modal](getAllQuery, scanModal, total, limit, offset)
	return res, err
}

func scanModal(rows pgx.Rows) (*Modal, error) {
	var modal Modal
	err := rows.Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Email, &modal.Password)
	if err != nil {
		return nil, err
	}
	return &modal, nil
}

func FindByEmailAndPassword(modalLoginRequest *ModalRequest) (*Modal, error) {
	var modal Modal
	replacedSql := strings.ReplaceAll(findByEmailAndPasswordQuery, "schemaName", "istikbal")
	err := database.DB.QueryRow(context.Background(), replacedSql, modalLoginRequest.Email, modalLoginRequest.Password).Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Email, &modal.Password)
	if err != nil {
		return nil, err
	}
	return &modal, nil
}

func GetByIdDB(id string) (*Modal, error) {
	var modal Modal
	replacedSql := strings.ReplaceAll(getByIdQuery, "schemaName", "istikbal")
	err := database.DB.QueryRow(context.Background(), replacedSql, id).Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Email, &modal.Password)
	if err != nil {
		return nil, err
	}
	return &modal, nil
}

func CreateDB(modal *Modal) error {
	replacedSql := strings.ReplaceAll(createQuery, "schemaName", "istikbal")
	err := database.DB.QueryRow(context.Background(), replacedSql, modal.Email, modal.Password).Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Email, &modal.Password)
	return err
}

func UpdateDB(modal *Modal, id string) error {
	replacedSql := strings.ReplaceAll(updateByIdQuery, "schemaName", "istikbal")
	err := database.DB.QueryRow(context.Background(), replacedSql, modal.Email, modal.Password, id).Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Email, &modal.Password)
	return err
}

func DeleteDB(modalDeleteRequest *ModalDeleteRequest) error {
	replacedSql := strings.ReplaceAll(deleteByIdsQuery, "schemaName", "istikbal")
	_, err := database.DB.Exec(context.Background(), replacedSql, modalDeleteRequest.IdList)
	return err
}
