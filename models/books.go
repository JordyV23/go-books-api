package models

import "gorm.io/gorm"

//Books es el modelo de libros
type Books struct {
	ID        uint    `gorm:"primary key;autoIncrement" json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:publisher`
}

//MigrateBooks crea la tabla de libros
func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Books{})
	return err
}
