package products

// Criação da Interface
type Service interface {
	GetAll() ([]Product, error)
	Store(name, category string, count int, price float64) (Product, error)
}

// Declaração da Estrutura que contém um Repository
type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

/* O método GetAll que se encarregará de passar a tarefa para o Repository e retornar um array de Produtos */
func (s *service) GetAll() ([]Product, error) {
	ps, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	return ps, nil
}

/*
O método Store ficará encarregado de passar a tarefa de obter o último ID e
salvar o produto no Repository, o serviço se encarregará de incrementar o ID

*/

func (s *service) Store(name, category string, count int, price float64) (Product, error) {
	lastID, err := s.repository.LastID()
	if err != nil {
		return Product{}, err
	}

	lastID++

	product, err := s.repository.Store(lastID, name, category, count, price)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}
