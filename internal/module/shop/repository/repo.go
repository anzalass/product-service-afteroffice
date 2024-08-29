package repository

import (
	"codebase-app/internal/module/shop/entity"
	"codebase-app/internal/module/shop/ports"
	"context"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ShopRepository = &shopRepository{}

type shopRepository struct {
	db *sqlx.DB
}

func NewShopRepository(db *sqlx.DB) *shopRepository {
	return &shopRepository{
		db: db,
	}
}

func (r *shopRepository) CreateShop(ctx context.Context, req *entity.CreateShopRequest) (*entity.CreateShopResponse, error) {
	var resp = new(entity.CreateShopResponse)
	// Your code here
	query := `
		INSERT INTO shops (user_id, name, description, terms)
		VALUES (?, ?, ?, ?) RETURNING id
	`

	err := r.db.QueryRowContext(ctx, r.db.Rebind(query),
		req.UserId,
		req.Name,
		req.Description,
		req.Terms).Scan(&resp.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::CreateShop - Failed to create shop")
		return nil, err
	}

	return resp, nil
}

func (r *shopRepository) GetShop(ctx context.Context, req *entity.GetShopRequest) (*entity.GetShopResponse, error) {
	var resp = new(entity.GetShopResponse)
	// Your code here
	query := `
		SELECT name, description, terms
		FROM shops
		WHERE id = ? AND deleted_at is NULL
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Id).StructScan(resp)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetShop - Failed to get shop")
		return nil, err
	}

	return resp, nil
}

func (r *shopRepository) DeleteShop(ctx context.Context, req *entity.DeleteShopRequest) error {
	queryshop := `UPDATE shops SET deleted_at = NOW() WHERE id = ? AND user_id = ?`
	queryproduct := `UPDATE product SET deleted_at = NOW() WHERE shop_id = ? AND user_id = ?`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(queryshop), req.Id, req.UserId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DeleteShop  - Failed to delete shop")
		return err
	}

	_, err2 := r.db.ExecContext(ctx, r.db.Rebind(queryproduct), req.Id, req.UserId)
	if err2 != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DeleteShop Product - Failed to delete shop product")
		return err
	}

	return nil
}

func (r *shopRepository) UpdateShop(ctx context.Context, req *entity.UpdateShopRequest) (*entity.UpdateShopResponse, error) {
	var resp = new(entity.UpdateShopResponse)

	query := `
		UPDATE shops
		SET name = ?, description = ?, terms = ?, updated_at = NOW()
		WHERE id = ? AND user_id = ?
		RETURNING id
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query),
		req.Name,
		req.Description,
		req.Terms,
		req.Id,
		req.UserId).Scan(&resp.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::UpdateShop - Failed to update shop")
		return nil, err
	}

	return resp, nil
}

func (r *shopRepository) GetShops(ctx context.Context, req *entity.ShopsRequest) (*entity.ShopsResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.ShopItem
	}

	var (
		resp = new(entity.ShopsResponse)
		data = make([]dao, 0, req.Paginate)
	)
	resp.Items = make([]entity.ShopItem, 0, req.Paginate)

	query := `
		SELECT
			COUNT(id) OVER() as total_data,
			id,
			name
		FROM shops
		WHERE
			deleted_at IS NULL
			AND user_id = ?
		LIMIT ? OFFSET ?
	`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query),
		req.UserId,
		req.Paginate,
		req.Paginate*(req.Page-1),
	)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetShops - Failed to get shops")
		return nil, err
	}

	if len(data) > 0 {
		resp.Meta.TotalData = data[0].TotalData
	}

	for _, d := range data {
		resp.Items = append(resp.Items, d.ShopItem)
	}

	resp.Meta.CountTotalPage(req.Page, req.Paginate, resp.Meta.TotalData)

	return resp, nil
}

