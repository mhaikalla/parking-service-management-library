package parkinglot

import (
	UsecaseParkingLot "github.com/mhaikalla/parking-service-management-library/components/usecase/usecaseParkingLot"
	"github.com/mhaikalla/parking-service-management-library/pkg/file"

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
