package models

import "time"

type Order struct {
    ID        int64     `json:"id"`
    ProductID int64     `json:"product_id"`
    BuyerID   string    `json:"buyer_id"`
    Quantity  int       `json:"quantity"`
    CreatedAt time.Time `json:"created_at"`
}
