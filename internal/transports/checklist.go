package handler_api

import (
	"net/http"
	"strconv"

	"github.com/AquaEngineering/AquaHub/internal/domain"

	grpcLog "github.com/o-sokol-o/grpc_log_server/pkg/domain"

	"github.com/gin-gonic/gin"
)

type ChecklistsResponse struct {
	Data []domain.Checklist `json:"data"`
}

// @Summary     Create Checklist
// @Security    ApiKeyAuth
// @Tags        Checklists
// @Description create checklist
// @ID          create-list
// @Accept      json
// @Produce     json
// @Param       input   body      domain.UpdateChecklist true "Checklist info"
// @Success     200     {object}  idResponse
// @Failure     400,404 {object} statusResponse
// @Failure     500     {object} statusResponse
// @Failure     default {object} statusResponse
// @Router      /api/lists [post]
func (h *Handler) createList(ctx *gin.Context) {

	// Получение ID пользователя из контекста и последующий его вывод в response
	userId, err := getUserIdFromContext(ctx)
	if err != nil {
		// добавим обработку случая когда в контексте нет id пользователя
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var input domain.UpdateChecklist
	if err := ctx.BindJSON(&input); err != nil {
		h.newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Вызовем метод создания списка,
	// в который передадим наши данные, полученные из токена аутентификации и тела запроса.
	id, err := h.serviceChecklist.Create(userId, input)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	h.GrpcLog.Send(grpcLog.LogRequest_CREATE, grpcLog.LogRequest_CHECKLIST, int64(id))

	// Добавим тело ответа при успешном запросе, в котором будем возвращать ID созданного списка.
	ctx.JSON(http.StatusOK, idResponse{
		ID: id,
	})
}

// @Summary     Get All Checklists
// @Security    ApiKeyAuth
// @Tags        Checklists
// @Description get all lists
// @ID          get-all-lists
// @Accept      json
// @Produce     json
// @Success     200     {object} ChecklistsResponse
// @Failure     400,404 {object} statusResponse
// @Failure     500     {object} statusResponse
// @Failure     default {object} statusResponse
// @Router      /api/lists [get]
func (h *Handler) getAllLists(ctx *gin.Context) {

	userId, err := getUserIdFromContext(ctx)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	lists, err := h.serviceChecklist.GetAllChecklistOfUser(userId)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	h.GrpcLog.Send(grpcLog.LogRequest_GET, grpcLog.LogRequest_CHECKLIST, 0)

	// Для response используем дополнительную структуру getAllListsResponse,
	// в которой будет поле дата, типа слайса списков.
	// Добавим ответ
	ctx.JSON(http.StatusOK, ChecklistsResponse{
		Data: lists,
	})
}

// @Summary     Get Checklist By Id
// @Security    ApiKeyAuth
// @Tags        Checklists
// @Description get list by id
// @ID          get-list-by-id
// @Accept      json
// @Produce     json
// @Param 		id path int true "Checklist ID"
// @Success     200     {object} domain.Checklist
// @Failure     400,404 {object} statusResponse
// @Failure     500     {object} statusResponse
// @Failure     default {object} statusResponse
// @Router      /api/lists/{id} [get]
func (h *Handler) getListById(ctx *gin.Context) {

	// Реализуем логику получения списка по id.

	userId, err := getUserIdFromContext(ctx)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// Вызовем у контекста метод Param, указав в качестве аргумента имя параметра.
	// И обернем эту функцию в Atoi, из стандартной библиотеки strconv, для преобразования строки в число.
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		h.newErrorResponse(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.serviceChecklist.GetById(userId, id)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	h.GrpcLog.Send(grpcLog.LogRequest_GET, grpcLog.LogRequest_CHECKLIST, int64(id))

	ctx.JSON(http.StatusOK, list)
}

// @Summary     Update Checklist By Id
// @Security    ApiKeyAuth
// @Tags        Checklists
// @Description get update by id
// @ID          get-update-by-id
// @Accept      json
// @Produce     json
// @Param 		id path int true "Checklist ID"
// @Param 		input body domain.UpdateChecklist true "Checklist info"
// @Success     200     {object} statusResponse
// @Failure     400,404 {object} statusResponse
// @Failure     500     {object} statusResponse
// @Failure     default {object} statusResponse
// @Router      /api/lists/{id} [put]
func (h *Handler) updateListById(ctx *gin.Context) {

	userId, err := getUserIdFromContext(ctx)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var input domain.UpdateChecklist
	if err = ctx.BindJSON(&input); err != nil {
		h.newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	input.ID, err = strconv.Atoi(ctx.Param("id"))
	if err != nil || input.ID == 0 {
		h.newErrorResponse(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.serviceChecklist.Update(userId, input); err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	h.GrpcLog.Send(grpcLog.LogRequest_UPDATE, grpcLog.LogRequest_CHECKLIST, int64(input.ID))

	ctx.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary     Delete Checklist By Id
// @Security    ApiKeyAuth
// @Tags        Checklists
// @Description get delete by id
// @ID          get-delete-by-id
// @Accept      json
// @Produce     json
// @Param       id path int     true  "Checklist ID"
// @Success     200     {object} statusResponse
// @Failure     400,404 {object} statusResponse
// @Failure     500     {object} statusResponse
// @Failure     default {object} statusResponse
// @Router      /api/lists/{id} [delete]
func (h *Handler) deleteListById(ctx *gin.Context) {

	userId, err := getUserIdFromContext(ctx)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		h.newErrorResponse(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.serviceChecklist.Delete(userId, id)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	h.GrpcLog.Send(grpcLog.LogRequest_DELETE, grpcLog.LogRequest_CHECKLIST, int64(id))

	ctx.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
