package parking

import (
	UsecaseParking "github.com/mhaikalla/parking-service-management-library/components/usecase/usecaseParking"
	"github.com/mhaikalla/parking-service-management-library/pkg/file"

	validation "github.com/go-playground/validator/v10"
)

type Handlers struct {
	Config         map[string]map[string]interface{}
	Validator      validation.Validate
	UsecaseParking UsecaseParking.IUsecaseParking
}

// NewMenuHandlers create a new `MenuHandlers` with `db` provided.
func NewParkingHandlers(
	config map[string]map[string]interface{},
	validator validation.Validate,
	path string,
) (handler *Handlers, err error) {
	defer func() {
		if r, ok := recover().(error); r != nil && ok {
			err = r
		}
	}()

	usecaseParking := UsecaseParking.NewParkingUsecase(file.NewFileSystem(path))

	return &Handlers{
		Config:         config,
		Validator:      validator,
		UsecaseParking: usecaseParking,
	}, nil
}
