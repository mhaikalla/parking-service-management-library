package usecaseParkingLot

import (
	"encoding/json"
	"strconv"
	"time"

	models "github.com/mhaikalla/parking-service-management-library/components/models/entity"
	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/components/models/response"
	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

func (ctx *usecaseObj) UpdateParkingLot(dc contexts.BearerContext, req request.UpdateParkingLotRequest) (*response.BaseMessageResponse, *errs.Errs) {
	resp := response.BaseMessageResponse{
		Message: "failed",
	}
	dateNow := time.Now().UTC()
	parkingLotData := []models.ParkingLot{}
	tableName := models.ParkingLotTableName
	if !ctx.FileSystem.IsFileExisting(tableName) {
		_, errCreate := ctx.FileSystem.CreateFile(tableName)
		if errCreate != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.InternalServerError).
				SetMessage(errCreate.Error())
		}
		parkingLotData = append(parkingLotData, models.ParkingLot{
			BaseEntity: models.BaseEntity{
				Id:        len(parkingLotData) + 1,
				CreatedAt: dateNow,
				UpdatedAt: dateNow,
				DeletedAt: nil,
			},
			Floor:    req.Floor,
			Name:     req.Name,
			IsParked: false,
		})
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
		for i, pl := range parkingLotData {
			id, errConv := strconv.Atoi(req.Id)
			if errConv != nil {
				return nil, errs.NewErrContext().
					SetCode(errs.InternalServerError).
					SetMessage(errConv.Error())
			}

			if pl.Id == id {
				if pl.DeletedAt != nil {
					return nil, errs.NewErrContext().
						SetCode(errs.NotFound).
						SetMessage("Data Not Found")
				}
				if pl.IsParked {
					return nil, errs.NewErrContext().
						SetCode(errs.NotFound).
						SetMessage("This Parking Area Has Filled")
				}

				parkingLotData[i].UpdatedAt = dateNow
				parkingLotData[i].Floor = req.Floor
				parkingLotData[i].Name = req.Name
				break
			}
		}
	}
	stat, err := ctx.FileSystem.SaveData(tableName, parkingLotData)
	if !stat && err != nil {
		return nil, errs.NewErrContext().
			SetCode(errs.InternalServerError).
			SetMessage(err.Error())
	}
	resp.Message = "Success, Data Updated"
	resp.Data = req
	return &resp, nil
}
