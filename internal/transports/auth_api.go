package handler_api

import (
	"net/http"

	grpcLog "github.com/o-sokol-o/grpc_log_server/pkg/domain"

	"github.com/AquaEngineering/AquaHub/internal/domain"
	"github.com/AquaEngineering/AquaHub/pkg/jwt_processing"

	"github.com/gin-gonic/gin"
)

// Обработчик регистрации

// @Summary     SignUp
// @Tags        Authentication
// @Description create account
// @ID          create-account
// @Accept      json
// @Produce     json
// @Param       input   body      domain.User true "account info"
// @Success     200     {integer} integer     1
// @Failure     400,404 {object} statusResponse
// @Failure     500     {object} statusResponse
// @Failure     default {object} statusResponse
// @Router      /auth/sign-up [post]
func (h *Handler) signUp(ctx *gin.Context) {

	// В структуру input будем записывать принятый json от пользователя
	var input domain.User

	// У gin.Context есть метод BindJSON, принимающий ссылку на структуру
	// в которую кладёт распарсенный и предварительно валидированный принятый json
	// Или 400 - ошибка на стороне клиента
	if err := ctx.BindJSON(&input); err != nil {
		h.newErrorResponse(ctx, http.StatusBadRequest, "User send invalid input body: "+err.Error()) // http.StatusBadRequest = 400
		return
	}

	if err := input.Validate(); err != nil {
		h.newErrorResponse(ctx, http.StatusBadRequest, "User send invalid input body: "+err.Error()) // http.StatusBadRequest = 400
		return
	}

	// Принимает структуру User в качестве аргумента
	// Возвращает ID созданного нового пользователя
	id, err := h.serviceAuthentications.CreateUser(input)
	if err != nil {
		// внутренняя ошибка на сервере
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error()) // http.StatusInternalServerError = 500
		return
	}

	jwt_token, err := jwt_processing.GenerateToken(id)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error()) // http.StatusInternalServerError = 500
		return
	}

	h.GrpcLog.Send(grpcLog.LogRequest_CREATE, grpcLog.LogRequest_USER, int64(id))

	// Если пользователь существует, то в ответе получаем токен.
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": "Bearer " + jwt_token,
	})
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Обработчик аутентификации

// @Summary     SignIn
// @Tags        Authentication
// @Description login
// @ID          login
// @Accept      json
// @Produce     json
// @Param       input   body     signInInput true "credentials"
// @Success     200     {string} string      "token"
// @Failure     400,404 {object} statusResponse
// @Failure     500     {object} statusResponse
// @Failure     default {object} statusResponse
// @Router      /auth/sign-in [post]
func (h *Handler) signIn(ctx *gin.Context) {

	var input signInInput

	if err := ctx.BindJSON(&input); err != nil {
		h.newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var jwt_token string
	user, err := h.serviceAuthentications.Authenticate(input.Email, input.Password)
	if err != nil {
		// Возвращаем ошибку SQL из БД
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	jwt_token, err = jwt_processing.GenerateToken(user.ID)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	h.GrpcLog.Send(grpcLog.LogRequest_CREATE, grpcLog.LogRequest_USER, int64(user.ID))

	// Если пользователь существует, то в ответе получаем токен.
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": "Bearer " + jwt_token,
	})
}
