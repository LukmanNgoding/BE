package domain

import "github.com/labstack/echo/v4"

type Core struct {
	ID       uint
	Username string
	Email    string
	Password string
	Photo    string
	Bio      string
}

type Repository interface {
	Insert(newUser Core) (Core, error)
	Delete(ID uint) error
	Update(updateData Core) (Core, error)
	Login(newUser Core) (Core, error)
}

type Service interface {
	AddUser(newUser Core) (Core, error)
	Delete(id uint) error
	UpdateProfile(updateData Core) (Core, error)
	LoginUser(newUser Core) (Core, error)
	GenerateToken(id uint) string
}

type Handler interface {
	AddUser() echo.HandlerFunc
	DeleteByID() echo.HandlerFunc
	updateUser() echo.HandlerFunc
	LoginUser() echo.HandlerFunc
}
