package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackyansen22/crud-category/internal/model"
)

type ReportRepository interface {
	GetReport(ctx context.Context, start, end time.Time) (*model.ReportResponse, error)
}

type reportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetReport(
	ctx context.Context,
	start, end time.Time,
) (*model.ReportResponse, error) {

	var report model.ReportResponse

	// ===============================
	// Total revenue & total transaksi
	// ===============================
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COALESCE(SUM(total_amount), 0),
			COUNT(*)
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`, start, end).Scan(
		&report.TotalRevenue,
		&report.TotalTransaksi,
	)
	if err != nil {
		return nil, err
	}

	// ===============================
	// Produk terlaris (qty terbanyak)
	// ===============================
	err = r.db.QueryRowContext(ctx, `
		SELECT
			p.nama,
			SUM(td.quantity) AS qty_terjual
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.nama
		ORDER BY qty_terjual DESC
		LIMIT 1
	`, start, end).Scan(
		&report.ProdukTerlaris.Nama,
		&report.ProdukTerlaris.QtyTerjual,
	)

	// kalau tidak ada transaksi sama sekali
	if err == sql.ErrNoRows {
		report.ProdukTerlaris = model.BestSeller{}
		return &report, nil
	}
	if err != nil {
		return nil, err
	}

	return &report, nil
}
