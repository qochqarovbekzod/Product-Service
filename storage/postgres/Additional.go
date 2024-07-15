package postgres

import (
	"database/sql"
	pb "product/generated/product"
	_ "github.com/lib/pq"

)

type AddRepo struct {
	DB *sql.DB
}

func NewAddRepo(db *sql.DB) AddRepo {
	return AddRepo{
		DB: db}
}

func (a *AddRepo) CreateCategory(in *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	var resp pb.CreateCategoryResponse
	err := a.DB.QueryRow(`
			INSERT INTO
			category(
				name,
				discription)
			VALUES(
				$1,
				$2
			)
			RETURNING
				id,
				name,
				discription,
				created_at
			`, in.Name, in.Description).Scan(
		&resp.Id, &resp.Name, &resp.Description, &resp.CreatedAt)

	return &resp, err
}
