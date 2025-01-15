package permission

import (
	"app/src/general/constants"
	"app/src/general/database"
	"context"
	"github.com/jackc/pgx/v5"
	"strings"
)

const getPermissionsByUserIdQuery = `
SELECT p.id, p.created, p.modified, p.name, p.description
FROM schemaName.permissions p
JOIN schemaName.role_permissions rp ON p.id = rp.permission_id
JOIN schemaName.user_roles ur ON rp.role_id = ur.role_id
WHERE ur.user_id = $1
`

const getPermissionNamesByUserIdQuery = `
SELECT p.name
FROM schemaName.permissions p
JOIN schemaName.role_permissions rp ON p.id = rp.permission_id
JOIN schemaName.user_roles ur ON rp.role_id = ur.role_id
WHERE ur.user_id = $1
`

func GetPermissionsByUserIdDB(schemaName string, userId int64) ([]*Modal, error) {
	replacedSql := strings.ReplaceAll(getPermissionsByUserIdQuery, constants.SchemaName, schemaName)
	rows, err := database.DB.Query(context.Background(), replacedSql, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*Modal
	for rows.Next() {
		permission, err := scanModal(rows)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func GetPermissionNamesByUserIdDB(schemaName string, userId int64) ([]*string, error) {
	replacedSql := strings.ReplaceAll(getPermissionNamesByUserIdQuery, constants.SchemaName, schemaName)

	rows, err := database.DB.Query(context.Background(), replacedSql, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissionNames []*string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		permissionNames = append(permissionNames, &name)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissionNames, nil
}

func scanModal(rows pgx.Rows) (*Modal, error) {
	var modal Modal
	err := rows.Scan(&modal.Id, &modal.Created, &modal.Modified, &modal.Name, &modal.Description)
	if err != nil {
		return nil, err
	}
	return &modal, nil
}
