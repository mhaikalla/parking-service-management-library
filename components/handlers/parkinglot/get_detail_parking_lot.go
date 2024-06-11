package parkinglot

import (
	"log"
	"strconv"

	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

func (h *Handlers) GetDetailParkingLot() func(i interface{}) error {
	return func(i interface{}) error {
		bc := i.(contexts.BearerContext)

		// Parse input request
		in := request.GetDetailParkingLotRequest{
			ParkingLotId: bc.Param("id"),
		}
		errValidateData := h.Validator.Struct(in)
		if errValidateData != nil {
			return bc.JSON(errs.BadRequest, errs.NewErrContext().
				SetCode(errs.BadRequest).
				SetError(errValidateData))
		}
		result, errResp := h.usecaseParkingLot.GetDetailParkingLot(bc, &in)
		if errResp != nil {
			log.Println(errResp)
			errCode, _ := strconv.Atoi(errResp.Code)
			return bc.JSON(errCode, errResp)
		}

		return bc.JSON(200, result)
	}
}
