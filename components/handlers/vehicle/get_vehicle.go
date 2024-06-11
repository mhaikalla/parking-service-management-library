package vehicle

import (
	"log"
	"strconv"

	"parking-service/components/handlers/validator"
	"parking-service/components/models/request"
	"parking-service/pkg/contexts"
	"parking-service/pkg/errs"
)

func (h *Handlers) GetVehicle() func(i interface{}) error {
	return func(i interface{}) error {
		bc := i.(contexts.BearerContext)

		resultValidation, errValidation := validator.ValidateGetListParams(h.Validator, bc)
		if errValidation != nil {
			return bc.JSON(errs.BadRequest, errValidation)
		}

		result, errResp := h.usecaseVehicle.GetVehicles(bc, &request.GetVehicleRequest{
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
