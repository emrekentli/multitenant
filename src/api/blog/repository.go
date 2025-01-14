package blog

import (
	"app/src/general/constants"
	"app/src/general/database"
	"app/src/general/util/query"
	"context"
	"github.com/jackc/pgx/v5"
	"strings"
)

const deleteByIdsQuery = `
with delete_blog_tag as (
    DELETE FROM schemaName.blog_tag WHERE blog_id = ANY($1::int[])
)
DELETE FROM schemaName.blog WHERE id = ANY($1::int[])
`

const getAllQuery = `
SELECT 
    blog.id,
    blog.created,
    blog.modified,
    body,
    image,
    slug,
    COALESCE(
        JSONB_AGG(json_build_object('id', t.id, 'name', t.name)) FILTER (WHERE t.id IS NOT NULL), '[]'
    ) as tags
FROM schemaName.blog
LEFT JOIN schemaName.blog_tag bt ON blog.id = bt.blog_id
LEFT JOIN schemaName.tag t ON bt.tag_id = t.id
GROUP BY blog.id
ORDER BY $3 LIMIT $1 OFFSET $2
`

const createQuery = `
WITH insert_blog AS (
    INSERT INTO schemaName.blog (body, image, slug) 
    VALUES ($1, $2, $3) 
    RETURNING *
),
insert_blog_tag AS (
    INSERT INTO schemaName.blog_tag (blog_id, tag_id)
    SELECT insert_blog.id, unnest($4::int[]) 
    FROM insert_blog 
    RETURNING *
)
SELECT 
    insert_blog.id,
    insert_blog.created,
    insert_blog.modified,
    insert_blog.slug,
    insert_blog.body,
    insert_blog.image,
    COALESCE(
        JSONB_AGG(json_build_object('id', t.id, 'name', t.name)) 
        FILTER (WHERE t.id IS NOT NULL), 
        '[]'
    ) AS tags
FROM insert_blog
LEFT JOIN schemaName.blog_tag bt 
    ON insert_blog.id = bt.blog_id
LEFT JOIN schemaName.tag t 
    ON bt.tag_id = t.id
GROUP BY 
    insert_blog.id, 
    insert_blog.created, 
    insert_blog.modified, 
    insert_blog.slug, 
    insert_blog.body, 
    insert_blog.image;

`

const updateQuery = `
WITH delete_blog_tag AS (
    DELETE FROM schemaName.blog_tag WHERE blog_id = $5
),
insert_blog_tag AS (
    INSERT INTO schemaName.blog_tag (blog_id, tag_id)
    SELECT $5, unnest($4::int[]) RETURNING *
),
update_blog AS (
    UPDATE schemaName.blog
    SET body = $1, image = $2, slug = $3 WHERE id = $5 RETURNING *
)
SELECT
    update_blog.id,
    update_blog.created,
    update_blog.modified,
    update_blog.body,
    update_blog.image,
    update_blog.slug,
    COALESCE(
        JSONB_AGG(json_build_object('id', t.id, 'name', t.name)) FILTER (WHERE t.id IS NOT NULL), '[]'
    ) AS tags
FROM update_blog
LEFT JOIN schemaName.blog_tag bt ON update_blog.id = bt.blog_id
LEFT JOIN schemaName.tag t ON bt.tag_id = t.id
GROUP BY 
    update_blog.id, 
    update_blog.created, 
    update_blog.modified, 
    update_blog.body, 
    update_blog.image, 
    update_blog.slug;

`

func GetAllDB(tenantSchema string, limit, offset int, order string) ([]*Modal, error) {
	replacedSQL := strings.ReplaceAll(getAllQuery, constants.SchemaName, tenantSchema)
	res, err := query.GetAll[Modal](replacedSQL, scanModal, limit, offset, order)
	return res, err
}

func scanModal(rows pgx.Rows) (*Modal, error) {
	var modal Modal
	err := rows.Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Body, &modal.Image, &modal.Slug, &modal.Tags)
	if err != nil {
		return nil, err
	}
	return &modal, nil
}

func CreateDB(schemaName string, modal *Modal) error {
	var tagIds = make([]int64, len(modal.Tags))
	for i, tag := range modal.Tags {
		tagIds[i] = tag.Id
	}
	replacedSQL := strings.ReplaceAll(createQuery, constants.SchemaName, schemaName)
	err := database.DB.QueryRow(context.Background(), replacedSQL, modal.Body, modal.Image, modal.Slug, tagIds).Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Slug, &modal.Body, &modal.Image, &modal.Tags)
	return err
}

func UpdateDB(schemaName string, modal *Modal, id string) error {
	var tagIds = make([]int64, len(modal.Tags))
	for i, tag := range modal.Tags {
		tagIds[i] = tag.Id
	}
	replacedSQL := strings.ReplaceAll(updateQuery, constants.SchemaName, schemaName)
	err := database.DB.QueryRow(context.Background(), replacedSQL, modal.Body, modal.Image, modal.Slug, tagIds, id).
		Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Body, &modal.Image, &modal.Slug, &modal.Tags)
	return err
}

func DeleteDB(schemaName string, modalDeleteRequest *ModalDeleteRequest) error {
	replacedSQL := strings.ReplaceAll(deleteByIdsQuery, constants.SchemaName, schemaName)
	_, err := database.DB.Exec(context.Background(), replacedSQL, modalDeleteRequest.IdList)
	return err
}
