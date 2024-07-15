package postgres

import (
	"database/sql"
	pb "product/generated/product"
	_ "github.com/lib/pq"

)

type OrderRepo struct {
	DB *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{DB: db}
}

func (o OrderRepo) CreateOrder(in *pb.CreateOrderResponse) (*pb.CreateOrderResponse, error) {

	err := o.DB.QueryRow(`
			INSERT INTO
			orders(
				user_id,
				total_amount,
				status,
				shipping_address
			)
			VALUES(
				$1,
				$2,
				$3,
				$4)
			RETURNING
				id,
				user_id,
				total_amount,
				status,
				shipping_address,
				created_at
			`, in.UserId, in.TotalAmount, in.Status, in.ShippingAddress).Scan(&in.Id, &in.UserId, &in.TotalAmount, &in.Status, &in.ShippingAddress, &in.CretedAt)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (o *OrderRepo) DeleteOrder(in *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {

	var resp pb.DeleteOrderResponse

	err := o.DB.QueryRow(`
			UPDATE orders
			SET
				deleted_at=CURRENT_TIMESTAMP
			WHERE
				id=$1
			AND
				deleted_at=0
			RETURNING
				id,
				status,
				deleted_at
			`, in.OrderId).Scan(&resp.Id, &resp.Status, &resp.DeletedAt)

	return &resp, err
}

func (o *OrderRepo) UpdateOrder(in *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {

	var resp pb.UpdateOrderResponse

	err := o.DB.QueryRow(`
			UPDATE orders
			SET
				status=&1,
				update_at=CURRENT_TIMESTAMP
			WHERE
				id=$1
			AND
				deleted_at=0
			RETURNING
				id,
				status,
				update_at
			`, in.Id).Scan(&resp.Id, &resp.Status, &resp.UpdatedAt)

	return &resp, err
}

func (o *OrderRepo) GetallOrder(in *pb.GetallOrderRequest) (*pb.GetallOrderResponse, error) {

	rows, err := o.DB.Query(`SELECT idFROM orders`)
	if err != nil {
		return nil, err
	}

	count := 0

	for rows.Next() {
		count++
	}

	resp, err := o.DB.Query(`
			SELECT 
				id,
				user_id,
				total_amount,
				status,
				shipping_address,
				created_at,
]			FROM orders
			WHERE limit=$1 offset=$2
				`, in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	var orders []*pb.Order

	for resp.Next() {
		var order pb.Order
		err = resp.Scan(&order.Id, &order.UserId, &order.TotalAmount, &order.Status, &order.ShippingAddress, &order.CreatedAt)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}
	return &pb.GetallOrderResponse{
		Orders: orders,
		Total:  int32(count),
		Limit:  in.Limit,
		Offset: in.Offset}, nil
}

func (o *OrderRepo) GetByIdOrder(in *pb.GetByIdOrderRerquest) (*pb.GetByIdOrderResponse, error) {

	var resp pb.GetByIdOrderResponse

	err := o.DB.QueryRow(`
			SELECT 
				id,
				user_id,
				total_amount,
				status,
				shipping_address,
				creted_at,
				uptaded_at
			FROM orders
			WHERE id=$1 and deleted_at=0`, in.OrderId).Scan(
		&resp.Id, &resp.UserId, &resp.TotalAmount, &resp.Status, &resp.ShippingAddress, &resp.CretedAt, &resp.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (o *OrderRepo) PaymentQuery(id string) (*pb.CreatePaymentResponse, error) {
	var resp pb.CreatePaymentResponse

	rows, err := o.DB.Exec(`
			UPDATE orders
			SET status='success,
				updated_at=CURRENT_TIMESTAMP
			WHERE id=$1 and deleted_at=0`, id)
	if err != nil {
		return nil, err
	}
	rowsaff, err := rows.RowsAffected()
	if err != nil || rowsaff == 0 {
		return nil, err
	}

	err = o.DB.QueryRow(`
				SLELCT 
					total_amount,
					status
				FROM orders
				WHERE id=$1 and deleted_at=0
					`, id).Scan(&resp.Amount, &resp.Status)
	return &resp, err
}

func (o *OrderRepo) PoductId(orederId string) (string, error) {
	var id string
	err := o.DB.QueryRow(`
			SELECT 
				product_id
			FROM order_item
			WHERE order_id=$1`, orederId).Scan(&id)

	return id, err
}
