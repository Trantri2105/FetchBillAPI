package dto

type Request struct {
	Start string `json:"start" validate:"required"`
	End   string `jdon:"end" validate:"required"`
}
