package model

// File представляет собой модель файла
type File struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Path       string `json:"path"`
	DocumentID uint   `gorm:"index"` // Внешний ключ для связи с Document
}
