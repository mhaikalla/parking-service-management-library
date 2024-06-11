package vehicle

import (
	"strconv"

	"parking-service/components/models/request"
	"parking-service/pkg/contexts"
	"parking-service/pkg/errs"
)

func (h *Handlers) UpdateVehicle() func(i interface{}) error {
	return func(i interface{}) error {
		bc := i.(contexts.BearerContext)

		in := request.UpdateVehicleRequest{}
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

		result, errResp := h.usecaseVehicle.UpdateVehicle(bc, in)
		if errResp != nil {
			errCode, _ := strconv.Atoi(errResp.Code)
			return bc.JSON(errCode, errResp)
		}

		return bc.JSON(200, result)
	}
}
