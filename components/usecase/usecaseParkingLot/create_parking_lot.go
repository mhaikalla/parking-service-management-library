package usecaseParkingLot

import (
	"encoding/json"
	"time"

	models "github.com/mhaikalla/parking-service-management-library/components/models/entity"
	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/components/models/response"
	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

func (ctx *usecaseObj) CreateParkingLot(dc contexts.BearerContext, req request.CreateParkingLotRequest) (*response.BaseMessageResponse, *errs.Errs) {
	resp := response.BaseMessageResponse{
		Message: "failed",
	}
	parkingLotData := []models.ParkingLot{}

	if !ctx.FileSystem.IsFileExisting(models.ParkingLotTableName) {
		_, errCreate := ctx.FileSystem.CreateFile(models.ParkingLotTableName)
		if errCreate != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errCreate.Error())
		}
	} else {
		parkingLot, errloadData := ctx.FileSystem.LoadFile(models.ParkingLotTableName)
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

	dateNow := time.Now().UTC()

	parkingLotData = append(parkingLotData, models.ParkingLot{
		BaseEntity: models.BaseEntity{
			Id:        len(parkingLotData) + 1,
			CreatedAt: dateNow,
			UpdatedAt: dateNow,
			DeletedAt: nil,
		},
		Floor: req.Floor,
		Name:  req.Name,
	})
	stat, err := ctx.FileSystem.SaveData(models.ParkingLotTableName, parkingLotData)
	if !stat && err != nil {
		return nil, errs.NewErrContext().
			SetCode(errs.InternalServerError).
			SetMessage(err.Error())
	}
	resp.Message = "Success"
	resp.Data = req
	return &resp, nil
}
