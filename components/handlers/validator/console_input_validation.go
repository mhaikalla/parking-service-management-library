package validator

import (
	"strconv"

	"parking-service/components/models/request"
	"parking-service/pkg/contexts"
	"parking-service/pkg/errs"

	validation "github.com/go-playground/validator/v10"
)

func ValidateGetListParams(validatorRequest validation.Validate, bc contexts.BearerContext) (*request.BaseGetListParams, error) {
	limit := bc.QueryParam("limit")
	offset := bc.QueryParam("offset")

	limitVal := 0
	offsetVal := 0
	if len(limit) > 0 {
		l, errLimitVal := strconv.Atoi(limit)
		if errLimitVal != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.BadRequest).
				SetMessage("Invalid Params limit")
		}
		limitVal = l
	}

	if len(offset) > 0 {
		o, errLimitVal := strconv.Atoi(offset)
		if errLimitVal != nil {
			return nil, errs.NewErrContext().
				SetCode(errs.BadRequest).
				SetMessage("Invalid Params limit")
		}
		offsetVal = o
	}

	inputs := request.BaseGetListParams{
		Search: bc.QueryParam("search"),
		Limit:  limitVal,
		Offset: offsetVal,
	}
	errValidate := validatorRequest.Struct(inputs)
	if errValidate != nil {
		return nil, errs.NewErrContext().
			SetCode(errs.BadRequest).
			SetMessage(errs.MessageByCode[strconv.Itoa(errs.BadRequest)])
	}
	return &inputs, nil

}
