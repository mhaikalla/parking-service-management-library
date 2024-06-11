package parkinglot

import (
	UsecaseParkingLot "parking-service/components/usecase/usecaseParkingLot"
	"parking-service/pkg/file"

	validation "github.com/go-playground/validator/v10"
)

type Handlers struct {
	Config            map[string]map[string]interface{}
	Validator         validation.Validate
	usecaseParkingLot UsecaseParkingLot.IUsecaseParkingLot
}

func NewParkingLotHandlers(
	config map[string]map[string]interface{},
	validator validation.Validate,
	path string,
) (handler *Handlers, err error) {
	defer func() {
		if r, ok := recover().(error); r != nil && ok {
			err = r
		}
	}()
	usecaseParkingLot := UsecaseParkingLot.NewParkingLotUsecase(file.NewFileSystem(path))
	return &Handlers{
		Config:            config,
		Validator:         validator,
		usecaseParkingLot: usecaseParkingLot,
	}, nil
}
