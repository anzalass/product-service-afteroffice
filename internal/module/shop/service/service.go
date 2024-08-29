package service

import (
	"codebase-app/internal/module/shop/entity"
	"codebase-app/internal/module/shop/ports"
	"context"

	"github.com/rs/zerolog/log"
)

var _ ports.ShopService = &shopService{}

type shopService struct {
	repo ports.ShopRepository
}

func NewShopService(repo ports.ShopRepository) *shopService {
	return &shopService{
		repo: repo,
	}
}

func (s *shopService) CreateShop(ctx context.Context, req *entity.CreateShopRequest) (*entity.CreateShopResponse, error) {
	return s.repo.CreateShop(ctx, req)
}

func (s *shopService) GetShop(ctx context.Context, req *entity.GetShopRequest) (*entity.GetShopResponse, error) {
	return s.repo.GetShop(ctx, req)
}

func (s *shopService) DeleteShop(ctx context.Context, req *entity.DeleteShopRequest) error {
	return s.repo.DeleteShop(ctx, req)
}

func (s *shopService) UpdateShop(ctx context.Context, req *entity.UpdateShopRequest) (*entity.UpdateShopResponse, error) {
	return s.repo.UpdateShop(ctx, req)
}

func (s *shopService) GetShops(ctx context.Context, req *entity.ShopsRequest) (*entity.ShopsResponse, error) {
	return s.repo.GetShops(ctx, req)
}
func (s *shopService) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.ProductResponse, error) {
	return s.repo.CreateProduct(ctx, req)
}
func (s *shopService) GetDetailShopAndProduct(ctx context.Context, id string, paginate int, page int) (*entity.DetailShopAndProduct, error) {
	return s.repo.GetDetailShopAndProduct(ctx, id, paginate, page)
}
func (s *shopService) GetAllProduct(ctx context.Context, req *entity.ProductFilter) (*entity.ProductsResponse, error) {
	return s.repo.GetAllProduct(ctx, req)
}
func (s *shopService) GetDetailProduct(ctx context.Context, id string) (*entity.ProductResponse, error) {
	return s.repo.GetDetailProduct(ctx, id)
}
func (s *shopService) DeleteProductByID(ctx context.Context, id string) error {
	return s.repo.DeleteProductByID(ctx, id)
}
func (s *shopService) UpdateProductByID(ctx context.Context, req *entity.UpdateProductRequest) (*entity.UpdateProductRequest, error) {
	log.Debug().Str("id", req.UserID).Msg("repository::Get Detail Product - ID User Input")

	product, err := s.repo.GetDetailProduct(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	log.Debug().Str("id", product.UserID).Msg("repository::Get Detail Product - ID User Product")
	var err1 error
	if product.UserID != req.UserID {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateProductByID - UserID is not equal")
		return nil, err1
	}

	resp, err2 := s.repo.UpdateProductByID(ctx, req)
	if err2 != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateProductByID - failed update product")
		return nil, err1
	}

	return resp, nil
}
