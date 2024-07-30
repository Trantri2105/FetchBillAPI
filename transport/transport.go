package transport

import (
	"backend/dto"
	endpoints "backend/endpoint"
	"backend/initilizer"
	"backend/repository"
	"backend/service"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
)

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request dto.Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func NewRouter() *gin.Engine {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	db := initilizer.DbConnect()
	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	billEndpoint := endpoints.NewBillEndpoint(service)

	getCamBillHandler := httptransport.NewServer(
		billEndpoint.GetCamBillEndpoint(),
		decodeRequest,
		encodeResponse,
		options...,
	)

	r := gin.Default()
	r.POST("/", gin.WrapH(getCamBillHandler))

	return r
}
