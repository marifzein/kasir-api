package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetDailyReport() (*models.DailyReport, error) {
	report := &models.DailyReport{}

	// Query 1: total revenue + jumlah transaksi hari ini (WIB)
	err := r.db.QueryRow(`
		SELECT 
			COALESCE(SUM(total_amount), 0) AS total_revenue,
			COUNT(*) AS total_transaksi
		FROM transactions
		WHERE DATE(created_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta') = CURRENT_DATE
	`).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// Query 2: produk terlaris hari ini
	var nama string
	var qty int
	err = r.db.QueryRow(`
		SELECT p.name AS nama, SUM(td.quantity) AS qty_terjual
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta') = CURRENT_DATE
		GROUP BY p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`).Scan(&nama, &qty)

	if err == sql.ErrNoRows {
		// Tidak ada penjualan hari ini → produk_terlaris null
		report.ProdukTerlaris = nil
	} else if err != nil {
		return nil, err
	} else {
		report.ProdukTerlaris = &models.ProdukTerlaris{
			Nama:       nama,
			QtyTerjual: qty,
		}
	}

	return report, nil
}

// GetReportByRange - untuk optional challenge
func (r *ReportRepository) GetReportByRange(startDate, endDate string) (*models.DailyReport, error) {
    report := &models.DailyReport{}

    // Parse tanggal (opsional, bisa langsung pakai string kalau yakin format benar)
    // Tapi untuk aman, kita pakai langsung di query

    // Query 1: total revenue + jumlah transaksi di range
    err := r.db.QueryRow(`
        SELECT 
            COALESCE(SUM(total_amount), 0),
            COUNT(*)
        FROM transactions
        WHERE DATE(created_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta') 
              BETWEEN $1 AND $2
    `, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaksi)
    if err != nil {
        return nil, err
    }

    // Query 2: produk terlaris di range
    var nama string
    var qty int
    err = r.db.QueryRow(`
        SELECT p.name, SUM(td.quantity)
        FROM transaction_details td
        JOIN products p ON td.product_id = p.id
        JOIN transactions t ON td.transaction_id = t.id
        WHERE DATE(t.created_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta') 
              BETWEEN $1 AND $2
        GROUP BY p.name
        ORDER BY SUM(td.quantity) DESC
        LIMIT 1
    `, startDate, endDate).Scan(&nama, &qty)

    if err == sql.ErrNoRows {
        report.ProdukTerlaris = nil
    } else if err != nil {
        return nil, err
    } else {
        report.ProdukTerlaris = &models.ProdukTerlaris{
            Nama:       nama,
            QtyTerjual: qty,
        }
    }

    return report, nil
}