package repositories

import (
    "backend-assignment/module/database"
    //"database/sql"
)

type JobRepo struct{}

func NewJobRepo() *JobRepo { return &JobRepo{} }

func (r *JobRepo) Create(jobID string) error {
    _, err := database.DB.Exec("INSERT INTO jobs (id, status, progress, processed, total) VALUES (?, 'QUEUED', 0, 0, 0)", jobID)
    return err
}

func (r *JobRepo) UpdateProgress(jobID string, processed, total int64, progress int, status string, resultPath *string) error {
    if resultPath != nil {
        _, err := database.DB.Exec("UPDATE jobs SET processed=?, total=?, progress=?, status=?, result_path=? WHERE id=?", processed, total, progress, status, *resultPath, jobID)
        return err
    }
    _, err := database.DB.Exec("UPDATE jobs SET processed=?, total=?, progress=?, status=? WHERE id=?", processed, total, progress, status, jobID)
    return err
}

func (r *JobRepo) GetStatusSummary(jobID string) (status string, processed, total int64, progress int, err error) {
    err = database.DB.QueryRow("SELECT status, processed, total, progress FROM jobs WHERE id=?", jobID).Scan(&status, &processed, &total, &progress)
    return
}

func (r *JobRepo) RequestCancel(jobID string) error {
    _, err := database.DB.Exec("UPDATE jobs SET status='CANCEL_REQUESTED' WHERE id=?", jobID)
    return err
}

func (r *JobRepo) IsCancelRequested(jobID string) (bool, error) {
    var status string
    if err := database.DB.QueryRow("SELECT status FROM jobs WHERE id=?", jobID).Scan(&status); err != nil {
        return false, err
    }
    return status == "CANCEL_REQUESTED", nil
}
