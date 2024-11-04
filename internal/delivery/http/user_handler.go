package http

import (
	"github.com/gin-gonic/gin"
	_ "product-wallet/internal/delivery/http/response"
	"product-wallet/internal/model"
	service "product-wallet/internal/services"
)

type UserHTTPHandler struct {
	Handler
	UserService service.UserService
}

func NewUserHTTPHandler(user service.UserService) *UserHTTPHandler {
	return &UserHTTPHandler{
		UserService: user,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Registers a new user with the provided username and password
// @Tags Users
// @Accept json
// @Produce json
// @Param register body model.CreateUserReq true "Registration Request"
// @Success 200 {object} response.DataResponse{data=model.CreateUserRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /auth/register [post]
func (h UserHTTPHandler) Register(ctx *gin.Context) {
	request := model.CreateUserReq{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	response, errException := h.UserService.Register(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, response)
}

// Login godoc
// @Summary User login
// @Description Authenticates the user and returns an access token
// @Tags Users
// @Accept json
// @Produce json
// @Param login body model.CreateUserReq true "Login Request"
// @Success 200 {object} response.DataResponse{data=model.LoginUserRes} "success"
// @Failure 400 {object} response.DataResponse "error"
// @Router /auth/login [post]
func (h UserHTTPHandler) Login(ctx *gin.Context) {
	request := model.CreateUserReq{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.BadRequestJSON(ctx, err.Error())
		return
	}
	result, errException := h.UserService.Login(ctx, &request)
	if errException != nil {
		h.ExceptionJSON(ctx, errException)
		return
	}

	h.DataJSON(ctx, result)
}
