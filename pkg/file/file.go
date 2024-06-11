package file

type IFileSystem interface {
	CreateFile(nameFile string) (*string, error)
	SaveData(nameFile string, data interface{}) (bool, error)
	IsFileExisting(nameFile string) bool
	LoadFile(fileName string) ([]byte, error)
}
