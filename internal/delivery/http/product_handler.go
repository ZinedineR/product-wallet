package http

import (
	"github.com/gin-gonic/gin"
	_ "product-wallet/internal/delivery/http/response"
	"product-wallet/internal/model"
	service "product-wallet/internal/services"
)

type ProductHTTPHandler struct {
	Handler
	ProductService service.ProductService
}

func NewProductHTTPHandler(productService service.ProductService) *ProductHTTPHandler {
	return &ProductHTTPHandler{
		ProductService: productService,
	}
}

// Create godoc
// @Summary Create a new product
// @Description Creates a new product in the catalog
// @Tags Products
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param product body model.CreateProductReq true "Create Product Request"
// @Success 200 {object} response.DataResponse{data=model.CreateProductRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /products [post]
func (h *ProductHTTPHandler) Create(ctx *gin.Context) {
	var request model.CreateProductReq
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	response, errException := h.ProductService.Create(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}

// Update godoc
// @Summary Update an existing product
// @Description Updates product details
// @Tags Products
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param id path string true "uuid format"
// @Param product body model.UpdateProductReq true "Update Product Request"
// @Success 200 {object} response.DataResponse{data=model.UpdateProductRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /products/{id} [put]
func (h *ProductHTTPHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var request model.UpdateProductReq
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	request.ID = id
	response, errException := h.ProductService.Update(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}

// Find godoc
// @Summary Get all products
// @Description Retrieves a list of all products with optional filters, pagination, and sorting
// @Tags Products
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param pageSize query string false "Number of items per page"
// @Param page query string false "Page number"
// @Param filter query string false "Filter rules"
// @Param sort query string false "Sort rules"
// @Success 200 {object} response.DataResponse{data=model.GetAllProductRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /products [get]
func (h *ProductHTTPHandler) Find(ctx *gin.Context) {
	page, sort, filter, err := h.ParsePaginationParams(ctx)
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	request := model.GetAllProductReq{
		Page:   page,
		Filter: filter,
		Sort:   sort,
	}
	response, errException := h.ProductService.Find(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}

// Detail godoc
// @Summary Get product details
// @Description Retrieves the details of a specific product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param id path string true "uuid format"
// @Success 200 {object} response.DataResponse{data=model.GetProductByIDRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /products/{id} [get]
func (h *ProductHTTPHandler) Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	request := model.GetProductByIDReq{
		ID: id,
	}
	response, errException := h.ProductService.Detail(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}

// Delete godoc
// @Summary Delete a product
// @Description Deletes a product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param id path string true "uuid format"
// @Success 200 {object} response.DataResponse{data=model.DeleteProductRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /products/{id} [delete]
func (h *ProductHTTPHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	request := model.DeleteProductReq{
		ID: id,
	}
	response, errException := h.ProductService.Delete(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}
