package controllers

import (
	"net/http"
	"strings"
	"unicode"

	"blogapp.com/database"
	"blogapp.com/models"
	"github.com/gin-gonic/gin"
)

// GET ALL
func GetAllBlogs(c *gin.Context) {
	var blogs []models.Blog
	database.DB.Find(&blogs)
	c.JSON(http.StatusOK, blogs)
}

// GET
func GetBlog(c *gin.Context) {
	id := c.Param("id")
	var blog models.Blog
	if err := database.DB.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}
	c.JSON(http.StatusOK, blog)
}

// SLUGIFY AND CREATE
func Slugify(title string) string {
	// Başlığı küçük harfe çevir, boşlukları tire ile değiştir ve özel karakterleri kaldır
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			return r
		}
		return -1
	}, slug)
	return slug
}

func CreateBlog(c *gin.Context) {
	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	blog.Slug = Slugify(blog.Title)
	database.DB.Create(&blog)
	c.JSON(http.StatusOK, blog)
}

// UPDATE
func UpdateBlog(c *gin.Context) {
	id := c.Param("id")
	var blog models.Blog
	if err := database.DB.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&blog)
	c.JSON(http.StatusOK, blog)
}

// DELETE
func DeleteBlog(c *gin.Context) {
	id := c.Param("id")
	var blog models.Blog
	if err := database.DB.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	database.DB.Delete(&blog)
	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted"})
}

// GETBLOG BY SLUG
func GetBlogBySlug(slug string) (models.Blog, error) {
	var blog models.Blog
	err := database.DB.Where("slug = ?", slug).First(&blog).Error
	if err != nil {
		return blog, err
	}
	return blog, nil
}

func GetBlogBySlugHandler(c *gin.Context) {
	slug := c.Param("slug")
	blog, err := GetBlogBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}
	c.JSON(http.StatusOK, blog)
}
