package main

import (
	"log"
	"net/http"
	"os"

	"github.com/anwardh/meliProject/cmd/server/handler"
	"github.com/anwardh/meliProject/docs"
	"github.com/anwardh/meliProject/internal/products"
	"github.com/anwardh/meliProject/pkg/store"
	"github.com/anwardh/meliProject/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Criação da Função Dummy
// func GetDummyEndpoint(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"hello": "world",
// 	})
// }

// func DummyMiddleware(c *gin.Context) {
// 	log.Println("Im a dummy!")
// 	// Pass on to the next-in-chain
// 	// Depois de imprimir na tela "Im a Dmummy", o request prosseguirá, senão, fica parado (consumindo recurso)
// 	// Cabe a implementação de uma timeOut( ) para encerrar a função
// 	c.Next()
// }

func respondWithError(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, web.NewResponse(code, nil, message))
}

func TokenAuthMiddleware() gin.HandlerFunc {
	requiredToken := os.Getenv("TOKEN")

	// Verificação do token
	if requiredToken == "" { // Se o valor do token estiver vazio
		log.Fatal("por favor, configure a variável de ambiente - token")
	}

	return func(c *gin.Context) {
		token := c.GetHeader("token")

		if token == "" { // Se token que estiver no Header for vazio
			respondWithError(c, http.StatusUnauthorized, "API token obrigatório")
			return
		}

		if token != requiredToken { // Se o token da Header for diferente

			respondWithError(c, http.StatusUnauthorized, "token do API inválido")
			return
		}
		c.Next()
	}
}

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
		pr.Use(TokenAuthMiddleware())

		pr.POST("/", p.Store())
		pr.GET("/", p.GetAll())
		pr.PUT("/:id", p.Update())
		pr.PATCH("/:id", p.UpdateName())
		pr.DELETE("/:id", p.Delete())
	}

	docs.SwaggerInfo.Host = os.Getenv("HOST")
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//r.Use(DummyMiddleware).GET("/dummy", GetDummyEndpoint) // Ednpoint da Função Dummy

	r.Run()
}
