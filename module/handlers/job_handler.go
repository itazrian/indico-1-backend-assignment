package handlers

import (
    "backend-assignment/module/services"
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
)

func JobRoutes(r *gin.Engine, svc *services.SettlementService) {
    h := NewJobHandler(svc)
    r.POST("/jobs/settlement", h.Create)
    r.GET("/jobs/:id", h.Status)
    r.POST("/jobs/:id/cancel", h.Cancel)
}

type JobHandler struct {
    Svc *services.SettlementService
}

func NewJobHandler(svc *services.SettlementService) *JobHandler { return &JobHandler{Svc: svc} }

type startReq struct {
    From string `json:"from" binding:"required"`
    To   string `json:"to" binding:"required"`
}

func (h *JobHandler) Create(c *gin.Context) {
    var r startReq
    if err := c.ShouldBindJSON(&r); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
        return
    }
    from, err := time.Parse("2006-01-02", r.From)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid from"}); return }
    to, err := time.Parse("2006-01-02", r.To)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid to"}); return }
    jobID, err := h.Svc.StartJob(from, to.AddDate(0,0,1))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusAccepted, gin.H{"job_id": jobID, "status": "QUEUED"})
}

func (h *JobHandler) Status(c *gin.Context) {
    id := c.Param("id")
    status, processed, total, progress, err := h.Svc.JobRepo.GetStatusSummary(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error":"NOT_FOUND"})
        return
    }
    resp := gin.H{"job_id": id, "status": status, "processed": processed, "total": total, "progress": progress}
    if status == "DONE" {
        resp["download_url"] = "/downloads/" + id + ".csv"
    }
    c.JSON(http.StatusOK, resp)
}

func (h *JobHandler) Cancel(c *gin.Context) {
    id := c.Param("id")
    if err := h.Svc.JobRepo.RequestCancel(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"failed"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"ok": true})
}
