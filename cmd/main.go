package main

import (
	"api-gateway/api"
	"api-gateway/api/handler"
	"api-gateway/cmd/server"
	"api-gateway/pkg/casbin"
	"api-gateway/pkg/config"
	"api-gateway/pkg/logs"
	"api-gateway/storage"
	redis2 "api-gateway/storage/redis"
	"log"
)

func main() {
	cfg := config.Load()
	logger := logs.InitLogger()

	enforcer, err := casbin.CasbinEnforcer(logger)
	if err != nil {
		logger.Error("error in creating casbin enforce", "error", err.Error())
	}

	redis := redis2.RabbitClient(cfg)
	store := storage.NewStorage(logger, redis)

	serviceManager := server.NewServiceManager(cfg, logger)
	queues := server.NewQueues(serviceManager)

	mainHandler := handler.NewHandler(serviceManager)
	router := api.NewRouter(mainHandler, cfg)

	router.InitRouter(queues, enforcer, store)

	log.Fatal(router.Run())
}
