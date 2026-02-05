package cmd

import (
	"gitlab.com/skeleton-golang/internal/infrastructure/container"
	"gitlab.com/skeleton-golang/internal/server"
)

func Run() {
	server.StartService(container.New())
}
