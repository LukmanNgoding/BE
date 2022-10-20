package delivery

import "github.com/ALTA-Group-Project-Social-Media-Apps/Social-Media-Apps/features/user/domain"

func SuccessDelete(msg string) map[string]interface{} {
	return map[string]interface{}{
		"message": msg,
	}
}

func FailResponse(msg string) map[string]string {
	return map[string]string{
		"message": msg,
	}
}

func SuccessResponse(msg string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"message": msg,
		"data":    data,
	}
}
func SuccessLogin(msg string, token string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"message": msg,
		"data":    data,
		"token":   token,
	}
}

type RegisterResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"nama"`
	Email    string `json:"email"`
}

type UpdateResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Photo    string `json:"photo"`
	Bio      string `json:"bio"`
}

func ToResponse(core interface{}, code string) interface{} {
	var res interface{}
	switch code {
	case "reg":
		cnv := core.(domain.Core)
		res = RegisterResponse{ID: cnv.ID, Username: cnv.Username, Email: cnv.Email}
	case "upd":
		cnv := core.(domain.Core)
		res = UpdateResponse{ID: cnv.ID, Username: cnv.Username, Email: cnv.Email, Photo: cnv.Photo, Bio: cnv.Bio}
	case "all":
		var arr []RegisterResponse
		cnv := core.([]domain.Core)
		for _, val := range cnv {
			arr = append(arr, RegisterResponse{ID: val.ID, Username: val.Username, Email: val.Email})
		}
		res = arr
	}

	return res
}
