package service

import (
	"backend/model"
	"backend/repository"
	"bytes"
	"encoding/csv"
	"fmt"
	"io/fs"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
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
	go saveBillsAndSendViaMail(camBills, start, end)
	return camBills, nil
}

func saveBillsAndSendViaMail(camBills []model.CamBill, start int64, end int64) {
	filePath := fmt.Sprintf("./csv/bills_%s-%s.csv", time.Unix(start, 0).Format("20060102_150405"), time.Unix(end, 0).Format("20060102_150405"))
	file, err := os.Create(filePath)
	if err != nil {
		log.Println(err.Error())
		return
	}

	writer := csv.NewWriter(file)
	header := []string{"TransactionID", "PurchaseDateTime", "CameraSN", "PackageType"}
	err = writer.Write(header)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, bill := range camBills {
		row := []string{bill.TransactionID, strconv.Itoa(int(bill.PurchaseDateTime)), bill.CameraSn, bill.PackageType}
		err = writer.Write(row)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	writer.Flush()
	file.Close()

	subject := "CSV file"
	body := fmt.Sprintf("Camera bill from %s to %s", time.Unix(start, 0).Format("2006-01-02 15:04:05"), time.Unix(end, 0).Format("2006-01-02 15:04:05"))
	to := []string{os.Getenv("TO")}
	err = sendCSVFileViaMail(subject, body, filePath, to)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func sendCSVFileViaMail(subject string, body string, filePath string, to []string) error {
	password := os.Getenv("PASSWORD")
	from := os.Getenv("FROM")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var fileInfo fs.FileInfo
	fileInfo, err = file.Stat()
	if err != nil {
		return err
	}
	fileContent := make([]byte, fileInfo.Size())
	file.Read(fileContent)
	defer file.Close()

	var email bytes.Buffer
	writer := multipart.NewWriter(&email)

	//Email header
	email.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	email.WriteString("MIME-Version: 1.0\r\n")
	email.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary()))
	email.WriteString("\r\n")

	//Email body
	bodyHeader := textproto.MIMEHeader{}
	bodyHeader.Set("Content-Type", "text/plain")
	bodyPart, err := writer.CreatePart(bodyHeader)
	if err != nil {
		return err
	}
	bodyPart.Write([]byte(body))

	//Email attachment
	attachmentHeader := textproto.MIMEHeader{}
	attachmentHeader.Set("Content-Type", "text/csv")
	attachmentHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name()))
	attachment, err := writer.CreatePart(attachmentHeader)
	if err != nil {
		return err
	}
	attachment.Write(fileContent)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), auth, from, to, email.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func NewService(repository repository.Repository) Service {
	return service{repository: repository}
}
