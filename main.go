package main

import (
	"log"
	"os"

	"individual_blog/config"
	"individual_blog/middleware"
	"individual_blog/models"
	"individual_blog/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化数据库
	db = models.InitDB()

	// 自动迁移数据库表
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	// 设置Gin模式
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin路由
	router := gin.Default()

	// 添加中间件
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// 设置路由
	routes.SetupRoutes(router)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在端口 %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
