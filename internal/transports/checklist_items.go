package handler_api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/o-sokol-o/hub/internal/domain"
)

func (h *Handler) createItem(ctx *gin.Context) {

	userId, err := getUserIdFromContext(ctx)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || listId == 0 {
		h.newErrorResponse(ctx, http.StatusBadRequest, "invalid list id param")
		return
	}

	var input domain.ChecklistItem
	if err := ctx.BindJSON(&input); err != nil {
		h.newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.serviceChecklistItem.Create(userId, listId, input)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllItems(ctx *gin.Context) {

	userId, err := getUserIdFromContext(ctx)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || listId == 0 {
		h.newErrorResponse(ctx, http.StatusBadRequest, "invalid list id param")
		return
	}

	items, err := h.serviceChecklistItem.GetAll(userId, listId)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if items == nil {
		items = make([]domain.ChecklistItem, 1)
	}

	ctx.JSON(http.StatusOK, items)
}

func (h *Handler) getItemById(ctx *gin.Context) {

	userId, err := getUserIdFromContext(ctx)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || listId == 0 {
		h.newErrorResponse(ctx, http.StatusBadRequest, "invalid list id param")
		return
	}

	itemId, err := strconv.Atoi(ctx.Param("item_id"))
	if err != nil || itemId == 0 {
		h.newErrorResponse(ctx, http.StatusBadRequest, "invalid item id param")
		return
	}

	item, err := h.serviceChecklistItem.GetById(userId, listId, itemId)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (h *Handler) updateItem(ctx *gin.Context) {

	userId, err := getUserIdFromContext(ctx)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || listId == 0 {
		h.newErrorResponse(ctx, http.StatusBadRequest, "invalid list id param")
		return
	}

	itemId, err := strconv.Atoi(ctx.Param("item_id"))
	if err != nil || itemId == 0 {
		h.newErrorResponse(ctx, http.StatusBadRequest, "invalid item id param")
		return
	}

	var input domain.UpdateChecklistItem
	if err := ctx.BindJSON(&input); err != nil {
		h.newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := input.Validate(); err != nil {
		h.newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.serviceChecklistItem.Update(userId, listId, itemId, input); err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteItem(ctx *gin.Context) {

	userId, err := getUserIdFromContext(ctx)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || listId == 0 {
		h.newErrorResponse(ctx, http.StatusBadRequest, "invalid list id param")
		return
	}

	itemId, err := strconv.Atoi(ctx.Param("item_id"))
	if err != nil || itemId == 0 {
		h.newErrorResponse(ctx, http.StatusBadRequest, "invalid item id param")
		return
	}

	err = h.serviceChecklistItem.Delete(userId, listId, itemId)
	if err != nil {
		h.newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{"ok"})
}
