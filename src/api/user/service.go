package user

import (
	"app/src/general/middlewares/jwt"
	"app/src/general/util/rest"
)

func GetAll(limit int, offset int) (*rest.Page[Modal], error) {
	res, err := GetAllDB(limit, offset)
	return res, err
}

func Login(modalLoginRequest *ModalRequest) (*JwtResponse, error) {
	res, err := FindByEmailAndPassword(modalLoginRequest)
	if err != nil {
		return nil, err
	}

	jwtStr, err := jwt.CreateJwt(res.Id)
	if err != nil {
		return nil, err
	}

	return &JwtResponse{Token: jwtStr}, nil
}

func Get(id string) (*Modal, error) {
	res, err := GetByIdDB(id)
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
