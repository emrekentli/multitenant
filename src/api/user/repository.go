package user

import (
	"app/src/general/constants"
	"app/src/general/database"
	"app/src/general/util/query"
	"app/src/general/util/rest"
	"context"
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

func GetAllDB(schemaName string, limit, offset int) (*rest.Page[Modal], error) {
	replacedCountQuery := strings.ReplaceAll(countQuery, constants.SchemaName, schemaName)
	total, err := query.CountQuery(replacedCountQuery)
	if err != nil {
		return nil, err
	}

	replacedGetAllQuery := strings.ReplaceAll(getAllQuery, constants.SchemaName, schemaName)
	res, err := query.GetAllDBPage[Modal](replacedGetAllQuery, scanModal, total, limit, offset)
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

func FindByEmailAndPassword(schemaName string, modalLoginRequest *ModalRequest) (*Modal, error) {
	var modal Modal
	replacedSql := strings.ReplaceAll(findByEmailAndPasswordQuery, constants.SchemaName, schemaName)
	err := database.DB.QueryRow(context.Background(), replacedSql, modalLoginRequest.Email, modalLoginRequest.Password).Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Email, &modal.Password)
	if err != nil {
		return nil, err
	}
	return &modal, nil
}

func GetByIdDB(schemaName string, id string) (*Modal, error) {
	var modal Modal
	replacedSql := strings.ReplaceAll(getByIdQuery, constants.SchemaName, schemaName)
	err := database.DB.QueryRow(context.Background(), replacedSql, id).Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Email, &modal.Password)
	if err != nil {
		return nil, err
	}
	return &modal, nil
}

func CreateDB(schemaName string, modal *Modal) error {
	replacedSql := strings.ReplaceAll(createQuery, constants.SchemaName, schemaName)
	err := database.DB.QueryRow(context.Background(), replacedSql, modal.Email, modal.Password).Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Email, &modal.Password)
	return err
}

func UpdateDB(schemaName string, modal *Modal, id string) error {
	replacedSql := strings.ReplaceAll(updateByIdQuery, constants.SchemaName, schemaName)
	err := database.DB.QueryRow(context.Background(), replacedSql, modal.Email, modal.Password, id).Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Email, &modal.Password)
	return err
}

func DeleteDB(schemaName string, modalDeleteRequest *ModalDeleteRequest) error {
	replacedSql := strings.ReplaceAll(deleteByIdsQuery, constants.SchemaName, schemaName)
	_, err := database.DB.Exec(context.Background(), replacedSql, modalDeleteRequest.IdList)
	return err
}
