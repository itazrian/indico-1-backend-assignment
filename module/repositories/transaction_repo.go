package repositories

import (
    "database/sql"
    "time"
)

type TxnRow struct {
    ID int64
    Merchant string
    Amount int64
    Fee int64
    PaidDate time.Time
}

type TransactionRepository struct{
    DB *sql.DB
}

func NewTransactionRepo(db *sql.DB) *TransactionRepository { return &TransactionRepository{DB: db} }

func (r *TransactionRepository) CountPaidBetween(from, to time.Time) (int64, error) {
    var n int64
    err := r.DB.QueryRow("SELECT COUNT(1) FROM transactions WHERE status='PAID' AND paid_at >= ? AND paid_at < ?", from, to).Scan(&n)
    return n, err
}

func (r *TransactionRepository) ReadPaidBatchAfterID(from, to time.Time, lastID int64, limit int) ([]TxnRow, int64, error) {
    rows, err := r.DB.Query(`SELECT id, merchant_id, amount_cents, fee_cents, paid_at FROM transactions WHERE status='PAID' AND paid_at >= ? AND paid_at < ? AND id > ? ORDER BY id LIMIT ?`, from, to, lastID, limit)
    if err != nil { return nil, lastID, err }
    defer rows.Close()
    out := make([]TxnRow,0)
    currMax := lastID
    for rows.Next() {
        var r0 TxnRow
        if err := rows.Scan(&r0.ID, &r0.Merchant, &r0.Amount, &r0.Fee, &r0.PaidDate); err != nil {
            return nil, lastID, err
        }
        if r0.ID > currMax { currMax = r0.ID }
        out = append(out, r0)
    }
    return out, currMax, rows.Err()
}
