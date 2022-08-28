package service

import (
	"github.com/AquaEngineering/AquaHub/internal/domain"
	"github.com/sirupsen/logrus"
)

// Констрктор сервисов.

// Внедрение зависимостей:
// Сервисы обращаются к БД с помощью интерфейсов, поэтому в конструкторе ждём интерфейсы к репозиторию

func NewServices(log *logrus.Logger, cache domain.Cache,

	a IStoreAuthorization,
	b IStoreChecklist,
	c IStoreChecklistItem,
	d IStoreAquahubs) (

	*logrus.Logger, domain.Cache,

	*AuthService,
	*ChecklistService,
	*ChecklistItemService,
	*AquahubListService) {

	return log, cache,

		NewAuthService(cache, a),
		NewChecklistService(b),
		NewChecklistItemService(c, b),
		NewAquahubListService(d)
}