func (r *shopRepository) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.ProductResponse, error) {
	var resp = new(entity.ProductResponse)

	queryproduct := `INSERT INTO product (user_id, shop_id, name, description, harga, stok, merek) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id, user_id, shop_id, name, description, harga, stok, merek`
	err1 := r.db.QueryRowContext(ctx, r.db.Rebind(queryproduct),
		req.UserID,
		req.ShopID,
		req.Name,
		req.Description,
		req.Harga,
		req.Stok,
		req.Merek,
	).Scan(&resp.ID, &resp.UserID, &resp.ShopID, &resp.Nama, &resp.Description, &resp.Harga, &resp.Stok, &resp.Merek)
	if err1 != nil {
		log.Error().Err(err1).Any("payload", req).Msg("repository::CreateProduct - Failed to create product")
		return nil, err1
	}

	queryKategori := `INSERT INTO kategori (product_id, name) VALUES (?, ?) RETURNING name`
	for _, v := range req.Kategori {
		_, err2 := r.db.ExecContext(ctx, r.db.Rebind(queryKategori), resp.ID, v.Name)

		if err2 != nil {
			log.Error().Err(err2).Any("payload", req).Msg("repository::CreateKategori - Failed to create kategori")
			return nil, err2
		}

	}

	resp.Kategori = req.Kategori

	return resp, nil

}

func (r *shopRepository) GetDetailShopAndProduct(ctx context.Context, id string, paginate int, page int) (*entity.DetailShopAndProduct, error) {
	var resp = new(entity.DetailShopAndProduct)

	type daoshop struct {
		Name        string `db:"name"`
		Description string `db:"description"`
		Terms       string `db:"terms"`
	}
	type daoproduct struct {
		ProductName       string `db:"product_name"`
		KategoriProductID string `db:"kategori_productid"`
		ProductDesc       string `db:"product_description"`
		ProductHarga      int    `db:"product_harga"`
		ProductStok       int    `db:"product_stok"`
		KategoriNames     string `db:"kategori_name"`
	}

	var datashop []daoshop
	var dataproduct []daoproduct
	var shopErr, productErr error

	// Channel untuk menampung hasil query
	shopChan := make(chan []daoshop, 1)
	productChan := make(chan []daoproduct, 1)

	// Goroutine untuk menjalankan query shop
	go func() {
		defer close(shopChan)
		shopErr = r.db.SelectContext(ctx, &datashop, r.db.Rebind(`SELECT name, description, terms FROM shops WHERE id = ?`), id)
		if shopErr != nil {
			log.Error().Err(shopErr).Any("payload", id).Msg("repository::GetDetailShopAndProduct Shop - Failed to get Get Detail Shop And Product")
		}
		shopChan <- datashop
	}()

	// Goroutine untuk menjalankan query product
	go func() {
		defer close(productChan)
		productErr = r.db.SelectContext(ctx, &dataproduct, r.db.Rebind(`SELECT product.name as product_name, product.description as product_description, product.harga as product_harga,
			product.stok as product_stok, kategori.product_id as kategori_productid, kategori.name as kategori_name 
			FROM product JOIN kategori ON product.id = kategori.product_id 
			WHERE shop_id = ? LIMIT ? OFFSET ?`), id, 4, 4*(page-1))
		if productErr != nil {
			log.Error().Err(productErr).Any("payload", id).Msg("repository::GetDetailShopAndProduct Product - Failed to get Get Detail Shop And Product")
		}
		productChan <- dataproduct
	}()

	// Menunggu hasil dari goroutine
	datashop = <-shopChan
	dataproduct = <-productChan

	if len(datashop) > 0 {
		resp.Name = datashop[0].Name
		resp.Description = datashop[0].Description
		resp.Terms = datashop[0].Terms
	}
	resp.Terjual = 0

	productMap := make(map[string]*entity.ProductResponseDetail)

	for _, row := range dataproduct {
		if _, exists := productMap[row.ProductName]; !exists {
			productMap[row.ProductName] = &entity.ProductResponseDetail{
				ID:          row.ProductName,
				Nama:        row.ProductName,
				Description: row.ProductDesc,
				Harga:       row.ProductHarga,
				Stok:        row.ProductStok,
			}
		}

		productMap[row.ProductName].Kategori = append(productMap[row.ProductName].Kategori, entity.KategoriRequest{
			ProductID: row.KategoriProductID,
			Name:      row.KategoriNames,
		})
	}

	for _, product := range productMap {
		resp.DaftarProduct = append(resp.DaftarProduct, *product)
	}
	resp.Meta.Page = page
	resp.Meta.Paginate = 4
	resp.Meta.TotalData = len(resp.DaftarProduct)
	resp.Meta.CountTotalPage(resp.Meta.Page, resp.Meta.Paginate, resp.Meta.TotalData)
	return resp, nil
}

