package products

// Adicionando a Estrutura Product e seus campos rotulados
type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Count    int     `json:"count"`
	Price    float64 `json:"price"`
}

/*Criação da variável para guardar os produtos
  Corresponde a nossa Camada de Persistência de Dados */
var ps []Product
var lastID int

// Criação da Iterface e Declaração dos Métodos
type Repository interface {
	GetAll() ([]Product, error)
	Store(id int, name, category string, count int, price float64) (Product, error)
	LastID() (int, error)
}

type repository struct{}

// Função que retornará o repositório um ponteiro para o repositório
func NewRepository() Repository {
	return &repository{}
}

// Métodos que serão utilizados sobre a estrutura repository
// quando for instanciada
func (r *repository) GetAll() ([]Product, error) {
	return ps, nil
}

func (r *repository) LastID() (int, error) {
	return lastID, nil
}

/* Store é o método que salvará as informações do produto,
atribuirá o último ID à variável e retornará a entidade Product */

func (r *repository) Store(id int, name, category string, count int, price float64) (Product, error) {
	p := Product{id, name, category, count, price}
	ps = append(ps, p)
	lastID = p.ID
	return p, nil
}
