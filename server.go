package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rjterradillos/golang_api2/config"
	"github.com/rjterradillos/golang_api2/controller"
	"github.com/rjterradillos/golang_api2/middleware"
	"github.com/rjterradillos/golang_api2/repository"
	"github.com/rjterradillos/golang_api2/service"
	"gorm.io/gorm"
)

var (
	db               *gorm.DB                    = config.SetupDatabaseConnection()
	userRepository   repository.UserRepository   = repository.NewUserRepository(db)
	courseRepository repository.CourseRepository = repository.NewCourseRepository(db)
	jwtService       service.JWTService          = service.NewJWTService()
	userService      service.UserService         = service.NewUserService(userRepository)
	courseService    service.CourseService       = service.NewCourseService(courseRepository)
	authService      service.AuthService         = service.NewAuthService(userRepository)
	authController   controller.AuthController   = controller.NewAuthController(authService, jwtService)
	userController   controller.UserController   = controller.NewUserController(userService, jwtService)
	courseController controller.CourseController = controller.NewCourseController(courseService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	courseRoutes := r.Group("api/courses", middleware.AuthorizeJWT(jwtService))
	{
		courseRoutes.GET("/", courseController.All)
		courseRoutes.POST("/", courseController.Insert)
		courseRoutes.GET("/:id", courseController.FindByID)
		courseRoutes.PUT("/:id", courseController.Update)
		courseRoutes.DELETE("/:id", courseController.Delete)
	}

	r.Run()
}
