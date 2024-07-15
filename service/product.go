package service

import (
	"context"
	"log/slog"
	pb "product/generated/product"
	"product/storage/postgres"
)

type ProductService struct {
	pb.UnimplementedProductServiceServer
	OrderRepo   *postgres.OrderRepo
	PaymentRepo *postgres.PaymentRepo
	ProductRepo *postgres.ProductRepo
	RetingRepo  *postgres.RetingRepo
	AddRepo     *postgres.AddRepo
	Log         *slog.Logger
}

func NewProductService(order *postgres.OrderRepo, pay *postgres.PaymentRepo, pro *postgres.ProductRepo, rat *postgres.RetingRepo, add *postgres.AddRepo, log *slog.Logger) ProductService {
	return ProductService{
		OrderRepo:   order,
		PaymentRepo: pay,
		ProductRepo: pro,
		RetingRepo:  rat,
		AddRepo:     add,
		Log:         log,
	}
}

func (p *ProductService) CreateProduct(ctx context.Context, in *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	return p.ProductRepo.CreateProduct(in)
}

func (p *ProductService) UpdateProduct(ctx context.Context, in *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	return p.ProductRepo.UpdateProduct(in)
}

func (p *ProductService) DeleteProduct(ctx context.Context, in *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	return p.ProductRepo.DeleteProduct(in)
}

func (p *ProductService) GetProduct(ctx context.Context, in *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	return p.ProductRepo.GetProduct(in)
}

func (p *ProductService) GetbyIdProduct(ctx context.Context, in *pb.GetbyIdProductRequest) (*pb.GetbyIdProductResponse, error) {
	return p.ProductRepo.GetbyIdProduct(in)
}

func (p *ProductService) GetallProducts(ctx context.Context, in *pb.GetallProductsRequest) (*pb.GetallProductsResponse, error) {
	return p.ProductRepo.GetallProducts(in)
}

func (p *ProductService) CreateRatingProducts(ctx context.Context, in *pb.CreateRatingProductsRequest) (*pb.CreateRatingProductsResponse, error) {
	return p.RetingRepo.CreateRatingProducts(in)
}

func (p *ProductService) GetProductRatings(ctx context.Context, in *pb.GetProductRatingsRequest) (*pb.GetProductRatingsResponse, error) {
	return p.RetingRepo.GetProductRatings(in)
}

func (p *ProductService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	var req pb.CreateOrderResponse
	req.UserId = in.UserId
	var items []*pb.ItemResponse
	sum := 0
	for _, val := range in.Items {
		var item pb.ItemResponse
		price, err := p.ProductRepo.OrderQuery(val.ProductId)
		if err != nil {
			return nil, err
		}

		item.Price = float32(val.Quantity) * float32(*price)
		item.Quantity = val.Quantity
		item.ProductId = val.ProductId
		items = append(items, &item)
		sum += int(item.Price)
	}
	req.TotalAmount = float32(sum)
	req.Status = "yaxshi"
	req.ShippingAddress = in.ShippingAddress
	response, err := p.OrderRepo.CreateOrder(&req)
	if err != nil {
		return nil, err
	}

	response.Items = items
	return response, err

}

func (p *ProductService) DeleteOrder(ctx context.Context, in *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	return p.OrderRepo.DeleteOrder(in)
}

func (p *ProductService) UpdateOrder(ctx context.Context, in *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	return p.OrderRepo.UpdateOrder(in)
}

func (p *ProductService) GetallOrder(ctx context.Context, in *pb.GetallOrderRequest) (*pb.GetallOrderResponse, error) {
	return p.OrderRepo.GetallOrder(in)
}

func (p *ProductService) GetByIdOrder(ctx context.Context, in *pb.GetByIdOrderRerquest) (*pb.GetByIdOrderResponse, error) {

	resp, err := p.OrderRepo.GetByIdOrder(in)
	if err != nil {
		return nil, err
	}
	items, err := p.ProductRepo.OrederGetbyIdQuery(resp.UserId)
	if err != nil {
		return nil, err
	}
	resp.Items = items

	return resp, nil

}

func (p *ProductService) CreatePayment(ctx context.Context, in *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	resp, err := p.OrderRepo.PaymentQuery(in.OrderId)
	if err != nil {
		return nil, err
	}

	var req pb.CreatePaymentResponse
	req.OrderId = in.OrderId
	req.Amount = resp.Amount
	req.Status = resp.Status
	req.PaymentMethod = in.PaymentMethod

	p.PaymentRepo.CreatePayment(&req)

	return p.PaymentRepo.CreatePayment(&req)
}

func (a *ProductService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	return a.AddRepo.CreateCategory(in)
}

func (a *ProductService) GetStatistics(ctx context.Context, in *pb.GetStatisticsRequest) (*pb.GetStatisticsResponse, error) {

	payment, err := a.PaymentRepo.GetStatisticsProduct(in)
	if err != nil {
		return nil, err
	}
	var sum float32
	var products []*pb.TopProduct
	var product pb.TopProduct

	for _, v := range payment {
		sum += v.Revenume
		productId, err := a.OrderRepo.PoductId(v.Id)
		if err != nil {
			return nil, err
		}

		top, err := a.ProductRepo.TopProduct(productId)
		if err != nil {
			return nil, err
		}

		product.Id = top.Id
		product.Name = top.Name
		product.SalesCount = int32(v.Revenume / top.Revenume)
		product.Revenume = v.Revenume
		products = append(products, &product)

	}
	var resp pb.GetStatisticsResponse
	resp.TotalSales = int32(len(products))
	resp.TotalRevenue = sum
	resp.TopProducts = products

	return &resp, err

}

func (p *ProductService) GetProductRecommendations(ctx context.Context, in *pb.GetProductRecommendationsRequest) (*pb.GetProductRecommendationsResponse, error) {
	return p.ProductRepo.GetProductRecommendations(in)
}
