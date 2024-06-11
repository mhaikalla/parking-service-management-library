package usecaseVehicle

import (
	"github.com/mhaikalla/parking-service-management-library/components/models/request"
	"github.com/mhaikalla/parking-service-management-library/components/models/response"
	"github.com/mhaikalla/parking-service-management-library/pkg/file"

	"github.com/mhaikalla/parking-service-management-library/pkg/contexts"
	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

type IUsecaseVehicle interface {
	CreateVehicle(dc contexts.BearerContext, req request.CreateVehicleRequest) (*response.BaseMessageResponse, *errs.Errs)
	UpdateVehicle(dc contexts.BearerContext, req request.UpdateVehicleRequest) (*response.BaseMessageResponse, *errs.Errs)
	DeleteVehicles(dc contexts.BearerContext, req *request.DeleteVehicleRequest) (*response.BaseMessageResponse, *errs.Errs)
	GetVehicles(dc contexts.BearerContext, req *request.GetVehicleRequest) (*response.GetVehiclesResponse, *errs.Errs)
	GetDetailVehicle(dc contexts.BearerContext, req *request.GetDetailVehicleRequest) (*response.GetDetailVehicleResponse, *errs.Errs)
}

type usecaseObj struct {
	FileSystem file.IFileSystem
}

func NewVehicleUsecase(ctx ...interface{}) IUsecaseVehicle {
	handle := usecaseObj{}
	for _, c := range ctx {
		switch c.(type) {
		case file.IFileSystem:
			handle.FileSystem = c.(file.IFileSystem)
		}
	}
	return &handle
}
