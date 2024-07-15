package postgres

import (
	"database/sql"
	pb "product/generated/product"
	_ "github.com/lib/pq"

)

type PaymentRepo struct {
	DB *sql.DB
}

func NewPaymentRepo(db *sql.DB) *PaymentRepo {
	return &PaymentRepo{DB: db}
}

func (p *PaymentRepo) CreatePayment(in *pb.CreatePaymentResponse) (*pb.CreatePaymentResponse, error) {
	err := p.DB.QueryRow(`
			INSERT INTO
			payments(
				order_id,
				amount,
				status,
				payment_method
			VALUES(
				$1,
				$2,
				$3,
				$4
			) 
			RETURNING
				order_id,
				amount,
				status,
				transaction_id,
				payment_method,
				created_at
				)`, in.OrderId, in.Amount, in.Status, in.PaymentMethod).Scan(
		&in.OrderId, &in.Amount, &in.Status, &in.TransactionId, &in.PaymentMethod, &in.CreatedAt)

	return in, err
}

func (p *PaymentRepo) PaymentStatus(in *pb.PaymentStatusRequest) (*pb.PaymentStatusResponse, error) {
	var resp pb.PaymentStatusResponse
	err := p.DB.QueryRow(`
			SELECT 
				id,
				order_id,
				amount,
				status.
				transaction_id,
				created_at
			FROM payments
			WHERE
				order_id=$1
			`, in.OrderId).Scan(
		&resp.PaymentId, &resp.OrderId, &resp.Amount, &resp.Status, &resp.TransactionId, &resp.CreatedAt)

	return &resp, err
}

func (p *PaymentRepo) GetStatisticsProduct(in *pb.GetStatisticsRequest) ([]*pb.TopProduct, error) {
	rows, err := p.DB.Query(`
			SELECT 
				user_id,
				amount
			FROM
				payments
			WHERE
				created_at>$1 AND created_at<$2
				`, in.StartDate, in.EndDate)
		if err != nil {
			return nil, err
		}

	var responses []*pb.TopProduct

	for rows.Next() {
		var resp pb.TopProduct
		err = rows.Scan(&resp.Id, &resp.Revenume)
		if err != nil {
			return nil, err
		}

		responses = append(responses, &resp)

	}

	return responses, nil
}
