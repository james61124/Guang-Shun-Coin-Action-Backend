package response

import (
	"time"
	"database/sql"
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

type AddProductResponse struct {
	ProductId  string `json:"productId"`
}

type ProductResponse struct {
	ProductId   string    `json:"productId"`
	ProductName string    `json:"productName"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	MinBidPrice float64   `json:"minBidPrice"`
	ImgUrl      string    `json:"imgUrl"`
	StartAt    time.Time `json:"createAt"`
	EndedAt     time.Time `json:"endedAt"`
	BidCount    int      `json:"bidCount"`
	IsStar      bool     `json:"isStar"`
}

type ProductListResponse struct {
    Products   []ProductResponse `json:"products"`
    TotalPages int                        `json:"totalPages"`
}

type UpdateProductPageResponse struct {
	ProductId   string    `json:"productId"`
	ProductName string    `json:"productName"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	MinBidPrice float64   `json:"minBidPrice"`
	ImgUrl      string    `json:"imgUrl"`
	StartAt    time.Time `json:"createAt"`
	EndedAt     time.Time `json:"endedAt"`
	BidCount    int      `json:"bidCount"`
}

type UpdateProductPageListResponse struct {
    Products   []UpdateProductPageResponse `json:"products"`
    TotalPages int                        `json:"totalPages"`
}

type TotalPagesOfProductResponse struct {
	Total   int    `json:"total"`
}

type BidHistory struct {
	Username string `json:"username"`
	BidPrice float64 `json:"bidPrice"`
	BidTime  time.Time `json:"bidTime"`
	Status sql.NullString `json:"status"`
}

type GetBidHistory struct {
	ProductID string `json:"productID"`
	ProductName string `json:"productName"`
	BidPrice float64 `json:"bidPrice"`
	BidTime  time.Time `json:"bidTime"`
	Status sql.NullString `json:"status"`
	ImageUrl  sql.NullString  `json:"imageUrl"`
}

type DetailResponse struct {
	Name        string      `json:"name"`
	Category        string      `json:"category"`
	CurrentPrice   float64     `json:"currentPrice"`
	IsStar        bool        `json:"isStar"`
	Price       float64     `json:"price"`
	MinBidPrice float64     `json:"minBidPrice"`
	StartTime   time.Time   `json:"startTime"`
	EndTime     time.Time   `json:"endTime"`
	Description string      `json:"description"`
	ImageUrl    []string    `json:"imageUrl"`
	TotalPageOfHistory    int   `json:"totalPageOfHistory"`
	History     []BidHistory `json:"history"`
}

type GetUserInfoResponse struct {
	RealName        string      `json:"realName"`
	NickName        string      `json:"nickName"`
	Cellphone        string      `json:"cellphone"`
	FbAccount        string      `json:"fbAccount"`
	Email        string      `json:"email"`
	Postcode        string      `json:"postcode"`
	ShippingAddr        string      `json:"shippingAddr"`
	Username        string      `json:"username"`
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
