package router

import (
	"blog_project/internal/handler"
	"blog_project/internal/repository"
	"blog_project/internal/service"
	"blog_project/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitRouter 初始化路由
func InitRouter(db *gorm.DB, JwtSecret string) *gin.Engine {
	//1、组装依赖链
	//用户
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	//分类
	categoryRepo := repository.NewCategoryRepo(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	//标签
	tagRepo := repository.NewTagRepo(db)
	tagService := service.NewTagService(tagRepo)
	tagHandler := handler.NewTagHandler(tagService)

	//2、创建Gin引擎
	r := gin.New()

	//3、挂载全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	//4、健康检查，不需要鉴权
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	//5、api v1路由组
	api := r.Group("/api/v1")
	{
		//----公开路由----
		//公开用户权限
		api.POST("/register", userHandler.RegisterHandler)
		api.POST("/login", userHandler.LoginHandler)
		//公开分类权限
		api.GET("/categories", categoryHandler.GetCategoryListHandler)
		api.GET("/categories/:id", categoryHandler.GetCategoryHandler)
		//公开标签权限
		api.GET("/tags", tagHandler.GetTagListHandler)
		api.GET("/tags/:id", tagHandler.GetTagHandler)
		//鉴权路由
		auth := api.Group("/")
		{
			//挂在JWT中间件
			auth.Use(middleware.JWT(JwtSecret))
			//用户模块
			auth.GET("/user/:id", userHandler.GetUserHandler)
			auth.GET("/users", userHandler.GetUserListHandler)
			auth.POST("/user/:id", userHandler.UpdateHandler)
			auth.DELETE("/user/:id", userHandler.DeleteHandler)
			//分类模块
			auth.POST("/categories", categoryHandler.CreateCategoryHandler)
			auth.PUT("/categories/:id", categoryHandler.UpdateCategoryHandler)
			auth.DELETE("/categories/:id", categoryHandler.DeleteCategoryHandler)
			//标签模块
			auth.POST("/tags", tagHandler.CreateTagHandler)
			auth.PUT("/tags/:id", tagHandler.UpdateTagHandler)
			auth.DELETE("/tags/:id", tagHandler.DeleteTagHandler)
		}
	}

	return r
}
