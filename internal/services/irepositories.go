package service

import (
	"github.com/o-sokol-o/hub/internal/domain"
)

//go:generate mockgen -source=irepositories.go -destination=mocks/mock.go

// Authorization - Авториза́ция — предоставление определённому лицу или группе лиц прав на
// выполнение определённых действий; а также процесс проверки данных прав при попытке
// выполнения этих действий. Часто можно услышать выражение, что какой-то
// человек «авторизован» для выполнения данной операции — это значит, что он имеет на неё право.

// Authentication - Аутентифика́ция — процедура проверки подлинности, например:
// проверка подлинности пользователя путём сравнения введённого им пароля с паролем,
// сохранённым в базе данных пользовательских логинов; подтверждение подлинности
// электронного письма путём проверки цифровой подписи письма по открытому ключу отправителя;

// Интерфейсы должны объявляться на том уровне,
// где они используются, а не на том где реализуется
// Определение интерфейсов к сущностям БД (идентичны сущночстям бизнес логики)

type IStoreAuthorization interface {
	GetUserByEmail(email string) (*domain.User, error)

	// Принимает структуру User в качестве аргумента
	// Возвращает ID созданного нового пользователя из БД
	CreateUser(user domain.User) (int, error)

	GetUser_LoginPassword(username, password string) (domain.User, error)

	GetUserHW_fromTokens(h_token, u_token string) ([]domain.SensorDataSet, error) // user_id, aquahub_id, []{device_id, sensor_id}
	GetAquahubId_fromTokens(h_token, u_token string) (int, error)
}

type IStoreChecklist interface {
	Create(userId int, list domain.CreateChecklist) (int, error)
	// GetAll_Checklist() ([]domain.Checklist, error)
	GetAll_ChecklistOfUser(userId int) ([]domain.Checklist, error)
	GetById(userId, listId int) (*domain.Checklist, error)
	Delete(userId, listId int) error
	Update(userId int, input domain.UpdateChecklist) error
}

type IStoreChecklistItem interface {
	Create(listId int, item domain.ChecklistItem) (int, error)
	GetAll(userId, listId int) ([]domain.ChecklistItem, error)
	GetById(userId, listId, itemId int) (domain.ChecklistItem, error)
	Delete(userId, listId, itemId int) error
	Update(userId, listId, itemId int, input domain.UpdateChecklistItem) error
}

type IStoreAquahubs interface {

	//-----------------------------------

	GetAquahubs_OfUser(userId int) ([]domain.AquahubList, error)
	GetDevices_OfAquahub(aquahubId int) ([]domain.AquahubList, error)
	GetSensors_OfDevice(deviceId int) ([]domain.AquahubList, error)
	GetDataSet_OfSensor(sensorId int) ([]domain.SensorDataSet, error)

	AppendData_OfSensor(list []domain.SensorDataSet) error
	Device_CreateOrUpdate(aquahub_id, device_local_id int, value string) error
	Sensor_CreateOrUpdate(aquahub_id, device_local_id, sensor_local_id int, value string) error

	GetName_DeviceSensor(sensor_id int) (domain.NameOfDeviceSensor, error)
}
