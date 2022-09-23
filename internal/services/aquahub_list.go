package service

import (
	"github.com/o-sokol-o/hub/internal/domain"
)

// Сервис для работы со списками

type AquahubListService struct {
	repo IStoreAquahubs
}

//-------------------------------------------------------------------------

// Также в нашем сервисе и понадобится репозиторий
// Добавим его в качестве поля нашей структуры и будем передавать в конструкторе.
func NewAquahubListService(repo IStoreAquahubs) *AquahubListService {
	return &AquahubListService{repo: repo}
}

/*
// При создании списка мы будем передавать данные на следующий уровень - в репозиторий,
// поэтому в сервисе мы лишь будем возвращать аналогичный метод репозитория.
// Дополнительной логики мы реализовывать не будем.
func (s *ChecklistListService) Create(userId int, list domain.ChecklistList) (int, error) {
	return s.repo.Create(userId, list)
}

// Метод GetAll, который будет возвращать слайс списка вместе с ошибкой.
func (s *ChecklistListService) GetAllChecklist() ([]domain.ChecklistList, error) {
	// В сервисе мы будем вызывать аналогичный метод репозитория, поскольку дополнительной бизнес логики тут нет.
	return s.repo.GetAll_Checklist()
}



func (s *ChecklistListService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}

func (s *ChecklistListService) Update(userId, listId int, input domain.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.repo.Update(userId, listId, input)
}
*/

// Метод GetAll, который будет принимать id пользователя
// и возвращать слайс списка вместе с ошибкой.
func (s *AquahubListService) GetAllAquahubOfUser(userId int) ([]domain.AquahubList, error) {
	// В сервисе мы будем вызывать аналогичный метод репозитория, поскольку дополнительной бизнес логики тут нет.
	return s.repo.GetAquahubs_OfUser(userId)
}

func (s *AquahubListService) GetDevicesOfAquahub(aquahubId int) ([]domain.AquahubList, error) {
	return s.repo.GetDevices_OfAquahub(aquahubId)
}

func (s *AquahubListService) GetSensorsOfDevice(deviceId int) ([]domain.AquahubList, error) {
	return s.repo.GetSensors_OfDevice(deviceId)
}

func (s *AquahubListService) GetDataSetOfSensor(sensorId int) ([]domain.SensorDataSet, error) {
	return s.repo.GetDataSet_OfSensor(sensorId)
}

func (s *AquahubListService) AppendDataOfSensor(list []domain.SensorDataSet) error {
	return s.repo.AppendData_OfSensor(list)
}

func (s *AquahubListService) DeviceCreateOrUpdate(aquahub_id int, device_local_id int, value string) error {
	return s.repo.Device_CreateOrUpdate(aquahub_id, device_local_id, value)
}

func (s *AquahubListService) SensorCreateOrUpdate(aquahub_id, device_local_id, sensor_local_id int, value string) error {
	return s.repo.Sensor_CreateOrUpdate(aquahub_id, device_local_id, sensor_local_id, value)
}

func (s *AquahubListService) GetNameOfDeviceSensor(sensor_id int) (domain.NameOfDeviceSensor, error) {
	return s.repo.GetName_DeviceSensor(sensor_id)
}
