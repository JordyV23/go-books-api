package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/jordyv23/go-books-api/models"
	"github.com/jordyv23/go-books-api/storage"
	"gorm.io/gorm"
)

type Book struct {
	Author  string `json:"author"`
	Title   string `json:"title"`
	Publish string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

//Funcion para crear un libro
func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}
	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "No se pudo crear el libro"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Libro creado exitosamnete"})

	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "El ID no puede estar vacio"})
		return nil
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "No se pudo eliminar el libro"})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Libro eliminado exitosamente"})
	return nil
}

//Funcion para obtener todos los libros
func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Lo siento, el libro no pudo ser obtenido"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Libro obtenido existosamente", "data": bookModels})

	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Books{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "El ID no puede estar vacio"})
		return nil
	}

	fmt.Println("El id es: ", id)

	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "No se pudo obtener el libro"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Libro obtenido exitosamente", "data": bookModel})
	return nil
}

//Funcion para declarar las rutas de la API
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("No se puede cargar la base de datos")
	}

	err = models.MigrateBooks(db)

	if err != nil {
		log.Fatal("No se puede migrar la base de datos")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
