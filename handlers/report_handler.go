package handlers

import (
	"encoding/json"
	"kasir-api/services"
	"net/http"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) HandleDailyReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	report, err := h.service.GetDailyReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// HandleReportRange - GET /api/report?start_date=...&end_date=...
func (h *ReportHandler) HandleReportRange(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    startDate := r.URL.Query().Get("start_date")
    endDate := r.URL.Query().Get("end_date")

    if startDate == "" || endDate == "" {
        http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
        return
    }

    // Optional: validasi format tanggal (sederhana)
    // Bisa pakai time.Parse kalau mau lebih ketat

    report, err := h.service.GetReportByRange(startDate, endDate)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(report)
}