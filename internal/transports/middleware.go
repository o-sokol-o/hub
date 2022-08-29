// Прослойка, которая парсит токен из запроса и предоставляет доступ для наших Энд-Поинтов

package handler_api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	userSession         = "userSession"
)

func (h *Handler) middleware_PrintHeader(ctx *gin.Context) {

	fmt.Printf("\n---------------- header begin -------------------\n")
	// fmt.Println(ctx.Request.Header)
	// fmt.Printf("\n---------------- header begin all -------------------\n\n")

	for k, vals := range ctx.Request.Header {
		fmt.Printf("%s\n", k)
		for _, v := range vals {
			fmt.Printf("\t%s\n", v)
		}
	}
	fmt.Printf("\n---------------- header end -------------------\n\n")
}

// Нам нужно получить значение из хедера авторизации, валидировать его,
// парсить токен и записать пользователя в контекст.
func (h *Handler) userIdentity_middleware(ctx *gin.Context) {

	session, _ := h.sessionStore.Get(ctx.Request, "session")

	// Check if user is authenticated
	userId, ok := session.Values["userId"].(int)
	if !ok || userId == 0 {
		h.newErrorResponse(ctx, http.StatusUnauthorized, "error: no userId into session")
		return
	}

	// Если операция ParseToken успешна - запишем значение id в контекст.
	// Это мы делаем для того чтобы иметь доступ к id пользователям (которые делают запрос)
	// в последующих обработчиках, которые вызываются после данной прослойки.
	ctx.Set(userCtx, userId)
}

// Функция, достающая ID пользователя из контекста, обрабатывает ошибки и выводит response.
func getUserIdFromContext(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok || idInt == 0 {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
