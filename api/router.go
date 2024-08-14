package api

import (
	"api-gateway/api/handler"
	"api-gateway/api/middleware"
	"api-gateway/pkg/config"
	"api-gateway/queue"
	"api-gateway/storage"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	_ "api-gateway/api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router interface {
	InitRouter(queues queue.Queues, enforce *casbin.Enforcer, st storage.Storage)
	Run() error
}

func NewRouter(hd handler.MainHandler, cfg config.Config) Router {
	router := gin.Default()
	return &routerImpl{hd, cfg, router}
}

type routerImpl struct {
	handler.MainHandler
	config.Config
	router *gin.Engine
}

// @title API-Gateway service
// @version 1.0
// @host localhost:8080
// @schemes http
func (r *routerImpl) InitRouter(queues queue.Queues, enforce *casbin.Enforcer, st storage.Storage) {
	r.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.router.Group("/api")

	api.Use(middleware.AuthenticationMiddleware())
	api.Use(middleware.Authorization(enforce))

	{
		user := api.Group("/user")
		usr := r.NewUserHandler()
		user.GET("", usr.GetUserProfile)
		user.PUT("", usr.UpdateUser)
		user.DELETE("", usr.DeleteUser)
		user.GET("/all", usr.GetAllUsers)    //admin
		user.POST("/create", usr.CreateUser) //admin
	}
	{
		medical := api.Group("/medical-records")
		medic := r.NewMedicalRecords(queues.MedicalRecords, st)
		medical.POST("", medic.AddMedicalRecord)
		medical.GET("", medic.ListMedicalRecords) // admin
		medical.GET("/:id", medic.GetMedicalRecord)
		medical.PUT("/:id", medic.UpdateMedicalRecord)
		medical.DELETE("/:id", medic.DeleteMedicalRecord)
	}
	{
		wearable := api.Group("/wearable-data")
		wear := r.NewWearableData(queues.Wearable)
		wearable.POST("", wear.AddWearableData)
		wearable.GET("", wear.ListWearableData) // admin
		wearable.GET("/:id", wear.GetWearableData)
		wearable.PUT("/:id", wear.UpdateWearableData)
		wearable.DELETE("/:id", wear.DeleteWearableData)
	}
	{
		lifestyle := api.Group("/lifestyle")
		life := r.NewLifestyle(queues.Lifestyle)
		lifestyle.POST("", life.AddLifestyleData)
		lifestyle.GET("", life.ListLifestyles) // admin
		lifestyle.GET("/:id", life.GetLifestyleData)
		lifestyle.PUT("/:id", life.UpdateLifestyleData)
		lifestyle.DELETE("/:id", life.DeleteLifestyleData)
	}
	{
		recommend := api.Group("/health-recommendations")
		health := r.NewMedicalRecommend()
		recommend.POST("", health.AddHealthRecommendation)
		recommend.GET("/:id", health.GetHealthRecommendation)
		recommend.GET("", health.GetAllHealthRecommendations)
	}
	{
		monitoring := api.Group("/monitoring")
		mon := r.NewMonitoring()
		monitoring.GET("/realtime", mon.GetMonitoringRealTime)
		monitoring.GET("/daily-summary", mon.GetMonitoringDailySummary)
		monitoring.GET("/weekly-summary", mon.GetMonitoringWeeklySummary)
	}
}

func (r *routerImpl) Run() error {
	return r.router.Run(r.GIN_SERVER_PORT)
}
