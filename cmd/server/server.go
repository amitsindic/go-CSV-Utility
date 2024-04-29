package server

import (
	"log"

	"ctlSolution.com/config"
	"ctlSolution.com/internal/handlers"
	"ctlSolution.com/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var(
	handler *handlers.Handler
	dataHandler *handlers.DataHandler
	deleteHandler *handlers.DeleteHandler
	updateHandler *handlers.UpdateHandler
	db *gorm.DB
)

func init(){
	// Initialize MySQL database connection
	db, err := gorm.Open(config.DbName, config.DbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisURL,
		Password: config.RedisPass,
		DB:       config.RedisDB,
	})

	// Initialize new db service
	dbSvc := services.NewDBService(db)

	// Initialize new redis service
	redisSvc := services.NewRedisService(redisClient)

	// Initialize CSV service
	service := services.NewCSVService(dbSvc, redisSvc)

	// Initialize CSV handler with the CSV service dependency injected
	handler = handlers.NewHandler(service)

	// Initialize data handler with the CSV service dependency injected
	dataHandler = handlers.NewDataHandler(dbSvc, redisSvc)

	// Initialize delete handler with the CSV service dependency injected
	deleteHandler = handlers.NewDeleteHandler(dbSvc, redisSvc)

	// Initialize update handler with the CSV service dependency injected
	updateHandler = handlers.NewUpdateHandler(dbSvc, redisSvc)

}


func StartServer(){
	// Initialize Gin router
	r := gin.Default()

	// API endpoint to handle CSV upload
	r.POST("/upload", handler.UploadCSV)

	// API endpoint to retrieve data
	r.GET("/data", dataHandler.GetData)

	// API endpoint to delete by name
	r.DELETE("/data", deleteHandler.DeleteData)

	// API endpoint to patch data
	r.PATCH("/data", updateHandler.UpdateData)

	defer db.Close()

	// Start the server
	r.Run(":8080")
}