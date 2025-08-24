package repositories

import (
    "backend-assignment/module/database"
    "errors"
)

type ProductRepo struct{}

func NewProductRepo() *ProductRepo { return &ProductRepo{} }

func (r *ProductRepo) ReduceStock(productID int64, qty int) error {
    tx, err := database.DB.Begin()
    if err != nil { return err }
    defer tx.Rollback()

    var stock int
    err = tx.QueryRow("SELECT stock FROM products WHERE id = ? FOR UPDATE", productID).Scan(&stock)
    if err != nil { return err }
    if stock < qty { return errors.New("OUT_OF_STOCK") }
    _, err = tx.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", qty, productID)
    if err != nil { return err }
    return tx.Commit()
}
