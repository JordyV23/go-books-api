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
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

//Funcion para crear un libro
func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}
	err := context.BodyParser(&book)
	//Si hay un error al parsear el body, se retorna un error
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}
	//Se crea el modelo de libro
	err = r.DB.Create(&book).Error
	//Si hay un error al crear el libro, se retorna un error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "No se pudo crear el libro"})
		return err
	}
	//Si no hay error, se retorna un mensaje de exito
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Libro creado exitosamnete"})
	return nil
}

//Funcion para eliminar un libro
func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	//Inicializa el modelo de libro
	bookModel := models.Books{}
	//Obtiene el id del libro a eliminar
	id := context.Params("id")
	//Si el id esta vacio, se retorna un error
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "El ID no puede estar vacio"})
		return nil
	}
	//Se elimina el libro
	err := r.DB.Delete(bookModel, id)
	//Si hay un error al eliminar el libro, se retorna un error
	if err.Error != nil {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "No se pudo eliminar el libro"})
		return err.Error
	}
	//Si no hay error, se retorna un mensaje de exito
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Libro eliminado exitosamente"})
	return nil
}

//Funcion para obtener todos los libros
func (r *Repository) GetBooks(context *fiber.Ctx) error {
	//Inicializa el modelo de libros
	bookModels := &[]models.Books{}
	//Obtiene todos los libros
	err := r.DB.Find(bookModels).Error
	//Si hay un error al obtener los libros, se retorna un error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Lo siento, el libro no pudo ser obtenido"})
		return err
	}
	//Si no hay error, se retorna un mensaje de exito
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Libro obtenido existosamente", "data": bookModels})
	return nil
}

//Funcion para obtener un libro por su id
func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	//Obtiene el id del libro a obtener
	id := context.Params("id")
	//Inicializa el modelo de libro
	bookModel := &models.Books{}
	//Si el id esta vacio, se retorna un error
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "El ID no puede estar vacio"})
		return nil
	}
	//Imprime el id del libro a obtener
	fmt.Println("El id es: ", id)
	//Obtiene el libro por su id
	err := r.DB.Where("id = ?", id).First(bookModel).Error
	//Si hay un error al obtener el libro, se retorna un error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "No se pudo obtener el libro"})
		return err
	}
	//Si no hay error, se retorna un mensaje de exito
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

//Funcion main
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	//Se crea la configuracion inicial
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	//Se crea la conexion a la base de datos
	db, err := storage.NewConnection(config)
	//Si hay un error al crear la conexion, se retorna un error
	if err != nil {
		log.Fatal("No se puede cargar la base de datos")
	}
	//Se migran los modelos
	err = models.MigrateBooks(db)
	//Si hay un error al migrar los modelos, se retorna un error
	if err != nil {
		log.Fatal("No se puede migrar la base de datos")
	}
	//Se crea el repositorio
	r := Repository{
		DB: db,
	}
	//Se crea la aplicacion
	app := fiber.New()
	//Se declaran las rutas
	r.SetupRoutes(app)
	//Se inicia la aplicacion
	app.Listen(os.Getenv("APP_PORT"))
}
