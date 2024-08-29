package entity

import "codebase-app/pkg/types"

type CreateShopRequest struct {
	UserId string `validate:"uuid" db:"user_id"`

	Name        string `json:"name" validate:"required" db:"name"`
	Description string `json:"description" validate:"required,max=255" db:"description"`
	Terms       string `json:"terms" validate:"required" db:"terms"`
}

type CreateShopResponse struct {
	Id string `json:"id" db:"id"`
}

type GetShopRequest struct {
	Id string `validate:"uuid" db:"id"`
}

type GetShopResponse struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Terms       string `json:"terms" db:"terms"`
}

type DeleteShopRequest struct {
	UserId string `prop:"user_id" validate:"uuid" db:"user_id"`

	Id string `validate:"uuid" db:"id"`
}

type UpdateShopRequest struct {
	UserId string `prop:"user_id" validate:"uuid" db:"user_id"`

	Id          string `params:"id" validate:"uuid" db:"id"`
	Name        string `json:"name" validate:"required" db:"name"`
	Description string `json:"description" validate:"required" db:"description"`
	Terms       string `json:"terms" validate:"required" db:"terms"`
}

type UpdateShopResponse struct {
	Id string `json:"id" db:"id"`
}

type ShopsRequest struct {
	UserId   string `prop:"user_id" validate:"uuid"`
	Page     int    `query:"page" validate:"required"`
	Paginate int    `query:"paginate" validate:"required"`
}

func (r *ShopsRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type ShopItem struct {
	Id   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type ShopsResponse struct {
	Items []ShopItem `json:"items"`
	Meta  types.Meta `json:"meta"`
}
type ProductsResponse struct {
	Product         []ProductResponseDashboard `json:"product"`
	Meta            types.Meta                 `json:"meta"`
	KategoriFilter  string                     `json:"kategorifilter" db:"kategori"`
	Name            string                     `json:"namefilter" db:"name"`
	MinHarga        int                        `json:"min_harga" db:"harga"`
	MaxHarga        int                        `json:"max_harga" db:"harga"`
	MerekFilter     string                     `json:"merekfilter" db:"merek"`
	PenilaianFilter int                        `json:"rating" db:"rating"`
	Page            int                        `json:"page" db:"page"`
	Pagination      int                        `json:"pagination" db:"pagination"`
}

type KategoriRequest struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name" db:"name" `
}

type CreateProductRequest struct {
	UserID      string            `validate:"uuid" db:"user_id"`
	ShopID      string            `db:"shop_id" json:"shop_id"`
	Name        string            `validate:"required" json:"name" db:"name"`
	Description string            `json:"description" validate:"required" db:"description"`
	Kategori    []KategoriRequest `validate:"required" json:"kategori"`
	Harga       int               `validate:"required" json:"harga" db:"harga"`
	Stok        int               `validate:"required" json:"stok" db:"stok"`
	Merek       string            `validate:"required" json:"merek" db:"merek"`
}
type ProductResponse struct {
	ID          string            `json:"id" db:"id" validate:"uuid"`
	UserID      string            `validate:"uuid" db:"user_id" json:"user_id"`
	ShopID      string            `validate:"uuid" db:"shop_name" json:"shop_name"`
	Nama        string            `validate:"required" json:"name" db:"name"`
	Description string            `validate:"required" json:"deskripsi" db:"deskripsi"`
	Kategori    []KategoriRequest `validate:"required" json:"kategori" db:"kategori"`
	Harga       int               `validate:"required" json:"harga" db:"harga"`
	Stok        int               `validate:"required" json:"stok" db:"stok"`
	Merek       string            `validate:"required" json:"merek" db:"merek"`
}
type ProductResponseDashboard struct {
	ID        string `json:"id" db:"id" validate:"uuid"`
	UserID    string `validate:"uuid" db:"user_id" json:"user_id"`
	ShopID    string `validate:"uuid" db:"shop_name" json:"shop_name"`
	Nama      string `validate:"required" json:"name" db:"name"`
	Kategori  string `validate:"required" json:"kategori" db:"kategori"`
	Harga     int    `validate:"required" json:"harga" db:"harga"`
	Stok      int    `validate:"required" json:"stok" db:"stok"`
	Penilaian int    `validate:"required" json:"penilaian" db:"penilaian"`
	Merek     string `validate:"required" json:"merek" db:"merek"`
}
type ProductResponseDetail struct {
	ID          string            `json:"id" db:"id" validate:"uuid"`
	Nama        string            `validate:"required" json:"nama" db:"nama"`
	Description string            `validate:"required" json:"deskripsi" db:"deskripsi"`
	Kategori    []KategoriRequest `validate:"required" json:"kategori" db:"kategori"`
	Harga       int               `validate:"required" json:"harga" db:"harga"`
	Stok        int               `validate:"required" json:"stok" db:"stok"`
}

type DetailShopAndProduct struct {
	ShopID        string                  `json:"shop_id"`
	Name          string                  `json:"name"`
	Description   string                  `json:"description"`
	Terms         string                  `json:"terms"`
	Terjual       int                     `json:"terjual"`
	DaftarProduct []ProductResponseDetail `json:"daftar_products"`
	Meta          types.Meta              `json:"meta"`
}

type ProductFilter struct {
	Kategori   string `json:"kategori" db:"kategori"`
	Name       string `json:"name" db:"name"`
	MinHarga   int    `json:"min_harga" db:"harga"`
	MaxHarga   int    `json:"max_harga" db:"harga"`
	Merek      string `json:"merek" db:"merek"`
	Penilaian  int    `json:"rating" db:"rating"`
	Page       int    `json:"page" db:"page"`
	Pagination int    `json:"pagination" db:"pagination"`
}

func (p *ProductFilter) SetDefaultFilter() {
	if p.MaxHarga < 1 {
		p.MaxHarga = 99999999999999
	}
}

type UpdateProductRequest struct {
	ID          string            `prop:"id" db:"id"`
	UserID      string            `json:"user_id" db:"user_id"`
	ShopID      string            `json:"shop_id" db:"shop_id"`
	Name        string            `json:"name" db:"name"`
	Description string            `json:"description" db:"description"`
	Kategori    []KategoriRequest `json:"kategori"`
	Harga       int               `json:"harga" db:"harga"`
	Stok        int               `json:"stok" db:"stok"`
	Merek       string            `json:"merek" db:"merek"`
}

//nama, deskripsi, kategori, harga, dan stok.
