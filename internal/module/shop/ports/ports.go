package ports

import (
	"codebase-app/internal/module/shop/entity"
	"context"
)

type ShopRepository interface {
	CreateShop(ctx context.Context, req *entity.CreateShopRequest) (*entity.CreateShopResponse, error)
	GetShop(ctx context.Context, req *entity.GetShopRequest) (*entity.GetShopResponse, error)
	DeleteShop(ctx context.Context, req *entity.DeleteShopRequest) error
	UpdateShop(ctx context.Context, req *entity.UpdateShopRequest) (*entity.UpdateShopResponse, error)
	GetShops(ctx context.Context, req *entity.ShopsRequest) (*entity.ShopsResponse, error)
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.ProductResponse, error)
	GetDetailShopAndProduct(ctx context.Context, id string, paginate int, page int) (*entity.DetailShopAndProduct, error)
	GetAllProduct(ctx context.Context, req *entity.ProductFilter) (*entity.ProductsResponse, error)
	GetDetailProduct(ctx context.Context, id string) (*entity.ProductResponse, error)
	DeleteProductByID(ctx context.Context, id string) error
	UpdateProductByID(ctx context.Context, req *entity.UpdateProductRequest) (*entity.UpdateProductRequest, error)
}

type ShopService interface {
	CreateShop(ctx context.Context, req *entity.CreateShopRequest) (*entity.CreateShopResponse, error)
	GetShop(ctx context.Context, req *entity.GetShopRequest) (*entity.GetShopResponse, error)
	DeleteShop(ctx context.Context, req *entity.DeleteShopRequest) error
	UpdateShop(ctx context.Context, req *entity.UpdateShopRequest) (*entity.UpdateShopResponse, error)
	GetShops(ctx context.Context, req *entity.ShopsRequest) (*entity.ShopsResponse, error)
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.ProductResponse, error)
	GetDetailShopAndProduct(ctx context.Context, id string, paginate int, page int) (*entity.DetailShopAndProduct, error)
	GetAllProduct(ctx context.Context, req *entity.ProductFilter) (*entity.ProductsResponse, error)
	GetDetailProduct(ctx context.Context, id string) (*entity.ProductResponse, error)
	DeleteProductByID(ctx context.Context, id string) error
	UpdateProductByID(ctx context.Context, req *entity.UpdateProductRequest) (*entity.UpdateProductRequest, error)
}
