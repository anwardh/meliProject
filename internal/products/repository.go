package products

import (
	"fmt"

	"github.com/anwardh/meliProject/pkg/store"
)

// Adicionando a Estrutura Product e seus campos rotulados
type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Count    int     `json:"count"`
	Price    float64 `json:"price"`
}

/*
Criação da variável para guardar os produtos

	Corresponde a nossa Camada de Persistência de Dados
*/
var ps []Product
var lastID int

// Criação da Iterface e Declaração dos Métodos
type Repository interface {
	GetAll() ([]Product, error)
	Store(id int, name, category string, count int, price float64) (Product, error)
	LastID() (int, error)
	// Declaração do Método Update - que cuidará de atualizar um dado
	Update(id int, name, productType string, count int, price float64) (Product, error)

	// Declaração do Método UpdateName
	UpdateName(id int, name string) (Product, error)

	// Declaração do Método Delete
	Delete(id int) error
}

type repository struct {
	db store.Store
}

// Função que retornará o repositório um ponteiro para o repositório
func NewRepository(db store.Store) Repository {
	return &repository{
		// Aqui estamos passando o "trabalhador" para a repository, que é do tipo Store
		db: db,
	}
}

// Métodos que serão utilizados sobre a estrutura repository
// quando for instanciada
func (r *repository) GetAll() ([]Product, error) {
	var ps []Product
	// estamos preenchendo a variavel "produtos" com a função read
	err := r.db.Read(&ps)

	// Se ocorrer um erro de leitura
	if err != nil {
		// retornamos o erro
		return nil, err
	}
	// Senão, retornamos os produtos lidos
	return ps, nil
}

func (r *repository) LastID() (int, error) {
	var ps []Product
	if err := r.db.Read(&ps); err != nil {
		// Caso ocorrra um erro na leitura do arquivo,retorna o lastId como 0
		return 0, err
	}

	// Caso ele leia o arquivo, mas o arquivo não tenha produtos, retorna o lastId como 0
	if len(ps) == 0 {
		return 0, nil
	}

	// Aqui obtemos o ultimo produto inserido
	ultimoProduto := ps[len(ps)-1]
	// Aqui retornamos o id do ultimo produto
	return ultimoProduto.ID, nil
}

/* Store é o método que salvará as informações do produto,
atribuirá o último ID à variável e retornará a entidade Product */

// Aqui está a implementação antiga da repository
// func (r *repository) Store(id int, name, category string, count int, price float64) (Product, error) {
// 	p := Product{id, name, category, count, price}
// 	ps = append(ps, p)
// 	lastID = p.ID
// 	return p, nil
// }

// Aqui está a nova implementação da repository

// para gravar num arquivo, precisamos ler o arquivo para pegar os produtos
// que já estavam nele, e adicionar mais um

func (r *repository) Store(id int, name, productType string, count int, price float64) (Product, error) {
	produtos := []Product{}

	// estamos preenchendo a variavel "produtos" com a função read
	r.db.Read(&produtos)

	// Criamos um novo produto com as informações que a pessoa passou na função
	p := Product{id, name, productType, count, price}
	// Agora a variavel produtos tem os produtos que estavam no JSON, mais o produto criado
	produtos = append(produtos, p)
	if err := r.db.Write(produtos); err != nil {
		return Product{}, err
	}
	return p, nil

}

// Criação do Método Update
/*
	Este método enviará os valores contidos nos campos que passarmos em 'p' e os alocará a um elemento, já existente,
encontrado por meio do Id que indicaros na busca (url).
	Com este Id encontrado, todos os elementos dos seus campos serão atualizados, caso contrário, não achando esse Id,
será nos enviada uma mensagem de - Produto não encontrado
*/
func (repository) Update(id int, name, productType string, count int, price float64) (Product, error) {
	p := Product{Name: name, Category: productType, Count: count, Price: price} // Instância de "p" para Update
	updated := false                                                            // Atribuição false para Updated - não foi realizado nenhum update até aqui
	for i := range ps {                                                         // Este For percorrerá a lista dos elementos criados no array para buscar o elemento com o Id que já existe
		if ps[i].ID == id { // Caso encontre esse Id ...
			p.ID = id      // ... o Id do novo produto será o mesmo do já existente (basicamente, o Id que passamos substituirá o já existente, só que são iguais)...
			ps[i] = p      // ... e aqui, irá atualizar (neste Id), todos os valores dos elementos que enviarmos no Put...
			updated = true // ... alterando o seu status para "True"
		}
	}
	if !updated { // Caso não tenha havido esse update, ou seja, se continuar como 'false'...
		return Product{}, fmt.Errorf("produto %d não encontrado", id) // ... nos será enviada uma mensagem de erro
	}
	return p, nil // Retorno do novo produto com um erro do tipo 'nil'
}

// Criação do Método updateName
func (repository) UpdateName(id int, name string) (Product, error) {
	var p Product       // Instância de "p" para UpdateName
	updated := false    // Atribuição false para Updated - não foi realizado nenhum update no Nome até aqui
	for i := range ps { // Este For percorrerá a lista dos elementos criados no array para buscar o elemento com o Id que já existe
		if ps[i].ID == id { // Caso encontre esse Id ...
			ps[i].Name = name // ... o Nome que indicarmos "modificará" o que já existe
			updated = true    // Alteração do sstatus para "true"...
			p = ps[i]         // ... e agora, o produto existente receberá o "novo nome"
		}
	}
	if !updated { // Caso não tenha havido esse update, ou seja, se continuar como 'false'...
		return Product{}, fmt.Errorf("produto %d não encontrado", id) // ... nos será enviada uma mensagem de erro
	}
	return p, nil // Retorno do produto com um novo Nome

}

// Criação do Método Delete
func (repository) Delete(id int) error {
	deleted := false
	var index int
	for i := range ps {
		if ps[i].ID == id {
			index = i
			deleted = true
		}
	}
	if !deleted {
		return fmt.Errorf("produto %d não encontrado", id)
	}
	/*
		Aqui, ps está separando a nossa 'lista de valores contidos' em Repository em duas partes
		Na primeira, estarão os valores do início até o índice (id) que buscamos
		Na segunda, adicionamos o valor 1 ao índice, que fará "pular" para o próximo índice, junto com os demais

		 0  1  2  3  4  5 - > Índices
		[1, 2, 3, 4, 5, 6] -> Valores

		Deletando o nº "3"

		[ :2] = [1, 2] -> Parte 1
		index + 1 -> 2 + 1 = 3
		[3: ] = [4, 5, 6] -> Parte 2

		append = ([1, 2], [4, 5, 6] ...)

		[1, 2, 4, 5, 6] -> FINAL
	*/
	ps = append(ps[:index], ps[index+1:]...)
	return nil
}
