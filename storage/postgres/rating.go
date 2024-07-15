package postgres

import (
	"database/sql"
	pb "product/generated/product"
	_ "github.com/lib/pq"

)

type RetingRepo struct {
	DB *sql.DB
}

func NewRetingRepo(db *sql.DB) *RetingRepo {
	return &RetingRepo{DB: db}
}

func (p *RetingRepo) CreateRatingProducts(in *pb.CreateRatingProductsRequest) (*pb.CreateRatingProductsResponse, error) {

	var resp pb.CreateRatingProductsResponse

	err := p.DB.QueryRow(`
			INSERT INTO
			ratings(
				product_id,
				user_id,
				reting,
				commit
			)
			VALUES(
				$1,
				$2,
				$3,
				$4) 
			RETURNING
				id,
				product_id,
				user_id,
				reting,
				commient,
				creted_at`, in.ProductId, in.UserId, in.Rating, in.Comment).Scan(
		&resp.Id, &resp.ProductId, &resp.UserId, &resp.Rating, &resp.Comment, &resp.CretedAt)

	return &resp, err
}

func (r *RetingRepo) GetProductRatings(in *pb.GetProductRatingsRequest) (*pb.GetProductRatingsResponse, error) {

	average, err := r.DB.Query(`
			SELECT sum(rating)/count(rating), count(rating)
			FROM
			ratings
			WHERE product_id=$1`, in.ProductId)

	if err != nil {
		return nil, err
	}
	count := 0
	var ave float64
	for average.Next() {

		err = average.Scan(&ave, &count)
		if err != nil {
			return nil, err
		}
	}
	resp, err := r.DB.Query(`
			SELECT 
				id,
				product_id,
				user_id,
				rating,
				comment,
				creted_at
			FROM
			ratings
			WHERE product_id=$1`, in.ProductId)
	if err != nil {
		return nil, err
	}

	var ratings []*pb.CreateRatingProductsResponse

	for resp.Next() {
		var rating pb.CreateRatingProductsResponse
		err := resp.Scan(&rating.Id, &rating.ProductId, &rating.UserId, &rating.Rating, &rating.Comment, &rating.CretedAt)
		if err != nil {
			return nil, err
		}

		ratings = append(ratings, &rating)
	}
	return &pb.GetProductRatingsResponse{
		Retings:       ratings,
		AverageRating: float32(ave),
		TotalRating:   float32(count),
	}, nil
}
