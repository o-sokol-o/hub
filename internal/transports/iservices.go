package handler_api

import (
	"github.com/AquaEngineering/AquaHub/internal/domain"
)

// Интерфейсы должны объявляться на том уровне абстракции (в том файле),
// где они используются, а не на том где реализуется

// Определение интерфейсов к сущностям бизнес логики

//go:generate mockgen -source=service.go -destination=mocks/mock.go

// Authorization - Авториза́ция — предоставление определённому лицу или группе лиц прав на
// выполнение определённых действий; а также процесс проверки данных прав при попытке
// выполнения этих действий. Часто можно услышать выражение, что какой-то
// человек «авторизован» для выполнения данной операции — это значит, что он имеет на неё право.

// Authentication - Аутентифика́ция — процедура проверки подлинности, например:
// проверка подлинности пользователя путём сравнения введённого им пароля с паролем,
// сохранённым в базе данных пользовательских логинов; подтверждение подлинности
// электронного письма путём проверки цифровой подписи письма по открытому ключу отправителя;

type IServiceAuthentications interface {
	// Принимает структуру User в качестве аргумента
	// Возвращает ID созданного нового пользователя из БД
	CreateUser(user domain.User) (int, error)

	Authenticate(email, password string) (*domain.User, error)

	// --old--

	// GetUser_LoginPassword(username, password string) (domain.User, error)
	// GenerateToken(username, password string) (string, error)
	// ParseToken(token string) (int, error)

	GetUserHWfromTokens(h_token, u_token string) ([]domain.SensorDataSet, error) // user_id, aquahub_id, []{device_id, sensor_id}
	GetAquahubIdfromTokens(h_token, u_token string) (int, error)
}

type IServiceChecklist interface {
	Create(userId int, list domain.UpdateChecklist) (int, error)
	// GetAllChecklist() ([]domain.Checklist, error)
	GetAllChecklistOfUser(userId int) ([]domain.Checklist, error)
	GetById(userId, listId int) (*domain.Checklist, error)
	Delete(userId, listId int) error
	Update(userId int, input domain.UpdateChecklist) error
}

type IServiceChecklistItem interface {
	Create(userId, listId int, item domain.ChecklistItem) (int, error)
	GetAll(userId, listId int) ([]domain.ChecklistItem, error)
	GetById(userId, listId, itemId int) (domain.ChecklistItem, error)
	Delete(userId, listId, itemId int) error
	Update(userId, listId, itemId int, input domain.UpdateChecklistItem) error
}

type IServiceAquahubList interface {

	//-------------------------
	AppendDataOfSensor(list []domain.SensorDataSet) error

	DeviceCreateOrUpdate(aquahub_id, device_local_id int, value string) error
	SensorCreateOrUpdate(aquahub_id, device_local_id, sensor_local_id int, value string) error

	// GetAllAquahubOfUser(userId int) ([]domain.AquahubList, error)
	// GetDevicesOfAquahub(aquahubId int) ([]domain.AquahubList, error)
	// GetSensorsOfDevice(deviceId int) ([]domain.AquahubList, error)
	// GetDataSetOfSensor(sensorId int) ([]domain.SensorDataSet, error)

	// GetNameOfDeviceSensor(sensor_id int) (domain.NameOfDeviceSensor, error)
}
