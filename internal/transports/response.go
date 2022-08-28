package handler_api // Работа с ошибками

import (
	"github.com/gin-gonic/gin"
)

type idResponse struct {
	ID int `json:"id"  example:"1"`
}

type statusResponse struct {
	Status string `json:"status"  example:"Ok"`
}

// Стандартная функция обработки ошибок принимающих обработчиков
func (h *Handler) newErrorResponse(ctx *gin.Context, statusCode int, message string) {
	// Здесь мы будем выводить сообщения ошибок в консоль
	h.log.Error(message)

	// А здесь вызываем метод нашего контекста AbortWithStatusJSON
	// который принимает http Status Code и тело ответа.
	// В качестве ответа функция принимает интерфейс, поэтому можем передать как структуру,
	// так и мапу со строкой в качестве ключа и интерфейс в качестве значения.
	// Поскольку у одного end-point_а (маршрута) может быть несколько последовательных обработчиков,
	// AbortWithStatusJSON блокирует выполнение последующих обработчиков и формирует ответ клиенту
	// в виде statuscode и тело сообщения в формате JSON
	ctx.AbortWithStatusJSON(statusCode, statusResponse{message})
}
