package tag

func GetAll(limit int, offset int) ([]*Modal, error) {
	res, err := GetAllDB(limit, offset)
	return res, err
}

func Create(modal *Modal) (*Modal, error) {
	err := CreateDB(modal)
	return modal, err
}

func Delete(modalDeleteRequest *ModalDeleteRequest) error {
	err := DeleteDB(modalDeleteRequest)
	return err
}
