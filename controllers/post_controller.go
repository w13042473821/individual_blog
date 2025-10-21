package controllers

import (
	"net/http"
	"time"

	"individual_blog/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostController struct {
	DB *gorm.DB
}

func NewPostController(db *gorm.DB) *PostController {
	return &PostController{DB: db}
}

type CreatePostInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// 创建文章
func (pc *PostController) CreatePost(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userID.(uint),
	}

	if err := pc.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	// 加载用户信息
	pc.DB.Preload("User").First(&post, post.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "文章创建成功",
		"post":    post,
	})
}

// 获取所有文章
func (pc *PostController) GetPosts(c *gin.Context) {
	var posts []models.Post
	if err := pc.DB.Preload("User").Order("created_at DESC").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	// 创建响应结构，包含评论计数
	type PostResponse struct {
		ID           uint        `json:"id"`
		Title        string      `json:"title"`
		Content      string      `json:"content"`
		UserID       uint        `json:"user_id"`
		User         models.User `json:"user"`
		CommentCount int64       `json:"comment_count"`
		CreatedAt    time.Time   `json:"created_at"`
		UpdatedAt    time.Time   `json:"updated_at"`
	}

	var response []PostResponse
	for _, post := range posts {
		var commentCount int64
		// 查询该文章的评论数量
		pc.DB.Model(&models.Comment{}).Where("post_id = ?", post.ID).Count(&commentCount)

		response = append(response, PostResponse{
			ID:           post.ID,
			Title:        post.Title,
			Content:      post.Content,
			UserID:       post.UserID,
			User:         post.User,
			CommentCount: commentCount,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": response,
	})
}

// 获取单个文章
func (pc *PostController) GetPost(c *gin.Context) {
	var post models.Post
	if err := pc.DB.Preload("User").Preload("Comments.User").First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

// 更新文章
func (pc *PostController) UpdatePost(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var post models.Post
	if err := pc.DB.First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权更新此文章"})
		return
	}

	var input UpdatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Content != "" {
		updates["content"] = input.Content
	}

	if err := pc.DB.Model(&post).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	pc.DB.Preload("User").First(&post, post.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "文章更新成功",
		"post":    post,
	})
}

// 删除文章
func (pc *PostController) DeletePost(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var post models.Post
	if err := pc.DB.First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此文章"})
		return
	}

	if err := pc.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文章删除成功",
	})
}
