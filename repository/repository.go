package repository

import (
	"backend/model"
	"database/sql"
	"errors"
	"log"
)

type Repository interface {
	GetCamBill(start int64, end int64) ([]model.CamBill, error)
}

type repository struct {
	db *sql.DB
}

func (r repository) GetCamBill(start int64, end int64) ([]model.CamBill, error) {
	row, err := r.db.Query(`select transaction_id, purchase_date_time, camera_sn, package_type 
						from cam_bills 
						where transaction_id is not null 
						and payment_method = 'VIETTELPAY' 
						and purchase_date_time > $1 AND purchase_date_time < $2
						and package_type in (select code from package_service where period > 3 and expired > 2595000)`, start, end)
	if err != nil {
		log.Println("Repository GetCamBill error : " + err.Error())
		return nil, errors.New("failed to fetch data")
	}
	camBills := []model.CamBill{}
	for row.Next() {
		bill := model.CamBill{}
		row.Scan(&bill.TransactionID, &bill.PurchaseDateTime, &bill.CameraSn, &bill.PackageType)
		camBills = append(camBills, bill)
	}
	return camBills, nil
}

func NewRepository(db *sql.DB) Repository {
	return repository{db: db}
}
