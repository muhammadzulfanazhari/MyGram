package controllers

import (
	"MyGram/database"
	"MyGram/helpers"
	"MyGram/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateComment(c *gin.Context) {
	db := database.GetDB()
	Comment := models.Comment{}
	contentType := helpers.GetContentType(c)

	if contentType == "application/json" {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	userData := c.MustGet("userData").(jwt.MapClaims)
	UserID := userData["id"].(float64)

	Comment.UserID = uint(UserID)

	Photo := models.Photo{}

	err := db.Where("id = ?", Comment.PhotoID).First(&Photo).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Data not found",
			"message": "Photo doesn't exist",
		})
		return
	}

	err = db.Create(&Comment).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Comment.ID,
		"message":    Comment.Message,
		"photo_id":   Comment.PhotoID,
		"user_id":    Comment.UserID,
		"created_at": Comment.CreatedAt,
	})
}

func GetComments2(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	userDataId := userData["id"].(float64)
	comments := []models.Comment{}
	err := db.Preload("User", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("ID", "email", "username", "created_at", "updated_at")
	}).Preload("Photo", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("ID", "title", "caption", "photo_url", "user_id", "created_at", "updated_at")
	}).Where("user_id = ?", uint(userDataId)).Find(&comments).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"data": gin.H{
				"error": err.Error(),
				"msg":   "Failed to Get Comment List",
			},
		})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func UpdateComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	Comment := models.Comment{}

	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserID = userID
	Comment.ID = uint(commentId)

	err := db.Model(&Comment).Where("id = ?", commentId).Updates(models.Comment{Message: Comment.Message}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Comment.ID,
		"message":    Comment.Message,
		"user_id":    Comment.UserID,
		"updated_at": Comment.UpdatedAt,
	})
}

func DeleteComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	Comment := models.Comment{}

	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userID := uint(userData["id"].(float64))

	Comment.UserID = userID
	Comment.ID = uint(commentId)

	err := db.Where("id = ?", commentId).Delete(&Comment).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}

func GetComments(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	userDataId := userData["id"].(float64)

	User := models.User{}
	err := db.Debug().Where("id", uint(userDataId)).First(&User).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	Comments := []models.Comment{}

	err = db.Debug().Where("user_id", uint(userDataId)).Find(&Comments).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Comments are not found",
		})
		return
	}

	resComments := []models.ResComment{}
	resUser := models.ResUser{
		ID:       User.ID,
		Username: User.Username,
		Email:    User.Email,
	}

	resComment := models.ResComment{}
	for _, comment := range Comments {
		resComment.ID = comment.ID
		resComment.Message = comment.Message
		resComment.PhotoID = comment.PhotoID
		resComment.UserID = comment.UserID
		resComment.CreatedAt = comment.CreatedAt
		resComment.UpdatedAt = comment.UpdatedAt
		resComment.User = &resUser

		Photo := models.Photo{}
		err = db.Debug().Where("id", comment.PhotoID).First(&Photo).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Photo is not found",
			})
			return
		}
		resPhoto := models.ResPhoto{
			ID:        Photo.ID,
			Title:     Photo.Title,
			Caption:   Photo.Caption,
			Photo_url: Photo.Photo_url,
			UserID:    Photo.UserID,
		}

		resComment.Photo = &resPhoto

		resComments = append(resComments, resComment)
	}

	c.JSON(http.StatusOK, resComments)
}
