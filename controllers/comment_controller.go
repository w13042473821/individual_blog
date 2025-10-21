package controllers

import (
	"net/http"
	"strconv"

	"individual_blog/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentController struct {
	DB *gorm.DB
}

func NewCommentController(db *gorm.DB) *CommentController {
	return &CommentController{DB: db}
}

type CreateCommentInput struct {
	Content string `json:"content" binding:"required"`
	PostID  uint   `json:"post_id" binding:"required"`
}

// 创建评论
func (cc *CommentController) CreateComment(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := cc.DB.First(&post, input.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	comment := models.Comment{
		Content: input.Content,
		UserID:  userID.(uint),
		PostID:  input.PostID,
	}

	if err := cc.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}

	// 加载用户信息
	cc.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "评论创建成功",
		"comment": comment,
	})
}

// 获取文章的所有评论
func (cc *CommentController) GetPostComments(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	// 先检查文章是否存在
	var post models.Post
	if err := cc.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	var comments []models.Comment
	// 使用正确的预加载方式
	if err := cc.DB.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, email") // 明确选择需要的字段
		}).
		Where("post_id = ?", postID).
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评论失败: " + err.Error()})
		return
	}

	// 如果评论为空，返回空数组而不是 null
	if comments == nil {
		comments = []models.Comment{}
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}
