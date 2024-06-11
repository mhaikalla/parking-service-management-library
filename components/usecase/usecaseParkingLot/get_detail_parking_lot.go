package usecaseParkingLot

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

func (ctx *usecaseObj) GetDetailParkingLot(dc contexts.BearerContext, req *request.GetDetailParkingLotRequest) (*response.GetDetailParkingLotResponse, *errs.Errs) {
	resp := response.GetDetailParkingLotResponse{}

	parkingLotData := []models.ParkingLot{}
	resultData := models.ParkingLot{}
	tableName := models.ParkingLotTableName
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
		parkingLot, errloadData := ctx.FileSystem.LoadFile(tableName)
		if errloadData != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errloadData.Error())
		}
		if err := json.Unmarshal(parkingLot, &parkingLotData); err != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(err.Error())
		}
	}
	id, errConv := strconv.Atoi(req.ParkingLotId)
	if errConv != nil {
		return nil, errs.NewErrContext().
			SetCode(errs.InternalServerError).
			SetMessage(errConv.Error())
	}

	for _, pld := range parkingLotData {
		if pld.Id == id {
			if pld.DeletedAt != nil {
				return nil, errs.NewErrContext().
					SetCode(errs.NotFound).
					SetMessage("Data Not Found")
			}
			resultData = pld
			break
		}
	}
	if reflect.ValueOf(resultData).IsZero() {
		return nil, errs.NewErrContext().
			SetCode(errs.NotFound).
			SetMessage("Data Not Found")
	}

	resp = response.GetDetailParkingLotResponse{
		BaseResponse: response.BaseResponse{
			Id:        resultData.Id,
			CreatedAt: resultData.CreatedAt,
			UpdatedAt: resultData.UpdatedAt,
		},
		Name:     resultData.Name,
		Floor:    resultData.Floor,
		IsParked: resultData.IsParked,
	}

	return &resp, nil
}