func (r *shopRepository) GetAllProduct(ctx context.Context, req *entity.ProductFilter) (*entity.ProductsResponse, error) {
	type dao struct {
		entity.ProductResponseDashboard
	}

	var data []dao
	var resp = new(entity.ProductsResponse)
	resp.Product = make([]entity.ProductResponseDashboard, 0, req.Pagination)

	var query = `SELECT DISTINCT
				product.id AS id,
				product.user_id AS user_id,
				shops.name AS shop_name,
				product.name AS name, 
				product.harga AS harga, 
				product.penilaian AS penilaian, 
				product.merek AS merek,
				product.stok AS stok,
				kategori.name AS kategori
			FROM 
				product
			JOIN 
				kategori ON product.id = kategori.product_id
			JOIN 
				shops ON shops.id = product.shop_id
			WHERE
				product.merek ILIKE '%' || ? || '%'
				AND product.name ILIKE '%' || ? || '%'
				AND CAST(product.harga AS numeric) >= ?
				AND CAST(product.harga AS numeric) <= ? 
				AND kategori.name ILIKE '%' || ? || '%'
				AND product.penilaian = ?
				AND product.deleted_at IS NULL
			LIMIT ? OFFSET ?;`

	var query2 = `SELECT DISTINCT
				product.id AS id,
				product.user_id AS user_id,
				shops.name AS shop_name,
				product.name AS name, 
				product.harga AS harga, 
				product.penilaian AS penilaian, 
				product.merek AS merek,
				product.stok AS stok,
				kategori.name AS kategori
			FROM 
				product
			JOIN 
				kategori ON product.id = kategori.product_id
			JOIN 
				shops ON shops.id = product.shop_id
			WHERE
				product.merek ILIKE '%' || ? || '%'
				AND product.name ILIKE '%' || ? || '%'
				AND CAST(product.harga AS numeric) >= ?
				AND CAST(product.harga AS numeric) <= ? 
				AND kategori.name ILIKE '%' || ? || '%'
				AND product.deleted_at IS NULL
			LIMIT ? OFFSET ?;`

	req.SetDefaultFilter()

	if req.Penilaian < 1 {
		err := r.db.SelectContext(ctx, &data, r.db.Rebind(query2),
			req.Merek,
			req.Name,
			req.MinHarga,
			req.MaxHarga,
			req.Kategori,

			req.Pagination,
			req.Pagination*(req.Page-1))

		if err != nil {
			log.Error().Err(err).Msg("repository::GetAllProducts - Failed to get products")
			return nil, err
		}
	} else if req.Penilaian > 0 {

		err := r.db.SelectContext(ctx, &data, r.db.Rebind(query),
			req.Merek,
			req.Name,
			req.MinHarga,
			req.MaxHarga,
			req.Kategori,
			req.Penilaian,
			req.Pagination,
			req.Pagination*(req.Page-1))

		if err != nil {
			log.Error().Err(err).Msg("repository::GetAllProducts - Failed to get products")
			return nil, err
		}
	}

	productMap := make(map[string]*entity.ProductResponseDashboard)

	for _, row := range data {
		if _, exists := productMap[row.ID]; !exists {
			productMap[row.ID] = &entity.ProductResponseDashboard{
				ID:        row.ID,
				UserID:    row.UserID,
				ShopID:    row.ShopID,
				Nama:      row.Nama,
				Kategori:  row.Kategori,
				Merek:     row.Merek,
				Penilaian: row.Penilaian,
				Harga:     row.Harga,
				Stok:      row.Stok,
			}
		}

	}

	for _, product := range productMap {
		resp.Product = append(resp.Product, *product)
	}

	resp.KategoriFilter = req.Kategori
	resp.MerekFilter = req.Merek
	resp.Name = req.Name
	resp.MaxHarga = req.MaxHarga
	resp.MinHarga = req.MinHarga
	resp.Page = req.Page
	resp.Pagination = req.Pagination
	resp.PenilaianFilter = req.Penilaian

	return resp, nil
}

