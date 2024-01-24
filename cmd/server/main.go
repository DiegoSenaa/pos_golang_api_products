package main

import (
	"github.com/DiegoSenaa/golang-api/configs"
	_ "github.com/DiegoSenaa/golang-api/docs"
	"github.com/DiegoSenaa/golang-api/internal/entity"
	"github.com/DiegoSenaa/golang-api/internal/infra/database"
	"github.com/DiegoSenaa/golang-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

// @title           Go API Example
// @version         1.0
// @description     Product API with auhtentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Fake
// @contact.email  fake-email@email.com

// @license.name   FAKE
// @license.url    FAKE

// @host      localhost:8000
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	conf := configs.LoadConfig(".")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, entity.User{})

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	c := chi.NewRouter()
	c.Use(middleware.Logger)

	c.Use(middleware.WithValue("jwt", conf.TokenAuth))
	c.Use(middleware.WithValue("jwtExpirationTime", conf.JwtExpirationTime))

	c.Route("/products", func(c chi.Router) {
		c.Use(jwtauth.Verifier(conf.TokenAuth))
		c.Use(jwtauth.Authenticator)
		c.Post("/", productHandler.CreateProduct)
		c.Get("/{id}", productHandler.GetProductById)
		c.Get("/", productHandler.GetAllProducts)
		c.Put("/{id}", productHandler.UpdateProduct)
		c.Delete("/{id}", productHandler.DeleteProduct)
	})
	c.Route("/users", func(u chi.Router) {
		u.Post("/", userHandler.CreateUser)
		u.Get("/", userHandler.GetUserByEmail)
		u.Post("/generate_token", userHandler.GetJwt)
	})

	c.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", c)
}

type ProductHandler struct {
	ProductDB database.ProductInterface
}
