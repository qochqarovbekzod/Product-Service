package postgres

import (
	"database/sql"
	"fmt"
	pb "product/generated/product"
	"product/pkg"
	_ "github.com/lib/pq"

)

type ProductRepo struct {
	DB *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{DB: db}
}

func (p *ProductRepo) CreateProduct(in *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	var resp pb.CreateProductResponse
	err := p.DB.QueryRow(`
			INSERT INTO
				products(
					name,
					description,
					price,
					catigory_Id,
					quantity
					)
				values(
					$1,
					$2,
					$3,
					$4,
					$5
				returning 
					name,
					description,
					price,
					catigory_Id,
					quantity
					creted_at
					)`, in.Name, in.Description, in.Price, in.CategoryId, in.Quantity).Scan(
		resp.Name, resp.Description, resp.Price, resp.CategoryId, resp.Quantity, resp.CreatedAd)

	return &resp, err
}

func (p *ProductRepo) UpdateProduct(in *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	var resp pb.UpdateProductResponse

	err := p.DB.QueryRow(`
			UPDATE products
			SET 
				name=$1,
				price=$2,
				updated_at=CURRENT_TIMESTAMP
			WHERE 
				di=$3 
			AND
				deleted_at=0
			RETURNING
				name,
				description,
				price,
				catigory_Id,
				quantity
				updated_at
				`, in.Name, in.Price, in.Id).Scan(
		resp.Name, resp.Description, resp.Price, resp.CategoryId, resp.Quantity, resp.UpdatedAd)

	return &resp, err
}

func (p ProductRepo) DeleteProduct(in *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	resp, err := p.DB.Exec(`
			UPDATE products
			SET 
				deleted_at=CURRENT_TIMESTAMP
			WHERE 
				di=$1
			AND
				deleted_at=0`, in.Id)
	if err != nil {
		return nil, err
	}
	rowse, err := resp.RowsAffected()
	if err != nil || rowse == 0 {
		return &pb.DeleteProductResponse{Success: "bu foydalanuvchi allaqachon ochirilgan"}, err
	}
	return &pb.DeleteProductResponse{Success: "ochirildi"}, nil
}

