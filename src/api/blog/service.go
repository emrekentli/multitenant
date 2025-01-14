package blog

func GetAll(tenantSchema string, limit int, offset int, order string) ([]*Modal, error) {
	res, err := GetAllDB(tenantSchema, limit, offset, order)
	return res, err
}

func Create(schemaName string, modal *Modal) (*Modal, error) {
	err := CreateDB(schemaName, modal)
	return modal, err
}

func Update(schemaName string, id string, modal *Modal) error {
	err := UpdateDB(schemaName, modal, id)
	return err
}

func Delete(schemaName string, modalDeleteRequest *ModalDeleteRequest) error {
	err := DeleteDB(schemaName, modalDeleteRequest)
	return err
}
