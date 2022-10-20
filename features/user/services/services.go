package services

import (
	"errors"
	"strings"
	"time"

	"github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/config"
	"github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/features/user/domain"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	qry domain.Repository
}

func New(repo domain.Repository) domain.Service {
	return &userService{
		qry: repo,
	}
}

func (us *userService) Delete(ID uint) error {
	err := us.qry.Delete(ID)
	if err != nil {
		log.Error(err.Error())
		return errors.New("no data")
	}

	return nil
}

func (us *userService) AddUser(newUser domain.Core) (domain.Core, error) {
	var temp string
	if newUser.Password != "" {
		generate, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
		newUser.Password = string(generate)
	}
	if newUser.Password == "" {
		temp = newUser.Password
	} else {
		temp = "error"
	}
	res, err := us.qry.Insert(newUser)

	if err != nil {
		if temp == "" {
			return domain.Core{}, errors.New("cannot encript password")
		}
		return domain.Core{}, errors.New(config.DUPLICATED_DATA)
	}

	return res, nil
}

func (us *userService) UpdateProfile(updatedData domain.Core) (domain.Core, error) {
	if updatedData.Password != "" {
		generate, _ := bcrypt.GenerateFromPassword([]byte(updatedData.Password), 10)
		updatedData.Password = string(generate)
	}

	res, err := us.qry.Update(updatedData)
	if err != nil {
		if strings.Contains(err.Error(), "column") {
			return domain.Core{}, errors.New("rejected from database")
		}
		return domain.Core{}, errors.New("rejected from database")
	}

	return res, nil
}

func (us *userService) LoginUser(newUser domain.Core) (domain.Core, error) {
	res, err := us.qry.Login(newUser)
	if err != nil {
		log.Error(err.Error())
		if strings.Contains(err.Error(), "table") {
			return domain.Core{}, errors.New("database error")
		} else if strings.Contains(err.Error(), "found") {
			return domain.Core{}, errors.New("no data")
		}
	}
	// token := GenerateToken(res.ID)
	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(newUser.Password))
	if err != nil {
		return domain.Core{}, errors.New("password tidak cocok")
	}
	return res, nil

}
func (us *userService) GenerateToken(id uint) string {
	claim := make(jwt.MapClaims)
	claim["authorized"] = true
	claim["id"] = id
	claim["exp"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	str, err := token.SignedString([]byte("R4hs!!a@"))
	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return str
}
