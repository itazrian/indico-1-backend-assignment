package repositories

import (
    "backend-assignment/module/database"
    "time"
)

type SettlementRepo struct{}

func NewSettlementRepo() *SettlementRepo { return &SettlementRepo{} }

func (r *SettlementRepo) Upsert(merchant string, date time.Time, gross, fee, net, cnt int64, runID string) error {
    _, err := database.DB.Exec(`
        INSERT INTO settlements (merchant_id, date, gross_cents, fee_cents, net_cents, txn_count, generated_at, unique_run_id)
        VALUES (?, ?, ?, ?, ?, ?, NOW(), ?)
        ON DUPLICATE KEY UPDATE
          gross_cents = gross_cents + VALUES(gross_cents),
          fee_cents   = fee_cents   + VALUES(fee_cents),
          net_cents   = net_cents   + VALUES(net_cents),
          txn_count   = txn_count   + VALUES(txn_count),
          generated_at = NOW(),
          unique_run_id = VALUES(unique_run_id)
    `, merchant, date.Format("2006-01-02"), gross, fee, net, cnt, runID)
    return err
}
