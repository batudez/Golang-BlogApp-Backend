package models

import "gorm.io/gorm"

type Blog struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	ImageUrl    string `json:"imageurl"`
	Slug        string `json:"slug"`
}
