package vehicle

import (
	"log"
	"strconv"

	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

func (h *Handlers) CreateVehicle() func(i interface{}) error {
	return func(i interface{}) error {
		bc := i.(contexts.BearerContext)

		in := request.CreateVehicleRequest{}
		if err := bc.Load(&in); err != nil {
			return bc.JSON(errs.BadRequest, errs.NewErrContext().
				SetCode(errs.BadRequest).
				SetError(err))
		}

		errValidateData := h.Validator.Struct(in)
		if errValidateData != nil {
			return bc.JSON(errs.BadRequest, errs.NewErrContext().
				SetCode(errs.BadRequest).
				SetError(errValidateData))
		}

		result, errResp := h.usecaseVehicle.CreateVehicle(bc, in)
		if errResp != nil {
			log.Println(errResp)
			errCode, _ := strconv.Atoi(errResp.Code)
			return bc.JSON(errCode, errResp)
		}

		return bc.JSON(201, result)
	}
}
