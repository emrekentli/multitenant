package permission

func GetPermissionsByUserId(schemaName string, userId int64) ([]*Modal, error) {
	return GetPermissionsByUserIdDB(schemaName, userId)
}

func GetPermissionNamesByUserId(schemaName string, userId int64) ([]*string, error) {
	return GetPermissionNamesByUserIdDB(schemaName, userId)
}
