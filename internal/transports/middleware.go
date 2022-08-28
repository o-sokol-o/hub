// Прослойка, которая парсит токен из запроса и предоставляет доступ для наших Энд-Поинтов

package handler_api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/AquaEngineering/AquaHub/pkg/jwt_processing"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
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
func (h *Handler) userIdentity_middleware(c *gin.Context) {

	// Получим хедер авторизации, валидируем его, что он не пустой.
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		h.newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	// Вызовем функцию Split в которой укажем разделить нашу строку по пробелам
	headerParts := strings.Split(header, " ")
	// При корректном хедере эта функция должна вернуть массив длиной в 2 элемента
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		// при ошибках возвращаем Status Code 401 - пользователь не авторизован
		h.newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	// теперь нужно распарсить token
	if len(headerParts[1]) == 0 {
		h.newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	// Метод ParseToken принимает token в качестве аргумента
	// и возвращать id пользователя при успешном парcинге
	userId, err := jwt_processing.ParseToken(headerParts[1])
	if err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Если операция ParseToken успешна - запишем значение id в контекст.
	// Это мы делаем для того чтобы иметь доступ к id пользователям (которые делают запрос)
	// в последующих обработчиках, которые вызываются после данной прослойки.
	c.Set(userCtx, userId)
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
