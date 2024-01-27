package user

import (
	"net/http"

	"github.com/bloodsteel/easynetes/internal/core"
	"github.com/bloodsteel/easynetes/internal/utils"
)

func ListUsers(userDao core.UserDao) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		users, err := userDao.List(ctx)
		if err != nil {
		}
		utils.RenderSuccess(writer, request, users)
	}
}

func CreateUser(userDao core.UserDao) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func HandlerUser(userDao core.UserDao) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
		case http.MethodPut:
		case http.MethodDelete:
		}
	}
}
