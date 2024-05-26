package response

import (
	"time"
)

type Response struct {
	Status  bool        `json:"Status"`
	Data    interface{} `json:"Data"`
	Message interface{} `json:"Message"`
}

type RegisterResponse struct {
	UUID  string `json:"UUID"`
	Email string `json:"Email"`
}

type LoginResponse struct {
	UUID  string `json:"UUID"`
	Token string `json:"Token"`
}

type ProductResponse struct {
	ProductId   string    `json:"productId"`
	ProductName string    `json:"productName"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	MinBidPrice float64   `json:"minBidPrice"`
	ImgUrl      string    `json:"imgUrl"`
	CreateAt    time.Time `json:"createAt"`
	EndedAt     time.Time `json:"endedAt"`
}

type ProductRespnseList struct {
	ProductInfo []ProductRespnseList `json: productInfo`
}


type ULIDResponse struct {
	ULID string `json:"ULID"`
}

type UTIDResponse struct {
	UTID string `json:"UTID"`
}

func New() *Response {
	return &Response{
		Status:  false,
		Data:    nil,
		Message: nil,
	}
}
