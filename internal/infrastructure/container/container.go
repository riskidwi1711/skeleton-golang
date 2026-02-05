package container

import (
	"context"
	"os"
	"sync"

	"gitlab.com/skeleton-golang/internal/config"
	"gitlab.com/skeleton-golang/internal/domain/repositories"
	mongodb "gitlab.com/skeleton-golang/internal/infrastructure/mongo"
	"gitlab.com/skeleton-golang/internal/pkg/log"
	"gitlab.com/skeleton-golang/internal/usecase/user"
)

type Container struct {
	Config       *config.DefaultConfig
	PostgresqlDB *config.PostgreSQLDB
	UserService  user.Service
}

func (c *Container) Validate() *Container {
	if c.Config == nil {
		panic("Config is nil")
	}
	if c.UserService == nil {
		panic("UserService is nil")
	}
	return c
}

func New(testingEnv ...string) *Container {
	if len(testingEnv) > 0 {
		config.Load(os.Getenv("env"), testingEnv[0])
	} else {
		config.Load(os.Getenv("env"), ".env")
	}

	defConfig := &config.DefaultConfig{
		Apps: config.Apps{
			Name:     config.GetString("appName"),
			Address:  config.GetString("address"),
			HttpPort: config.GetString("port"),
		},
	}
	mongoCfg := &config.MongoDB{
		URI: config.GetString("sa.mongodb.uri"),
	}

	log.New()

	ctx := context.Background()
	var wg = sync.WaitGroup{}
	mainMongo := mongodb.NewDB(*mongoCfg, ctx, &wg)

	// * Repositories
	userRepo := repositories.NewUser(mainMongo)

	// * Wrapper
	// pantherWrapper := calypso.NewWrapper().SetupRequestHeader().Setup(*pantherConfid)

	// * Services
	userService := user.NewService(
		userRepo,
	)

	// * Brokers

	// * Workers

	container := &Container{
		Config:      defConfig,
		UserService: userService,
	}
	container.Validate()
	return container

}
