package main

import (
	"log"
	"os"

	"github.com/anwardh/meliProject/cmd/server/handler"
	"github.com/anwardh/meliProject/docs"
	"github.com/anwardh/meliProject/internal/products"
	"github.com/anwardh/meliProject/pkg/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/*
Instanciamos cada camada do domínio Products e usaremos os métodos do controlador para cada endpoint.
*/
// @title MELI Bootcamp API
// @version 1.0
// @description This API Handle MELI Products.
// @termsOfService https://developers.mercadolibre.com.ar/es_ar/terminos-y-condiciones

// @contact.name API Support
// @contact.url https://developers.mercadolibre.com.ar/support

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {

	/*  Importação do GoDotEnv
	Load lerá seu(s) arquivo(s) env e os carregará no ENV para esse processo.
	Deve ficar sempre no início da aplicação
	Se você chamar Load sem nenhum argumento, o padrão será carregar .env no caminho atual. */
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error ao carregar o arquivo .env")
	}

	// Captura dos valores das Variáveis
	// usuario := os.Getenv("MY_USER")
	// password := os.Getenv("MY_PASS")

	// log.Println("User: ", usuario)
	// log.Println("Password: ", password)
	store := store.Factory("arquivo", "products.json")
	if store == nil {
		log.Fatal("Não foi possivel criar a store")
	}
	repo := products.NewRepository(store)
	service := products.NewService(repo)
	p := handler.NewProduct(service)

	r := gin.Default()
	pr := r.Group("/products")
	{
		pr.POST("/", p.Store())
		pr.GET("/", p.GetAll())
		pr.PUT("/:id", p.Update())
		pr.PATCH("/:id", p.UpdateName())
		pr.DELETE("/:id", p.Delete())
	}

	docs.SwaggerInfo.Host = os.Getenv("HOST")
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}
