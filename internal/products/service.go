package products

// Criação da Interface
type Service interface {
	GetAll() ([]Product, error)
	Store(name, category string, count int, price float64) (Product, error)
	// Declaração do Método Update
	Update(id int, name, productType string, count int, price float64) (Product, error)

	// Declaração do Método UpdateName
	UpdateName(id int, name string) (Product, error)

	// Declaração do Método Delete
	Delete(id int) error
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

// Criação do Método Update
func (s service) Update(id int, name, productType string, count int, price float64) (Product, error) {
	product, err := s.repository.Update(id, name, productType, count, price)

	return product, err
}

// Criação do Método UpdateName
func (s service) UpdateName(id int, name string) (Product, error) {
	product, err := s.repository.UpdateName(id, name)

	return product, err

}

// Criação do Método Delete
func (s service) Delete(id int) error {
	err := s.repository.Delete(id)

	return err
}
