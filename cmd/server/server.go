package server

import (
	"errors"
	"log"
	"net/http"

	"ctlSolution.com/config"
	"ctlSolution.com/internal/handlers"
	"ctlSolution.com/internal/services"
	"github.com/dgrijalva/jwt-go"
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
	r *gin.Engine
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { 
				return nil, errors.New("invalid signing method")
			}
			// Return the secret key used to sign the token
			return []byte(config.SecretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if the token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set the user ID in the context for further use
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		id, ok := claims["id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID in token"})
			c.Abort()
			return
		}
		password, ok := claims["password"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password in token"})
			c.Abort()
			return
		}

		if id != config.AuthID || password != config.AuthPass {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID in token"})
			c.Abort()
			return
		}

		//c.Set("userID", userID)

		c.Next()
	}
}

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

	// Initialize JWT authentication middleware
	authMiddleware := AuthMiddleware()

	// Initialize CSV handler with the CSV service dependency injected
	handler = handlers.NewHandler(service)

	// Initialize data handler with the CSV service dependency injected
	dataHandler = handlers.NewDataHandler(dbSvc, redisSvc)

	// Initialize delete handler with the CSV service dependency injected
	deleteHandler = handlers.NewDeleteHandler(dbSvc, redisSvc)

	// Initialize update handler with the CSV service dependency injected
	updateHandler = handlers.NewUpdateHandler(dbSvc, redisSvc)

	r = gin.Default()

	r.Use(authMiddleware)


}


func StartServer(){
	// Initialize Gin router
	//r := gin.Default()

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
	r.Run(":8081")
}