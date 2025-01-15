package user

import (
	"app/app/middlewares/jwt"
	"app/src/api/permission"
	"app/src/general/util/hash"
	"app/src/general/util/rest"
	"errors"
)

func GetAll(schemaName string, limit int, offset int) (*rest.Page[Modal], error) {
	res, err := GetAllDB(schemaName, limit, offset)
	return res, err
}

func Login(schemaName string, modalLoginRequest *ModalRequest) (*JwtResponse, error) {
	user, err := FindByEmail(schemaName, modalLoginRequest)
	if err != nil {
		return nil, err
	}
	isMatch := hash.Match(modalLoginRequest.Password, user.Password)
	if !isMatch {
		return nil, errors.New("invalid credentials")
	}

	permissions, err := permission.GetPermissionNamesByUserId(schemaName, user.Id)
	if err != nil {
		return nil, err
	}

	jwtStr, err := jwt.CreateJwt(user.Id, permissions)
	if err != nil {
		return nil, err
	}

	return &JwtResponse{Token: jwtStr}, nil
}

func Get(schemaName string, id string) (*Modal, error) {

	res, err := GetByIdDB(schemaName, id)
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
