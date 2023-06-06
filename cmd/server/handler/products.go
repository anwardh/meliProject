package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/anwardh/meliProject/internal/products"
	"github.com/anwardh/meliProject/pkg/web"
	"github.com/gin-gonic/gin"
)

// Declaração da Estrutura Request e seus campos rotulados
type request struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Count    int     `json:"count"`
	Price    float64 `json:"price"`
}

// Estrutura Product
type Product struct {
	service products.Service
}

// Função que recebe um Service (do pacote interno) e retorna o controller instanciado
func NewProduct(p products.Service) *Product {
	return &Product{
		service: p,
	}
}

/* O método de obtenção de produtos se encarregará de validar a solicitação,
passar a tarefa ao Service e devolver a resposta correspondente ao cliente */

func (c *Product) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, web.NewResponse(http.StatusUnauthorized, nil, "token inválido"))
			return
		}

		p, err := c.service.GetAll()
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, nil, "não há produtos armazenados"))
			return
		}
		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, p, ""))
	}
}

// Método Store
func (c *Product) Store() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}
		var req request
		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Validação Semêntica dos Campos da nossa Requisição

		if req.Name == "" {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "o nome do produto é obrigatório"))
			return
		}

		if req.Category == "" {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "o nome do tipo é obrigatório"))
			return
		}

		if req.Count == 0 {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "a quantidade é obrigatória"))
			return
		}

		if req.Price == 0 {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, nil, "o preço do produto é obrigatório"))
			return
		}

		p, err := c.service.Store(req.Name, req.Category, req.Count, req.Price)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, p)
	}
}

func (c *Product) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Validação do Token
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}
		// Validação do Id, convertido para inteiro
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		// VALIDAÇÃO DAS ATRIBUIÇÕES DOS CAMPOS DA REQUEST
		/*
			Se algum dos atributos for vazio, o Update não ocorrerá - aqui estão as Regras de Negócio para op Update de um produto
			Este Controller serve, justamente, para que os dados coletados na requisição não sejam, diretamente, armazendos
		no Banco de Dados */

		// Validação da Vinculação dos parâmetros para a Estrutura Request
		var req request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validação do Nome do Produto
		if req.Name == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "O nome do produto é obrigatório"})
			return
		}

		// Validação do Tipo do Produto
		if req.Category == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "O tipo do produto é obrigatório"})
			return
		}

		// Validação da Quantidade do Produto
		if req.Count == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "A quantidade é obrigatória"})
			return
		}

		// Validação do Preço do Produto
		if req.Price == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "O preço é obrigatório"})
			return
		}

		// Quando estiver 'OK', será chamado o método Update, do Service

		p, err := c.service.Update(int(id), req.Name, req.Category, req.Count, req.Price)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return // Retorno do erro do Service
		}
		ctx.JSON(http.StatusOK, p) // Retorno "OK" do Service

	}
}

func (c *Product) UpdateName() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}

		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		var req request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Validação do Nome
		if req.Name == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "o nome do produto é obrigatório"})
			return
		}

		p, err := c.service.UpdateName(int(id), req.Name)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, p)
	}
}

func (c *Product) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}

		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}

		err = c.service.Delete(int(id))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": fmt.Sprintf("O produto %d foi removido", id)})
	}
}
