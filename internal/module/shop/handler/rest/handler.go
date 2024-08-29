package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/middleware"
	"codebase-app/internal/module/shop/entity"
	"codebase-app/internal/module/shop/ports"
	"codebase-app/internal/module/shop/repository"
	"codebase-app/internal/module/shop/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type shopHandler struct {
	service ports.ShopService
}

func NewShopHandler() *shopHandler {
	var (
		handler = new(shopHandler)
		repo    = repository.NewShopRepository(adapter.Adapters.ShopeefunPostgres)
		service = service.NewShopService(repo)
	)
	handler.service = service

	return handler
}

func (h *shopHandler) Register(router fiber.Router) {
	router.Get("/shops", middleware.UserIdHeader, h.GetShops)
	router.Get("/shops/:id", h.GetShop)
	router.Post("/shops", middleware.UserIdHeader, h.CreateShop)
	router.Delete("/shops/:id", middleware.UserIdHeader, h.DeleteShop)
	router.Patch("/shops/:id", middleware.UserIdHeader, h.UpdateShop)
	router.Post("/product", middleware.UserIdHeader, h.CreateProduct)
	router.Post("/detailshop/:id", h.GetDetailShopAndProduct)
	router.Post("/product-all", h.GetAllProduct)
	router.Get("/product/:id", h.GetDetailProduct)
	router.Patch("/delete/:id", h.DeleteProductByID)
	router.Put("/update/:id", middleware.UserIdHeader, h.UpdateProductByID)

}

func (h *shopHandler) CreateShop(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateShopRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateShop - Parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = l.UserId

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::CreateShop - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.CreateShop(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))

}

func (h *shopHandler) GetShop(c *fiber.Ctx) error {
	var (
		req = new(entity.GetShopRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetShop - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetShop(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *shopHandler) DeleteShop(c *fiber.Ctx) error {
	var (
		req = new(entity.DeleteShopRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)
	req.UserId = l.UserId
	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::DeleteShop - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	err := h.service.DeleteShop(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}

func (h *shopHandler) UpdateShop(c *fiber.Ctx) error {
	var (
		req = new(entity.UpdateShopRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::UpdateShop - Parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = l.UserId
	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::UpdateShop - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.UpdateShop(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *shopHandler) GetShops(c *fiber.Ctx) error {
	var (
		req = new(entity.ShopsRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetShops - Parse request query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.UserId = l.UserId
	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetShops - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetShops(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))

}

func (h *shopHandler) CreateProduct(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateProductRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserID = l.UserId

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateProduct - Parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::CreateProduct - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.CreateProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}
	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))

}

func (h *shopHandler) GetDetailShopAndProduct(c *fiber.Ctx) error {
	var (
		id            = c.Params("id")
		querypage     = c.Query("page", "10")
		querypaginate = c.Query("paginate", "10")
		ctx           = c.Context()
	)

	page, _ := strconv.Atoi(querypage)
	paginate, _ := strconv.Atoi(querypaginate)

	resp, err := h.service.GetDetailShopAndProduct(ctx, id, paginate, page)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}
	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))
}

func (h *shopHandler) GetAllProduct(c *fiber.Ctx) error {
	var (
		req = new(entity.ProductFilter)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateProduct - Parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	resp, err := h.service.GetAllProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}
	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))

}
func (h *shopHandler) GetDetailProduct(c *fiber.Ctx) error {
	var (
		req = c.Params("id")
		ctx = c.Context()
	)

	resp, err := h.service.GetDetailProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}
	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))

}
func (h *shopHandler) DeleteProductByID(c *fiber.Ctx) error {
	var (
		req = c.Params("id")
		ctx = c.Context()
	)

	err := h.service.DeleteProductByID(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}
	return c.Status(fiber.StatusCreated).JSON(response.Success("Berhasil menghapus", ""))

}

func (h *shopHandler) UpdateProductByID(c *fiber.Ctx) error {
	var (
		req = new(entity.UpdateProductRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
		l   = middleware.GetLocals(c)
	)

	req.UserID = l.UserId
	req.ID = c.Params("id")

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateProduct - Parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::CreateProduct - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.UpdateProductByID(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}
	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))

}
