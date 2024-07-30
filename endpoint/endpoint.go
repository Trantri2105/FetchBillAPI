package endpoints

import (
	"backend/dto"
	"backend/model"
	"backend/service"
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-playground/validator/v10"
)

type BillEndpoint interface {
	GetCamBillEndpoint() endpoint.Endpoint
}

type billEndpoint struct {
	service service.Service
}

func (b billEndpoint) GetCamBillEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(dto.Request)
		validate := validator.New()
		err := validate.Struct(r)
		if err != nil {
			return nil, errors.New("start or end date or time zone is missing")
		}
		start, err := time.Parse("02-01-2006 15:04:05", r.Start)
		if err != nil {
			log.Println("Endpoint error : " + err.Error())
			return nil, err
		}
		var timeZone int
		timeZone, err = strconv.Atoi(r.TimeZone)
		if err != nil {
			log.Println("Endpoint error : " + err.Error())
			return nil, err
		}
		timeZone = -timeZone
		start = start.Add(time.Duration(timeZone) * time.Hour)
		var end time.Time
		end, err = time.Parse("02-01-2006 15:04:05", r.End)
		if err != nil {
			log.Println("Endpoint error : " + err.Error())
			return nil, err
		}
		end = end.Add(time.Duration(timeZone) * time.Hour)
		var camBills []model.CamBill
		camBills, err = b.service.GetCamBill(start.Unix(), end.Unix())
		if err != nil {
			return nil, err
		}
		return dto.Response{CamBills: camBills}, nil
	}
}

func NewBillEndpoint(service service.Service) BillEndpoint {
	return billEndpoint{service: service}
}