func (r *shopRepository) GetDetailProduct(ctx context.Context, id string) (*entity.ProductResponse, error) {
	resp := &entity.ProductResponse{}

	type dao struct {
		ID          string `db:"id_product"`
		UserID      string `db:"product_user_id"`
		NamaToko    string `db:"nama_toko"`
		Name        string `db:"name_product"`
		Harga       string `db:"harga_product"`
		Description string `db:"description_product"`
		Stok        int    `db:"stok_product"`
		Rating      int    `db:"rating"`
		Kategori    string `db:"kategori_product"`
		Merek       string `db:"merek_product"`
	}

	var data []dao

	query := `select product.id as id_product,
					product.user_id as product_user_id,
					 shops.name as nama_toko, 
					 product.name as name_product, 
					 product.harga as harga_product, 
					 product.description as description_product, 
					 product.stok as stok_product, 
					 product.penilaian as rating,
					 product.merek as merek_product,
					 kategori.name as kategori_product
				from product
				join shops on shops.id = product.shop_id
				join kategori on product.id = kategori.product_id
				where product.id = ? and product.deleted_at is null`

	log.Debug().Str("id", id).Msg("repository::Get Detail Product - ID Value")

	if id == "" {
		log.Error().Msg("repository::GetDetailProduct - ID is empty")
		return nil, fmt.Errorf("invalid ID: ID cannot be empty")
	}

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), id)
	if err != nil {
		log.Error().Err(err).Any("payload", id).Msg("repository::Get Detail Product - Failed to get shops")
		return nil, err
	}

	for _, d := range data {
		resp.Kategori = append(resp.Kategori, entity.KategoriRequest{
			ProductID: d.ID,
			Name:      d.Kategori,
		})
	}

	harga, _ := strconv.Atoi(data[0].Harga)
	resp.UserID = data[0].UserID
	resp.Nama = data[0].Name
	resp.ShopID = data[0].NamaToko
	resp.Harga = harga
	resp.Description = data[0].Description
	resp.Merek = data[0].Merek
	resp.Stok = data[0].Stok
	resp.ID = data[0].ID

	return resp, nil

}

func (r *shopRepository) DeleteProductByID(ctx context.Context, id string) error {
	queryproduct := `update product set deleted_at = NOW() where id = ?`
	querykategori := `update kategori set deleted_at = NOW() where product_id = ?`

	errChan := make(chan error, 2)

	// Jalankan querykategori di goroutine
	go func() {
		_, err := r.db.ExecContext(ctx, r.db.Rebind(querykategori), id)
		if err != nil {
			log.Error().Err(err).Any("payload", id).Msg("repository::DeleteShop - Failed to delete shop")
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// Jalankan queryproduct di goroutine
	go func() {
		_, err := r.db.ExecContext(ctx, r.db.Rebind(queryproduct), id)
		if err != nil {
			log.Error().Err(err).Any("payload", id).Msg("repository::DeleteShop Product - Failed to delete shop product")
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// Tunggu hasil dari kedua goroutine
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

func (r *shopRepository) UpdateProductByID(ctx context.Context, req *entity.UpdateProductRequest) (*entity.UpdateProductRequest, error) {
	var resp = new(entity.UpdateProductRequest)
	queryproduct := `update product set name = ?, description = ?, harga = ?, stok = ?, merek = ? where id = ? and deleted_at is null returning id, name, description, harga, stok, merek`

	err1 := r.db.QueryRowContext(ctx, r.db.Rebind(queryproduct),
		req.Name,
		req.Description,
		req.Harga,
		req.Stok,
		req.Merek,
		req.ID).Scan(
		&resp.ID,
		&resp.Name,
		&resp.Description,
		&resp.Harga,
		&resp.Stok,
		&resp.Merek)
	if err1 != nil {
		log.Error().Err(err1).Any("payload", req).Msg("repository::UpdateProductByID - Failed to update product")
		return nil, err1
	}

	// Delete existing categories for the product
	queryDelete := `DELETE FROM kategori WHERE product_id = ?`
	_, errDelete := r.db.ExecContext(ctx, r.db.Rebind(queryDelete), req.ID)
	if errDelete != nil {
		log.Error().Err(errDelete).Any("payload", req).Msg("repository::UpdateProductByID kategori - Failed to delete existing categories")
		return nil, errDelete
	}

	// Insert new categories
	for _, v := range req.Kategori {
		queryInsert := `INSERT INTO kategori (name, product_id) VALUES (?, ?)`
		_, errInsert := r.db.ExecContext(ctx, r.db.Rebind(queryInsert), v.Name, req.ID)
		if errInsert != nil {
			log.Error().Err(errInsert).Any("payload", req).Msg("repository::UpdateProductByID kategori - Failed to insert new categories")
			return nil, errInsert
		}
		resp.Kategori = append(resp.Kategori, entity.KategoriRequest{
			ProductID: req.ID,
			Name:      v.Name,
		})
	}

	return resp, nil

}
