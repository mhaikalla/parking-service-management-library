package parking

import (
	"strconv"

	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

func (h *Handlers) GetParkingData() func(i interface{}) error {
	return func(i interface{}) error {
		bc := i.(contexts.BearerContext)

		in := request.GetParkingData{}
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

		result, errResp := h.UsecaseParking.GetParkingData(bc, &in)
		if errResp != nil {
			errCode, _ := strconv.Atoi(errResp.Code)
			return bc.JSON(errCode, errResp)
		}

		return bc.JSON(200, result)
	}
}
