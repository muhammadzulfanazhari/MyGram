package controllers

import (
	"MyGram/database"
	"MyGram/helpers"
	"MyGram/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreatePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Photo := models.Photo{}
	UserID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	Photo.UserID = UserID

	err := db.Debug().Create(&Photo).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         Photo.ID,
		"title":      Photo.Title,
		"caption":    Photo.Caption,
		"photo_url":  Photo.Photo_url,
		"user_id":    Photo.UserID,
		"created_at": Photo.CreatedAt,
	})
}

func UpdatePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	Photo := models.Photo{}

	photoId, _ := strconv.Atoi(c.Param("photoId"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	Photo.UserID = userID
	Photo.ID = uint(photoId)

	err := db.Model(&Photo).Where("id = ?", photoId).Updates(models.Photo{Title: Photo.Title, Caption: Photo.Caption, Photo_url: Photo.Photo_url}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         Photo.ID,
		"title":      Photo.Title,
		"caption":    Photo.Caption,
		"photo_url":  Photo.Photo_url,
		"user_id":    Photo.UserID,
		"updated_at": Photo.UpdatedAt,
	})
}

func DeletePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	Photo := models.Photo{}

	photoId, _ := strconv.Atoi(c.Param("photoId"))
	userID := uint(userData["id"].(float64))

	Photo.UserID = userID
	Photo.ID = uint(photoId)

	err := db.Where("id = ?", photoId).Delete(&Photo).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})
}

func GetPhotos(c *gin.Context) {
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

	Photos := []models.Photo{}

	err = db.Debug().Where("user_id", uint(userDataId)).Find(&Photos).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Photos are not found",
		})
		return
	}

	resPhotos := []models.ResPhoto{}
	resUser := models.ResUser{
		Username: User.Username,
		Email:    User.Email,
	}

	resPhoto := models.ResPhoto{}
	for _, photo := range Photos {
		resPhoto.ID = photo.ID
		resPhoto.Title = photo.Title
		resPhoto.Caption = photo.Caption
		resPhoto.Photo_url = photo.Photo_url
		resPhoto.UserID = photo.UserID
		resPhoto.CreatedAt = photo.CreatedAt
		resPhoto.UpdatedAt = photo.UpdatedAt
		resPhoto.User = &resUser

		resPhotos = append(resPhotos, resPhoto)
	}

	c.JSON(http.StatusOK, resPhotos)
}
