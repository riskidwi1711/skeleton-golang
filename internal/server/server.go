package server

import (
	"gitlab.com/skeleton-golang/internal/infrastructure/container"
	"gitlab.com/skeleton-golang/internal/server/http"
)

func StartService(container *container.Container) {
	http.StartH2CServer(container)
}
