package main

import (
	"log"
	"net"
	pb "product/generated/product"
	"product/service"
	"product/storage/postgres"

	"google.golang.org/grpc"
	_ "github.com/lib/pq"

)

func main() {
	db, err := postgres.ConnectDb()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	listener, err := net.Listen("tcp", ":50050")

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	product := service.ProductService{
		OrderRepo:   postgres.NewOrderRepo(db),
		PaymentRepo: postgres.NewPaymentRepo(db),
		ProductRepo: postgres.NewProductRepo(db),
		RetingRepo:  postgres.NewRetingRepo(db),
	}

	pb.RegisterProductServiceServer(s, &product)

	log.Println("server is running on :50050 ...")
	if err = s.Serve(listener); err != nil {
		log.Fatal(err)
	}

}
