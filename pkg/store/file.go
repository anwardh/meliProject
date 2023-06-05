package store

import (
	"encoding/json"
	"os"
)

// Declaramos a inteface da nossa Store, nela definimos os métodos que a store deve ter
type Store interface {
	Read(data interface{}) error
	Write(data interface{}) error
}

// Estamos declarando um ipo personalizado, que se chama "Type" e o tipo desse tipo é string
// type Type string

// Criamos uma constante definindo no que iremos gravar, no caso falamos que vamos gravar num file ("arquivo")
const (
	// Caso quisermos criar um alias para string, podemos usar o tipo que criamos anterioremnte
	// FileType Type = "file"
	FileType string = "arquivo"
)

// Definimos nossa struct FileStore
type FileStore struct {
	FileName string
}

// A nossa Store é como se fosse um trabalhador que precisa saber o nome do arquivo que ele vai trabalhar
// Ao passarmos para ele o nome do arquivo, poderemos gravar e ler esse arquivo

// Aqui definimos no que o trabalhador vai gravar, no caso num "fiarquivole", e qual o nome desse "arquivo"
func Factory(store string, fileName string) Store {
	switch store {
	case FileType:
		return &FileStore{fileName}
	}
	return nil
}

// Ensinamos ao trabalhador como gravar num arquivo
func (fs *FileStore) Write(data interface{}) error {
	// A função MarshalIndent faz a mesma coisa que a Marshal, porém ela "indenta" o jso também
	fileData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	// Usamos a função WriteFile para escrever no arquivo com o nome que a pessoa definiu
	return os.WriteFile(fs.FileName, fileData, 0644)
}

// Ensinamos ao trabalhador como ler um arquivo
// products []Product
func (fs *FileStore) Read(data interface{}) error {
	// Usamos a função ReadFile para ler o arquivo com o nome que a pessoa definiu
	file, err := os.ReadFile(fs.FileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, data)
}
