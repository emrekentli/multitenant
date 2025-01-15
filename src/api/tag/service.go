package tag

func GetAll(schemaName string, limit int, offset int) ([]*Modal, error) {
	res, err := GetAllDB(schemaName, limit, offset)
	return res, err
}

func Create(schemaName string, modal *Modal) (*Modal, error) {
	err := CreateDB(schemaName, modal)
	return modal, err
}

func Delete(schemaName string, modalDeleteRequest *ModalDeleteRequest) error {
	err := DeleteDB(schemaName, modalDeleteRequest)
	return err
}
