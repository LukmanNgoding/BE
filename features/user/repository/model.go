package repository

import (
	"github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/features/user/domain"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
	Photo    string
	Bio      string
}

func FromDomain(du domain.Core) User {
	return User{
		Model:    gorm.Model{ID: du.ID},
		Username: du.Username,
		Email:    du.Email,
		Password: du.Password,
		Photo:    du.Photo,
		Bio:      du.Bio,
	}
}

func ToDomain(u User) domain.Core {
	return domain.Core{

		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Photo:    u.Photo,
		Bio:      u.Bio,
	}
}

func ToDomainArray(au []User) []domain.Core {
	var res []domain.Core
	for _, val := range au {
		res = append(res, domain.Core{ID: val.ID, Username: val.Username, Email: val.Email, Password: val.Password, Photo: val.Photo, Bio: val.Bio})
	}

	return res
}
