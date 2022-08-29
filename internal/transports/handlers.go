package handler_api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"

	_ "github.com/AquaEngineering/AquaHub/docs"
	"github.com/AquaEngineering/AquaHub/internal/domain"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ====================================================================================================

type Handler struct {
	serviceAuthentications IServiceAuthentications
	serviceChecklist       IServiceChecklist
	serviceChecklistItem   IServiceChecklistItem
	serviceAquahubList     IServiceAquahubList

	Router       *gin.Engine
	cache        domain.Cache
	log          *logrus.Logger
	sessionStore *sessions.CookieStore
}

// Внедрение зависимостей:
// Обработчики будут обращаться к Сервисам, поэтому в конструкторе ждём интерфейсы к Сервисам

func NewHandler(log *logrus.Logger, cache domain.Cache, a IServiceAuthentications, b IServiceChecklist, c IServiceChecklistItem, d IServiceAquahubList) *Handler {
	return &Handler{
		log:                    log,
		cache:                  cache,
		serviceAuthentications: a,
		serviceChecklist:       b,
		serviceChecklistItem:   c,
		serviceAquahubList:     d,
	}
}

// Инициализация роутера с end-point_ами
func (h *Handler) InitRoutes(sessionKey string) error {

	h.sessionStore = sessions.NewCookieStore([]byte(sessionKey))

	// Set the router as the default one provided by Gin
	router := gin.Default()

	router.Use(h.middleware_PrintHeader)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Методы авторизации
	auth := router.Group("/auth") // группа маршрутов "/auth"
	{
		auth.POST("/sign-up", h.signUp) // маршрут (end-point) "/auth/sign-up"
		auth.POST("/sign-in", h.signIn) // маршрут (end-point) "/auth/sign-in"
	}

	// Методы работы со списком и итемами
	api := router.Group("/api", h.userIdentity_middleware) // Для группы маршрутов "/api" зададим middleware обработчик
	{
		lists := api.Group("/lists") // группа маршрутов "/api/lists"
		{
			lists.POST("/", h.createList)          // end-point "/api/lists/"
			lists.GET("/", h.getAllLists)          // end-point "/api/lists/"
			lists.GET("/:id", h.getListById)       // end-point "/api/lists/:id", после ":" имя параметра, к которому можно получить доступ ( "id" )
			lists.PUT("/:id", h.updateListById)    // end-point "/api/lists/:id", есть параметр "id"
			lists.DELETE("/:id", h.deleteListById) // end-point "/api/lists/:id", есть параметр "id"

			items := lists.Group(":id/items") // группа маршрутов "/api/lists/:id/items"
			{
				items.POST("/", h.createItem) // end-point "/api/lists/:id/items/"
				items.GET("/", h.getAllItems) // end-point "/api/lists/:id/items/"
				items.GET("/:item_id", h.getItemById)
				items.PUT("/:item_id", h.updateItem)
				items.DELETE("/:item_id", h.deleteItem)
			}
		}
	}

	// Сгруппировать вместе маршруты, связанные с AquaHub API v1
	apiRoutes := router.Group("/v1", h.checkApiKey_middleware)
	{
		apiRoutes.GET("/sensor", h.sensorDataStore)

		apiRoutes.GET("/device/add", h.api_DeviceAdd)
		apiRoutes.GET("/sensor/add", h.api_SensorAdd)

		apiRoutes.GET("/device/meta", h.api_DeviceMeta)
		apiRoutes.GET("/sensor/meta", h.api_SensorMeta)
	}

	h.Router = router
	return nil
}
