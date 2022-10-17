package model

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"
)

type OrderModel struct {
	Uid  string
	Json OrderJson
}

type OrderJson struct {
	OrderUid    string `json:"order_uid" validate:"required,min=1"`
	TrackNumber string `json:"track_number" validate:"required"`
	Entry       string `json:"entry" validate:"required"`
	Delivery    struct {
		Name    string `json:"name" validate:"required"`
		Phone   string `json:"phone" validate:"required"`
		Zip     string `json:"zip" validate:"required"`
		City    string `json:"city" validate:"required"`
		Address string `json:"address" validate:"required"`
		Region  string `json:"region" validate:"required"`
		Email   string `json:"email" validate:"required"`
	} `json:"delivery" validate:"required"`
	Payment struct {
		Transaction  string `json:"transaction" validate:"required"`
		RequestId    string `json:"request_id"`
		Currency     string `json:"currency" validate:"required"`
		Provider     string `json:"provider" validate:"required"`
		Amount       int    `json:"amount" validate:"required"`
		PaymentDt    int    `json:"payment_dt" validate:"required"`
		Bank         string `json:"bank" validate:"required"`
		DeliveryCost int    `json:"delivery_cost" validate:"required"`
		GoodsTotal   int    `json:"goods_total" validate:"required"`
		CustomFee    int    `json:"custom_fee" validate:"required"`
	} `json:"payment" validate:"required"`
	Items []struct {
		ChrtId      int    `json:"chrt_id" validate:"required"`
		TrackNumber string `json:"track_number" validate:"required"`
		Price       int    `json:"price" validate:"required"`
		Rid         string `json:"rid" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Sale        int    `json:"sale" validate:"required"`
		Size        string `json:"size" validate:"required"`
		TotalPrice  int    `json:"total_price" validate:"required"`
		NmId        int    `json:"nm_id" validate:"required"`
		Brand       string `json:"brand" validate:"required"`
		Status      int    `json:"status" validate:"required"`
	} `json:"items" validate:"required"`
	Locale            string    `json:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	Shardkey          string    `json:"shardkey" validate:"required"`
	SmId              int       `json:"sm_id" validate:"required"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required"`
}

func NewFromByte(value []byte) (*OrderModel, error) {
	var oj OrderJson
	var om OrderModel
	err := json.Unmarshal(value, &oj)

	fields := reflect.ValueOf(&oj).Elem()
	for i := 0; i < fields.NumField(); i++ {
		tag := fields.Type().Field(i).Tag.Get("validate")
		if strings.Contains(tag, "required") && fields.Field(i).IsZero() {
			return nil, errors.New("required field is missing")
		}

	}

	om.Uid = oj.OrderUid
	om.Json = oj
	return &om, err
}
