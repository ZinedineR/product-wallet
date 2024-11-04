package http

import (
	"github.com/gin-gonic/gin"
	_ "product-wallet/internal/delivery/http/response"
	"product-wallet/internal/model"
	service "product-wallet/internal/services"
)

type TransactionHTTPHandler struct {
	Handler
	TransactionService service.TransactionService
}

func NewTransactionHTTPHandler(transactionService service.TransactionService) *TransactionHTTPHandler {
	return &TransactionHTTPHandler{
		TransactionService: transactionService,
	}
}

// Create godoc
// @Summary Create a new transaction
// @Description Creates a new transaction record
// @Tags Transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param transaction body model.CreateTransactionReq true "Create Transaction Request"
// @Success 200 {object} response.DataResponse{data=model.CreateTransactionRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /transactions [post]
func (h *TransactionHTTPHandler) Create(ctx *gin.Context) {
	var request model.CreateTransactionReq
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	response, errException := h.TransactionService.Create(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}

// Detail godoc
// @Summary Get transaction details
// @Description Retrieves details of a specific transaction by ID
// @Tags Transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param id path string true "Transaction ID"
// @Success 200 {object} response.DataResponse{data=model.GetTransactionByIDRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /transactions/{id} [get]
func (h *TransactionHTTPHandler) Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	request := model.GetTransactionByIDReq{
		ID: id,
	}
	response, errException := h.TransactionService.Detail(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}

// Find godoc
// @Summary Get all transactions
// @Description Retrieves all transactions with optional filters, pagination, and sorting
// @Tags Transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param pageSize query string false "Number of items per page"
// @Param page query string false "Page number"
// @Param filter query string false "Filter rules<br><br>### Rules Filter<br>rule:<br>  * {Name of Field}:{value}:{Symbol}<br><br>Symbols:<br>  * eq (=)<br>  * lt (<)<br>  * gt (>)<br>  * lte (<=)<br>  * gte (>=)<br>  * in (in)<br>  * like (like)"
// @Param sort query string false "Sort rules:<br><br>### Rules Sort<br>rule:<br>  * {Name of Field}:{Symbol}<br><br>Symbols:<br>  * asc<br>  * desc<br><br>"
// @Success 200 {object} response.DataResponse{data=model.GetAllTransactionRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /transactions [get]
func (h *TransactionHTTPHandler) Find(ctx *gin.Context) {
	// Parse pagination, sorting, and filter parameters
	page, sort, filter, err := h.ParsePaginationParams(ctx)
	if err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}

	// Construct the request object
	request := model.GetAllTransactionReq{
		Page:   page,
		Filter: filter,
		Sort:   sort,
	}

	// Call the Find method on TransactionService
	response, errException := h.TransactionService.Find(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	// Return the response data as JSON
	h.DataJSON(ctx, response)
}

// Credit godoc
// @Summary Credit transaction
// @Description Credits a specific wallet with an amount
// @Tags Transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param credit body model.CreditTransactionReq true "Credit Transaction Request"
// @Success 200 {object} response.DataResponse{data=model.CreditTransactionRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /transactions/credit [post]
func (h *TransactionHTTPHandler) Credit(ctx *gin.Context) {
	var request model.CreditTransactionReq
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	response, errException := h.TransactionService.Credit(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}

// Transfer godoc
// @Summary Transfer transaction
// @Description Transfers an amount from one wallet to another
// @Tags Transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param transfer body model.TransferTransactionReq true "Transfer Transaction Request"
// @Success 200 {object} response.DataResponse{data=model.TransferTransactionRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /transactions/transfer [post]
func (h *TransactionHTTPHandler) Transfer(ctx *gin.Context) {
	var request model.TransferTransactionReq
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	response, errException := h.TransactionService.Transfer(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}

// Delete godoc
// @Summary Delete a transaction
// @Description Deletes a transaction by ID
// @Tags Transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization JWT input: Bearer <Token>"
// @Param id path string true "Transaction ID"
// @Success 200 {object} response.DataResponse{data=model.DeleteTransactionRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /transactions/{id} [delete]
func (h *TransactionHTTPHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	request := model.DeleteTransactionReq{
		ID: id,
	}
	response, errException := h.TransactionService.Delete(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}
	h.DataJSON(ctx, response)
}
