package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Ansalps/genzone-cart-svc/pkg/db"
	"github.com/Ansalps/genzone-cart-svc/pkg/models"
	cartpb "github.com/Ansalps/genzone-cart-svc/pkg/pb"
	productpb "github.com/Ansalps/genzone-product-svc/pkg/pb"
	"google.golang.org/grpc"
)

type Server struct {
	H db.Handler
	cartpb.UnimplementedCartServiceServer
}

func (s *Server) AddToCart(ctx context.Context, req *cartpb.CreateCartRequest) (*cartpb.CreateCartResponse, error) {
	var cart models.Cart
	fmt.Println("req", req)
	cart.UserID = req.Userid
	cart.ProductID = req.Productid
	cart.Qty = uint(req.Quantity)
	// Fetch product details
	product, err := getProductDetails(req.Productid)
	if err != nil {
		log.Printf("Error fetching product details: %v", err)
	}
	// Calculate total amount
	cart.Price = product.Price
	totalAmount := product.Price * float64(req.Quantity)
	cart.Amount = totalAmount
	//var product models.Product
	if result := s.H.DB.Create(&cart); result.Error != nil {
		return &cartpb.CreateCartResponse{
			Status: http.StatusConflict,
			Error:  result.Error.Error(),
		}, nil
	}
	return &cartpb.CreateCartResponse{
		Status: http.StatusCreated,
		Id:     int64(cart.ID),
	}, nil
}
func getProductDetails(productID string) (*productpb.GetProductResponse, error) {
	// Connect to Product Service
	conn, err := grpc.Dial("localhost:50054", grpc.WithInsecure()) // Replace with proper address
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %v", err)
	}
	defer conn.Close()

	client := productpb.NewProductServiceClient(conn)

	// Call GetProduct RPC
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := client.GetProduct(ctx, &productpb.GetProductRequest{ProductId: productID})
	if err != nil {
		return nil, fmt.Errorf("failed to get product details: %v", err)
	}
	fmt.Println("price",response.Price)
	return response, nil
}