func (p ProductRepo) GetProduct(in *pb.GetProductRequest) (*pb.GetProductResponse, error) {

	total, err := p.DB.Query(`
			SELECT 
			*
			FROM products`)

	if err != nil {
		return nil, err
	}

	count := 0
	for total.Next() {

		count++
	}

	resp, err := p.DB.Query(`
			SELECT 
				id,
				name,
				price,
				catigory_id
			FROM
				products 
			WHERE
				limit=$1
				offset=$2 
			`, in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	var products []*pb.Product

	for resp.Next() {
		var product pb.Product
		err := resp.Scan(&product.Id, &product.Name, &product.Price, &product.CategoryId)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return &pb.GetProductResponse{
		Products: products,
		Total:    int32(count),
		Limit:    in.Limit,
		Offset:   in.Offset,
	}, nil

}

func (p *ProductRepo) GetbyIdProduct(in *pb.GetbyIdProductRequest) (*pb.GetbyIdProductResponse, error) {

	var resp pb.GetbyIdProductResponse

	err := p.DB.QueryRow(`
			SELECT
				id,
				name,
				description,
				price,
				catigoty_id,
				quantity,
				creted_at,
				updated_at
			FROM
				products
			WHERE
				id=$1
			AND
				deleted_at=0`, in.Id).Scan(
		&resp.Id, &resp.Name, &resp.Description, &resp.Price, &resp.CategoryId, resp.Quantity, &resp.CretedAt, &resp.UpdatedAd)

	return &resp, err
}

func (p ProductRepo) GetallProducts(in *pb.GetallProductsRequest) (*pb.GetallProductsResponse, error) {
	var (
		params   = make(map[string]interface{})
		arrCount []interface{}
		arr      []interface{}
		filter   string
		total    string
	)

	query := `
		SELECT 
			p.id,
			p.name,
			p.price,
			p.cotegory_id
		FROM 
			products as p 
		JOIN 
			product_categories AS pc 
		ON 
			p.cotigory_id=pc.id 
		WHERE true `
	total = query
	if len(in.Category) > 0 {
		params["cotegory"] = in.Category
		filter += " and pc.cotegory = :cotegory "
		total += " and pc.cotegory = :cotegory "
	}

	if in.MaxPrice > 0 {
		params["max_price"] = in.MaxPrice
		filter += " and p.price < :max_price "
		total += " and p.price < :max_price "
	}

	if in.MinPrice > 0 {
		params["min_price"] = in.MinPrice
		filter += " and price > :min_price "
		total += " and p.price < :max_price "
	}

	total, arrCount = pkg.ReplaceQueryParams(total, params)
	t, err := p.DB.Query(total, arrCount...)
	if err != nil {
		return nil, err
	}
	count := 0
	for t.Next() {
		count++
	}

	if in.Offset > 0 {
		params["offset"] = in.Offset
		filter += " OFFSET :offset"
	}

	if in.Limit > 0 {
		params["limit"] = in.Limit
		filter += " LIMIT :limit"
	}
	query = query + filter

	query, arr = pkg.ReplaceQueryParams(query, params)
	fmt.Println(query, arr)
	rows, err := p.DB.Query(query, arr...)
	fmt.Println(err, query)
	if err != nil {
		return nil, err
	}

	var products []*pb.Product
	for rows.Next() {
		var product pb.Product
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.CategoryId)

		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &pb.GetallProductsResponse{
		Products: products,
		Total:    int32(count),
		Limit:    in.Limit,
		Offset:   in.Offset}, nil
}

func (p *ProductRepo) OrderQuery(id string) (*float64, error) {
	var price float64

	err := p.DB.QueryRow("SELECT price FROM products WHERE id=$1", id).Scan(&price)

	return &price, err
}

func (p *ProductRepo) OrederGetbyIdQuery(user_id string) ([]*pb.ItemResponse, error) {
	var items []*pb.ItemResponse

	resp, err := p.DB.Query(`
			SELECT 
				id,
				quantity,
				price
			FROM
			WHERE
				user_id=$1 and deleted_at=0`, user_id)

	for resp.Next() {
		var item pb.ItemResponse
		err := resp.Scan(&item.ProductId, &item.Quantity, &item.Price)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	return items, err
}

func (p *ProductRepo) ProductQuery(id string) (*pb.TopProduct, error) {
	var resp pb.TopProduct
	err := p.DB.QueryRow(`
			SELECT 
				id,
				name,
				price
			FROM
			products
			WHERE
				user_id=&1 and deleted_at=0
				`, id).Scan(
		resp.Id, resp.Name, resp.Revenume)
	return &resp, err
}

func (p *ProductRepo) TopProduct(id string) (*pb.TopProduct, error) {
	var resp pb.TopProduct
	err := p.DB.QueryRow(`
			SELECT 
				id,
				name,
				price
			from products
			WHERE id=$1`, id).Scan(&resp.Id, &resp.Name, &resp.Revenume)
	return &resp, err
}

func (p *ProductRepo) GetProductRecommendations(in *pb.GetProductRecommendationsRequest) (*pb.GetProductRecommendationsResponse, error) {
	rows, err := p.DB.Query(`
			SELECT
				id,
				name,
				price,
				category_id
			FROM products
			WHERE 
				user_id=$1 limit $2
				`, in.UserId, in.Limit)
	if err != nil {
		return nil, err
	}

	var products []*pb.Recommendation
	var product pb.Recommendation

	for rows.Next() {
		err = rows.Scan(&product.Id, &product.Name, &product.Price, &product.CategoryId)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	return &pb.GetProductRecommendationsResponse{
		Recommendations: products,
	}, nil

}
