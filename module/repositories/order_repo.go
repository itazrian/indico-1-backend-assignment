package repositories

import (
    "backend-assignment/module/database"
    "backend-assignment/module/models"
)

type OrderRepo struct{}

func NewOrderRepo() *OrderRepo { return &OrderRepo{} }

func (r *OrderRepo) Create(o *models.Order) error {
    res, err := database.DB.Exec("INSERT INTO orders (product_id, buyer_id, quantity) VALUES (?,?,?)", o.ProductID, o.BuyerID, o.Quantity)
    if err != nil { return err }
    id, _ := res.LastInsertId()
    o.ID = id
    return nil
}

func (r *OrderRepo) GetByID(id int64) (*models.Order, error) {
    row := database.DB.QueryRow("SELECT id, product_id, buyer_id, quantity, created_at FROM orders WHERE id = ?", id)
    o := &models.Order{}
    if err := row.Scan(&o.ID, &o.ProductID, &o.BuyerID, &o.Quantity, &o.CreatedAt); err != nil {
        return nil, err
    }
    return o, nil
}
