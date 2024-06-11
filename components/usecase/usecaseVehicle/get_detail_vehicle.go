package usecaseVehicle

import (
	"encoding/json"
	"reflect"
	"strconv"

	models "github.com/mhaikalla/parking-service-management-library/components/models/entity"
	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/components/models/response"
	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

func (ctx *usecaseObj) GetDetailVehicle(dc contexts.BearerContext, req *request.GetDetailVehicleRequest) (*response.GetDetailVehicleResponse, *errs.Errs) {
	resp := response.GetDetailVehicleResponse{}

	vehicleData := []models.Vehicle{}
	resultData := models.Vehicle{}
	tableName := models.VehicleTableName
	if !ctx.FileSystem.IsFileExisting(tableName) {
		_, errCreate := ctx.FileSystem.CreateFile(tableName)
		if errCreate != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errCreate.Error())
		}
		return nil, errs.NewErrContext().
			SetCode(errs.NotFound).
			SetMessage("Data Not Found")
	} else {
		vehicle, errloadData := ctx.FileSystem.LoadFile(tableName)
		if errloadData != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errloadData.Error())
		}
		if err := json.Unmarshal(vehicle, &vehicleData); err != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(err.Error())
		}
	}
	id, errConv := strconv.Atoi(req.VehicleId)
	if errConv != nil {
		return nil, errs.NewErrContext().
			SetCode(errs.InternalServerError).
			SetMessage(errConv.Error())
	}

	for _, vd := range vehicleData {
		if vd.Id == id {
			if vd.DeletedAt != nil {
				return nil, errs.NewErrContext().
					SetCode(errs.NotFound).
					SetMessage("Data Not Found")
			}
			resultData = vd
			break
		}
	}
	if reflect.ValueOf(resultData).IsZero() {
		return nil, errs.NewErrContext().
			SetCode(errs.NotFound).
			SetMessage("Data Not Found")
	}

	resp = response.GetDetailVehicleResponse{
		BaseResponse: response.BaseResponse{
			Id:        resultData.Id,
			CreatedAt: resultData.CreatedAt,
			UpdatedAt: resultData.UpdatedAt,
		},
		Name:                resultData.Name,
		Type:                resultData.Type,
		FirstHourPrice:      resultData.FirstHourPrice,
		PricePerHourPercent: resultData.PricePerHourPercent,
	}

	return &resp, nil
}
