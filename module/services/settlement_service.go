package services

import (
    "backend-assignment/module/repositories"
    "database/sql"
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strconv"
    "sync"
    "time"

    "github.com/google/uuid"
)

type SettlementService struct {
    DB *sql.DB
    TxnRepo *repositories.TransactionRepository
    SettleRepo *repositories.SettlementRepo
    JobRepo *repositories.JobRepo
    Workers int
    BatchSize int
    OutputFolder string
}

func NewSettlementService(db *sql.DB) *SettlementService {
    return &SettlementService{
        DB: db,
        TxnRepo: repositories.NewTransactionRepo(db),
        SettleRepo: repositories.NewSettlementRepo(),
        JobRepo: repositories.NewJobRepo(),
        Workers: 4,
        BatchSize: 5000,
        OutputFolder: "/tmp/settlements",
    }
}

func (s *SettlementService) StartJob(from, to time.Time) (string, error) {
    jobID := "job_" + uuid.New().String()
    if err := s.JobRepo.Create(jobID); err != nil {
        return "", err
    }
    _ = os.MkdirAll(s.OutputFolder, 0o755)
    go func() {
        if err := s.run(jobID, from, to); err != nil {
            log.Printf("job %s failed: %v", jobID, err)
            _ = s.JobRepo.UpdateProgress(jobID, 0, 0, 0, "FAILED", nil)
        }
    }()
    return jobID, nil
}

func (s *SettlementService) run(jobID string, from, to time.Time) error {
    if err := s.JobRepo.UpdateProgress(jobID, 0, 0, 0, "RUNNING", nil); err != nil { return err }
    total, err := s.TxnRepo.CountPaidBetween(from, to)
    if err != nil { return err }
    if err := s.JobRepo.UpdateProgress(jobID, 0, total, 0, "RUNNING", nil); err != nil { return err }

    csvPath := filepath.Join(s.OutputFolder, fmt.Sprintf("%s.csv", jobID))
    f, err := os.Create(csvPath)
    if err != nil { return err }
    w := csv.NewWriter(f)
    _ = w.Write([]string{"merchant_id","date","gross","fee","net","txn_count"})
    w.Flush()
    defer func() { w.Flush(); f.Close() }()

    batchesCh := make(chan []repositories.TxnRow, s.Workers*2)
    resultsCh := make(chan map[string]aggRow, s.Workers*2)
    var wg sync.WaitGroup

    for i:=0; i<s.Workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for rows := range batchesCh {
                if cancel, _ := s.JobRepo.IsCancelRequested(jobID); cancel { return }
                local := make(map[string]aggRow)
                for _, r := range rows {
                    key := r.Merchant + "|" + r.PaidDate.Format("2006-01-02")
                    a := local[key]
                    a.Merchant = r.Merchant
                    a.Date = r.PaidDate
                    a.Gross += r.Amount
                    a.Fee += r.Fee
                    a.Net += (r.Amount - r.Fee)
                    a.Count++
                    local[key] = a
                }
                select {
                case resultsCh <- local:
                default:
                }
            }
        }()
    }

    var lastID int64 = 0
    var produced int64 = 0
produce:
    for {
        if cancel, _ := s.JobRepo.IsCancelRequested(jobID); cancel {
            _ = s.JobRepo.UpdateProgress(jobID, produced, total, int(percent(produced,total)), "CANCELED", nil)
            break
        }
        rows, maxID, err := s.TxnRepo.ReadPaidBatchAfterID(from, to, lastID, s.BatchSize)
        if err != nil {
            close(batchesCh); wg.Wait(); close(resultsCh); return err
        }
        if len(rows) == 0 { break }
        lastID = maxID
        select {
        case batchesCh <- rows:
            produced += int64(len(rows))
            _ = s.JobRepo.UpdateProgress(jobID, produced, total, int(percent(produced,total)), "RUNNING", nil)
        default:
            time.Sleep(10 * time.Millisecond)
            continue produce
        }
    }
    close(batchesCh)
    wg.Wait()
    close(resultsCh)

    finalAgg := make(map[string]aggRow)
    var processed int64 = 0
    for m := range resultsCh {
        for k, v := range m {
            e := finalAgg[k]
            e.Merchant = v.Merchant
            e.Date = v.Date
            e.Gross += v.Gross
            e.Fee += v.Fee
            e.Net += v.Net
            e.Count += v.Count
            finalAgg[k] = e
        }
    }

    tx, err := s.DB.Begin()
    if err != nil { return err }
    runID := uuid.New().String()
    for _, v := range finalAgg {
        _ = w.Write([]string{v.Merchant, v.Date.Format("2006-01-02"), strconv.FormatInt(v.Gross,10), strconv.FormatInt(v.Fee,10), strconv.FormatInt(v.Net,10), strconv.FormatInt(v.Count,10)})
        _, err := tx.Exec(`
            INSERT INTO settlements (merchant_id, date, gross_cents, fee_cents, net_cents, txn_count, generated_at, unique_run_id)
            VALUES (?, ?, ?, ?, ?, ?, NOW(), ?)
            ON DUPLICATE KEY UPDATE
              gross_cents = gross_cents + VALUES(gross_cents),
              fee_cents   = fee_cents   + VALUES(fee_cents),
              net_cents   = net_cents   + VALUES(net_cents),
              txn_count   = txn_count   + VALUES(txn_count),
              generated_at = NOW(),
              unique_run_id = VALUES(unique_run_id)
        `, v.Merchant, v.Date.Format("2006-01-02"), v.Gross, v.Fee, v.Net, v.Count, runID)
        if err != nil { tx.Rollback(); return err }
    }
    if err := tx.Commit(); err != nil { return err }
    w.Flush()
    _ = s.JobRepo.UpdateProgress(jobID, processed, total, 100, "DONE", &csvPath)
    return nil
}

type aggRow struct {
    Merchant string
    Date time.Time
    Gross int64
    Fee int64
    Net int64
    Count int64
}

func percent(a,b int64) float64 {
    if b==0 { return 100 }
    return (float64(a)/float64(b))*100
}
