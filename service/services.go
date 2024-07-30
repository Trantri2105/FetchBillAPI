package service

import (
	"backend/model"
	"backend/repository"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Service interface {
	GetCamBill(start int64, end int64) ([]model.CamBill, error)
}

type service struct {
	repository repository.Repository
}

func (s service) GetCamBill(start int64, end int64) ([]model.CamBill, error) {
	camBills, err := s.repository.GetCamBill(start, end)
	if err != nil {
		return nil, err
	}
	err = saveBills(camBills, start, end)
	if err != nil {
		log.Println("Service error : " + err.Error())
	}
	return camBills, nil
}

func saveBills(camBills []model.CamBill, start int64, end int64) error {
	file, err := os.Create(fmt.Sprintf("./csv/queryResultFrom%sTo%s.csv", time.Unix(start, 0).UTC().Format("02-01-2006T15_04_05"), time.Unix(end, 0).UTC().Format("02-01-2006T15_04_05")))
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"TransactionID", "PurchaseDateTime", "CameraSN", "PackageType"}
	err = writer.Write(header)
	if err != nil {
		return err
	}
	for _, bill := range camBills {
		row := []string{bill.TransactionID, strconv.Itoa(int(bill.PurchaseDateTime)), bill.CameraSn, bill.PackageType}
		err = writer.Write(row)
		if err != nil {
			return nil
		}
	}
	return nil
}
func NewService(repository repository.Repository) Service {
	return service{repository: repository}
}
