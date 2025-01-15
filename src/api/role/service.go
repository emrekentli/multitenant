package role

func GetAll(schemaName string, limit, offset int) ([]*Modal, error) {
	res, err := GetAllDB(schemaName, limit, offset)
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
