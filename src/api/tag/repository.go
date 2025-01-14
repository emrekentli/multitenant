package tag

import (
	"app/src/general/database"
	"app/src/general/util/query"
	"context"
	"github.com/jackc/pgx/v5"
	"strings"
)

const deleteByIdsQuery = `
    with delete_blog_tag as (
        delete from schemaName.blog_tag
        where tag_id = ANY($1)
    )
    delete from schemaName.tag
    where id = ANY($1)
`

const getAllQuery = `SELECT * FROM schemaName.tag ORDER BY id LIMIT $1 OFFSET $2`

const createQuery = `INSERT INTO schemaName.tag (name) VALUES ($1) RETURNING *`

func GetAllDB(limit, offset int) ([]*Modal, error) {

	res, err := query.GetAll[Modal](getAllQuery, scanModal, limit, offset)
	return res, err
}

func scanModal(rows pgx.Rows) (*Modal, error) {
	var modal Modal
	err := rows.Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Name)
	if err != nil {
		return nil, err
	}
	return &modal, nil
}

func CreateDB(modal *Modal) error {
	replacedSQL := strings.ReplaceAll(createQuery, "schemaName", "istikbal")
	err := database.DB.QueryRow(context.Background(), replacedSQL, modal.Name).Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Name)
	return err
}

func DeleteDB(modalDeleteRequest *ModalDeleteRequest) error {
	replacedSQL := strings.ReplaceAll(deleteByIdsQuery, "schemaName", "istikbal")
	_, err := database.DB.Exec(context.Background(), replacedSQL, modalDeleteRequest.IdList)
	return err
}
