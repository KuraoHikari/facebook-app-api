package main

import (
	"github.com/KuraoHikari/facebook-app-res-api/config"
	"github.com/KuraoHikari/facebook-app-res-api/controller"
	"github.com/KuraoHikari/facebook-app-res-api/middleware"
	"github.com/KuraoHikari/facebook-app-res-api/repository"
	"github.com/KuraoHikari/facebook-app-res-api/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = config.SetupDatabaseConnection()
	userRepository 	repository.UserRepository 	= repository.NewUserRepository(db)
	postRepository 	repository.PostRepository 	= repository.NewPostRepository(db)
	jwtService 		service.JWTService 			= service.NewJWTService()
	userService 	service.UserService 		= service.NewUserService(userRepository)
	postService    	service.PostService       	= service.NewBookService(postRepository)
	userController 	controller.UserController 	= controller.NewUserController(userService,jwtService)
	postController 	controller.PostController 	= controller.NewPostController(postService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	
	r := gin.Default()
	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", userController.Login)
		authRoutes.POST("/register", userController.Register)
	}
	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}
	postRoutes := r.Group("api/posts", middleware.AuthorizeJWT(jwtService))
	{
		postRoutes.GET("/", postController.All)
		postRoutes.POST("/", postController.Insert)
		postRoutes.GET("/:id", postController.FindByID)
		postRoutes.PUT("/:id", postController.Update)
		postRoutes.DELETE("/:id", postController.Delete)
	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}