package vehicle

import (
	UsecaseVehicle "parking-service/components/usecase/usecaseVehicle"
	"parking-service/pkg/file"

	validation "github.com/go-playground/validator/v10"
)

type Handlers struct {
	Config         map[string]map[string]interface{}
	Validator      validation.Validate
	usecaseVehicle UsecaseVehicle.IUsecaseVehicle
}

func NewVehicleHandlers(
	config map[string]map[string]interface{},
	validator validation.Validate,
	path string,
) (handler *Handlers, err error) {
	defer func() {
		if r, ok := recover().(error); r != nil && ok {
			err = r
		}
	}()
	usecaseVehicle := UsecaseVehicle.NewVehicleUsecase(file.NewFileSystem(path))
	return &Handlers{
		Config:         config,
		Validator:      validator,
		usecaseVehicle: usecaseVehicle,
	}, nil
}
