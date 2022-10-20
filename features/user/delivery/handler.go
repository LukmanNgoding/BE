package delivery

import (
	"context"
	"net/http"
	"os"
	"strings"

	jwt "github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/features/JWT"
	"github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/features/user/domain"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type userHandler struct {
	srv domain.Service
}

func New(e *echo.Echo, srv domain.Service) {
	handler := userHandler{srv: srv}
	e.POST("/users", handler.AddUser())
	e.POST("/login", handler.LoginUser())
	e.PUT("/users/update", handler.updateUser(), middleware.JWT([]byte("R4hs!!a@")))
	e.DELETE("/users", handler.DeleteByID(), middleware.JWT([]byte("R4hs!!a@")))
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
		input.Username = c.FormValue("username")
		input.Email = c.FormValue("email")
		input.Password = c.FormValue("password")
		input.Bio = c.FormValue("bio")

		file, err := c.FormFile("photo")
		if err != nil {
			return err
		}

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		s3Config := &aws.Config{
			Region:      aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_USER"), os.Getenv("AWS_KEY"), ""),
		}

		s3Session := session.New(s3Config)

		uploader := s3manager.NewUploader(s3Session)
		inputData := &s3manager.UploadInput{
			Bucket: aws.String("sosialtabucket"),           // bucket's name
			Key:    aws.String("myfiles/" + file.Filename), // files destination location
			Body:   src,                                    // content of the file
			// ACL:    aws.String("Objects - List"),
			// ContentType: aws.String("image/png"), // content type
		}
		_, _ = uploader.UploadWithContext(context.Background(), inputData)
		input.Photo = "" + strings.ReplaceAll(file.Filename, " ", "+")
		input.ID = jwt.ExtractToken(c)
		if input.ID == 0 {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "cannot validate token",
			})
		}

		cnv := ToDomain(input)
		res, err := us.srv.UpdateProfile(cnv)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, FailResponse(err.Error()))
		}

		return c.JSON(http.StatusCreated, SuccessResponse("berhasil update", ToResponse(res, "upd")))
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
