package delivery

import (
	"net/http"

	jwt "github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/features/JWT"
	"github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/features/user/domain"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	srv domain.Service
}

func New(e *echo.Echo, srv domain.Service) {
	handler := userHandler{srv: srv}
	e.DELETE("/users", handler.DeleteByID())
	e.POST("/users", handler.AddUser())
	e.GET("/users/update", handler.updateUser())
}

func (us *userHandler) DeleteByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := jwt.ExtractToken(c)
		if id == 0 {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "cannot validate token",
			})
		}
		err := us.srv.Delete(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
		}

		return c.JSON(http.StatusOK, SuccessDelete("success delete user"))
	}
}

func (us *userHandler) AddUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input RegisterFormat
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, FailResponse("cannot bind input"))
		}
		cnv := ToDomain(input)
		res, err := us.srv.AddUser(cnv)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
		}

		return c.JSON(http.StatusCreated, SuccessResponse("berhasil register", ToResponse(res, "reg")))
	}

}

func (us *userHandler) updateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input UpdateFormat
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, FailResponse("cannot bind input"))
		}
		cnv := ToDomain(input)
		res, err := us.srv.UpdateProfile(cnv)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
		}

		return c.JSON(http.StatusCreated, SuccessResponse("berhasil update", ToResponse(res, "reg")))
	}

}

func (us *userHandler) LoginUser() echo.HandlerFunc {
	//autentikasi user login
	return func(c echo.Context) error {
		var resQry LoginFormat
		if err := c.Bind(&resQry); err != nil {
			return c.JSON(http.StatusBadRequest, FailResponse("cannot bind input"))
		}

		cnv := ToDomain(resQry)
		res, err := us.srv.LoginUser(cnv)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
		}
		token := us.srv.GenerateToken(res.ID)
		return c.JSON(http.StatusCreated, SuccessLogin("berhasil register", token, ToResponse(res, "reg")))
	}
}
