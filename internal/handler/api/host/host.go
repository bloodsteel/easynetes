package host

import (
	"net/http"

	"github.com/bloodsteel/easynetes/internal/core"
)

func ListHosts(hostDao core.HostInstanceDao) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
	}
}

func CreateHost(hostDao core.HostInstanceDao) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}

func HandlerHost(hostDao core.HostInstanceDao) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
