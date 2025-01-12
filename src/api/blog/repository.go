package blog

import (
	"context"
	"github.com/emrekentli/multitenant-boilerplate/config/database"
	"github.com/emrekentli/multitenant-boilerplate/src/util/query"
	"github.com/jackc/pgx/v5"
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
with insert_blog as (
    INSERT INTO schemaName.blog (body, image, slug) VALUES ($1, $2, $3) RETURNING *
),
insert_blog_tag as (
    INSERT INTO schemaName.blog_tag (blog_id, tag_id)
    SELECT insert_blog.id, unnest($4::int[]) FROM insert_blog RETURNING *
)
SELECT 
    insert_blog.id,
    insert_blog.created,
    insert_blog.modified,
    insert_blog.slug,
    body,
    image,
    COALESCE(
        JSONB_AGG(json_build_object('id', t.id, 'name', t.name)) FILTER (WHERE t.id IS NOT NULL), '[]'
    ) as tags
FROM insert_blog
LEFT JOIN schemaName.blog_tag bt ON insert_blog.id = bt.blog_id
LEFT JOIN schemaName.tag t ON bt.tag_id = t.id
GROUP BY insert_blog.id
`

const updateQuery = `
with delete_blog_tag as (
    DELETE FROM schemaName.blog_tag WHERE blog_id = $5
),
insert_blog_tag as (
    INSERT INTO schemaName.blog_tag (blog_id, tag_id)
    SELECT $5, unnest($4::int[]) RETURNING *
),
update_blog as (
    UPDATE schemaName.blog
    SET body = $1, image = $2, slug = $3 WHERE id = $5 RETURNING *
)
SELECT
    update_blog.id,
    update_blog.created,
    update_blog.modified,
    body,
    image,
    slug,
    COALESCE(
        JSONB_AGG(json_build_object('id', t.id, 'name', t.name)) FILTER (WHERE t.id IS NOT NULL), '[]'
    ) as tags
FROM update_blog
LEFT JOIN schemaName.blog_tag bt ON update_blog.id = bt.blog_id
LEFT JOIN schemaName.tag t ON bt.tag_id = t.id
GROUP BY update_blog.id
`

func GetAllDB(limit, offset int, order string) ([]*Modal, error) {
	res, err := query.GetAll[Modal](getAllQuery, scanModal, limit, offset, order)
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

func CreateDB(modal *Modal) error {
	var tagIds = make([]int64, len(modal.Tags))
	for i, tag := range modal.Tags {
		tagIds[i] = tag.Id
	}

	err := database.DB.QueryRow(context.Background(), createQuery, modal.Body, modal.Image, modal.Slug, tagIds).Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Body, &modal.Image, &modal.Slug, &modal.Tags)
	return err
}

func UpdateDB(modal *Modal, id string) error {
	var tagIds = make([]int64, len(modal.Tags))
	for i, tag := range modal.Tags {
		tagIds[i] = tag.Id
	}

	err := database.DB.QueryRow(context.Background(), updateQuery, modal.Body, modal.Image, modal.Slug, tagIds, id).
		Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Body, &modal.Image, &modal.Slug, &modal.Tags)
	return err
}

func DeleteDB(modalDeleteRequest *ModalDeleteRequest) error {
	_, err := database.DB.Exec(context.Background(), deleteByIdsQuery, modalDeleteRequest.IdList)
	return err
}
