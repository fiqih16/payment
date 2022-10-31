package main

import (
	"api-payment/config"
	"api-payment/controller"
	"api-payment/repository"
	"api-payment/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = config.SetupDatabaseConnection()

	// repository
	userRepository repository.UserRepository = repository.NewUserRepository(db)

	// service
	jwtService service.JWTService = service.NewJWTService()
	authService service.AuthService = service.NewAuthService(userRepository)

	// controller
	authController controller.AuthController = controller.NewAuthController(jwtService, authService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	authRoutes := r.Group("api/v1")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	r.Run()
}