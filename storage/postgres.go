package storage

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Config es la configuracion de la base de datos
type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

//Crea una nueva conexion a la base de datos
func NewConnection(config *Config) (*gorm.DB, error) {
	//Se crea el string de conexion
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)
	//Se crea la conexion
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//Si hay un error al crear la conexion, se retorna un error
	if err != nil {
		return db, err
	}
	//Si no hay error, se retorna la conexion
	return db, nil
}
