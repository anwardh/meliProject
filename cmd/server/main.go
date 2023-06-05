package main

import (
	"github.com/anwardh/meliProject/cmd/server/handler"
	"github.com/anwardh/meliProject/internal/products"
	"github.com/gin-gonic/gin"
)

/*
Instanciamos cada camada do domínio Products e usaremos os métodos do controlador para cada endpoint.
*/

func main() {
	repo := products.NewRepository()
	service := products.NewService(repo)
	p := handler.NewProduct(service)

	r := gin.Default()
	pr := r.Group("/products")
	pr.POST("/", p.Store())
	pr.GET("/", p.GetAll())

	pr.PUT("/:id", p.Update())

	pr.PATCH("/:id", p.UpdateName())

	pr.DELETE("/:id", p.Delete())
	r.Run()
}
