package http

import (
	"github.com/gin-gonic/gin"
	_ "product-wallet/internal/delivery/http/response"
	"product-wallet/internal/model"
	service "product-wallet/internal/services"
)

type WalletHTTPHandler struct {
	Handler
	WalletService service.WalletService
}

// NewWalletHTTPHandler initializes the Wallet HTTP handler
func NewWalletHTTPHandler(walletService service.WalletService) *WalletHTTPHandler {
	return &WalletHTTPHandler{
		WalletService: walletService,
	}
}

// Create godoc
// @Summary Create a new wallet
// @Description Creates a new wallet for the user
// @Tags Wallets
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param wallet body model.CreateWalletReq true "Create Wallet Request"
// @Param wallet body model.CreateWalletReq true "Update Wallet Request"
// @Success 200 {object} response.DataResponse{data=model.CreateWalletRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /wallets [post]
func (h WalletHTTPHandler) Create(ctx *gin.Context) {
	var request model.CreateWalletReq
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	request.UserId = h.ParseGetKey(ctx, "user_id")
	response, errException := h.WalletService.Create(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}

// Update godoc
// @Summary Update an existing wallet
// @Description Updates wallet information
// @Tags Wallets
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param id path string true "Wallet ID"
// @Param wallet body model.UpdateWalletReq true "Update Wallet Request"
// @Success 200 {object} response.DataResponse{data=model.UpdateWalletRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /wallets/{id} [put]
func (h WalletHTTPHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var request model.UpdateWalletReq
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	request.ID = id
	request.UserId = h.ParseGetKey(ctx, "user_id")
	response, errException := h.WalletService.Update(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}

// Find godoc
// @Summary Get all wallets
// @Description Retrieves all wallets for the user
// @Tags Wallets
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param pageSize query string false "Number of items per page"
// @Param page query string false "Page number"
// @Param filter query string false "Filter rules<br><br>### Rules Filter<br>rule:<br>  * {Name of Field}:{value}:{Symbol}<br><br>Symbols:<br>  * eq (=)<br>  * lt (<)<br>  * gt (>)<br>  * lte (<=)<br>  * gte (>=)<br>  * in (in)<br>  * like (like)"
// @Param sort query string false "Sort rules:<br><br>### Rules Sort<br>rule:<br>  * {Name of Field}:{Symbol}<br><br>Symbols:<br>  * asc<br>  * desc<br><br>"
// @Success 200 {object} response.DataResponse{data=model.GetAllWalletRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /wallets [get]
func (h WalletHTTPHandler) Find(ctx *gin.Context) {
	page, sort, filter, err := h.ParsePaginationParams(ctx)
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	request := model.GetAllWalletReq{
		Page:   page,
		Filter: filter,
		Sort:   sort,
	}
	response, errException := h.WalletService.Find(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}

// Detail godoc
// @Summary Get wallet details
// @Description Retrieves details of a specific wallet by ID
// @Tags Wallets
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param id path string true "Wallet ID"
// @Success 200 {object} response.DataResponse{data=model.GetWalletByIDRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /wallets/{id} [get]
func (h WalletHTTPHandler) Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	request := model.GetWalletByIDReq{
		ID: id,
	}
	response, errException := h.WalletService.Detail(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}

// DetailWalletTransaction godoc
// @Summary Get wallet transaction details
// @Description Retrieves details of a specific wallet, including transactions within a date range
// @Tags Wallets
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param id path string true "Wallet ID"
// @Param from query string false "Start date for transactions in YYYY-MM-DD format"
// @Param to query string false "End date for transactions in YYYY-MM-DD format"
// @Success 200 {object} response.DataResponse{data=model.GetWalletByTransactionRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /wallets/transaction/{id} [get]
func (h WalletHTTPHandler) DetailWalletTransaction(ctx *gin.Context) {
	// Extract wallet ID from the path
	id := ctx.Param("id")

	// Parse from and to date parameters
	fromDate, toDate, err := h.ParseDateParam(ctx)
	if err != nil {
		h.BadRequestJSON(ctx, "Invalid date format. Use YYYY-MM-DD.")
		return
	}

	// Create the request for the service
	request := model.GetWalletByTransactionReq{
		ID:   id,
		From: fromDate,
		To:   toDate,
	}

	// Call the service method
	response, errException := h.WalletService.DetailWalletTransaction(ctx, request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	// Return the response as JSON
	h.DataJSON(ctx, response)
}

// Delete godoc
// @Summary Delete a wallet
// @Description Deletes a wallet by ID
// @Tags Wallets
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param id path string true "Wallet ID"
// @Success 200 {object} response.DataResponse{data=model.DeleteWalletRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /wallets/{id} [delete]
func (h WalletHTTPHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	request := model.DeleteWalletReq{
		ID: id,
	}
	response, errException := h.WalletService.Delete(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}
