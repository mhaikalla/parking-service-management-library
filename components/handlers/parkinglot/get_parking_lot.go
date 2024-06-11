package parkinglot

import (
	"log"
	"strconv"

	"github.com/mhaikalla/parking-service-management-library/components/handlers/validator"
	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

func (h *Handlers) GetParkingLot() func(i interface{}) error {
	return func(i interface{}) error {
		bc := i.(contexts.BearerContext)

		resultValidation, errValidation := validator.ValidateGetListParams(h.Validator, bc)
		if errValidation != nil {
			return bc.JSON(errs.BadRequest, errValidation)
		}

		result, errResp := h.usecaseParkingLot.GetParkingLots(bc, &request.GetParkingLotRequest{
			BaseGetListParams: request.BaseGetListParams{
				Search: resultValidation.Search,
				Limit:  resultValidation.Limit,
				Offset: resultValidation.Offset,
			},
		})
		if errResp != nil {
			log.Println(errResp)
			errCode, _ := strconv.Atoi(errResp.Code)
			return bc.JSON(errCode, errResp)
		}

		return bc.JSON(200, result)
	}
}
