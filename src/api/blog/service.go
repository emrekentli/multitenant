package blog

func GetAll(limit int, offset int, order string) ([]*Modal, error) {
	res, err := GetAllDB(limit, offset, order)
	return res, err
}

func Create(modal *Modal) (*Modal, error) {
	err := CreateDB(modal)
	return modal, err
}

func Update(id string, modal *Modal) error {
	err := UpdateDB(modal, id)
	return err
}

func Delete(modalDeleteRequest *ModalDeleteRequest) error {
	err := DeleteDB(modalDeleteRequest)
	return err
}
