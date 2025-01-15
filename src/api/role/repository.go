package role

import (
	"app/src/general/constants"
	"app/src/general/database"
	"app/src/general/util/query"
	"context"
	"github.com/jackc/pgx/v5"
	"strings"
)

const getAllQuery = `SELECT id, created, modified, name, description FROM schemaName.roles ORDER BY id LIMIT $1 OFFSET $2`

const createQuery = `
    INSERT INTO schemaName.roles(name, description) 
    VALUES($1, $2) RETURNING id, created, modified, name, description
`

const updateByIdQuery = `UPDATE schemaName.roles SET name = $1, description = $2 WHERE id = $3 RETURNING id, created, modified, name, description`

const deleteByIdsQuery = `DELETE FROM schemaName.roles WHERE id = ANY($1)`

func GetAllDB(schemaName string, limit, offset int) ([]*Modal, error) {
	replacedSql := strings.ReplaceAll(getAllQuery, constants.SchemaName, schemaName)
	res, err := query.GetAll[Modal](replacedSql, scanModal, limit, offset)
	return res, err
}

func scanModal(rows pgx.Rows) (*Modal, error) {
	var modal Modal
	err := rows.Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Name, &modal.Description)
	if err != nil {
		return nil, err
	}
	return &modal, nil
}

func CreateDB(schemaName string, modal *Modal) error {
	replacedSql := strings.ReplaceAll(createQuery, constants.SchemaName, schemaName)
	err := database.DB.QueryRow(context.Background(), replacedSql, modal.Name, modal.Description).Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Name, &modal.Description)
	return err
}

func UpdateDB(schemaName string, modal *Modal, id string) error {
	replacedSql := strings.ReplaceAll(updateByIdQuery, constants.SchemaName, schemaName)
	err := database.DB.QueryRow(context.Background(), replacedSql, modal.Name, modal.Description, id).Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Name, &modal.Description)
	return err
}

func DeleteDB(schemaName string, modalDeleteRequest *ModalDeleteRequest) error {
	replacedSql := strings.ReplaceAll(deleteByIdsQuery, constants.SchemaName, schemaName)
	_, err := database.DB.Exec(context.Background(), replacedSql, modalDeleteRequest.IdList)
	return err
}
