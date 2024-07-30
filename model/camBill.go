package model

type CamBill struct {
	TransactionID    string `json:"transactionId"`
	PurchaseDateTime int64    `json:"purchaseDateTime"`
	CameraSn         string `json:"cameraSn"`
	PackageType      string `json:"packageType"`
}
