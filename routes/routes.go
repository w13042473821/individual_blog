package routes

import (
	"individual_blog/controllers"
	"individual_blog/middleware"
	"individual_blog/models"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	db := models.GetDB()

	// 初始化控制器
	userController := controllers.NewUserController(db)
	postController := controllers.NewPostController(db)       // 文章控制器
	commentController := controllers.NewCommentController(db) // 评论控制器

	// 公开路由
	public := router.Group("/api")
	{
		public.POST("/register", userController.Register)
		public.POST("/login", userController.Login)
		public.GET("/posts", postController.GetPosts)                        // 文章列表
		public.GET("/posts/:id", postController.GetPost)                     // 文章详情
		public.GET("/posts/:id/comments", commentController.GetPostComments) // 文章评论
	}

	// 需要认证的路由
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// 用户相关
		protected.GET("/user", userController.GetCurrentUser)

		// 文章相关
		protected.POST("/posts", postController.CreatePost)
		protected.PUT("/posts/:id", postController.UpdatePost)
		protected.DELETE("/posts/:id", postController.DeletePost)

		// 评论相关
		protected.POST("/comments", commentController.CreateComment)
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
