package domain

import (
	"errors"
	"time"

	"github.com/lib/pq"
)

// Поля полностью совпадают с БД
// Добавлены JSON теги чтобы корректно принимать и выводить данные в HTTP-запросах
// Тег binding:"required" - валидирует наличие данного поля в теле запроса (является реализацией фремворка gin)

// ID         string          `json:"id" db:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
// Name       string          `json:"name"  validate:"required" example:"Rocket Launch"`
type AquahubList struct {
	ID          int           `json:"id" db:"id" validate:"required"`
	AccountID   int           `json:"account_id" db:"account_id" validate:"required" truss:"api-create"`
	Title       string        `json:"title" db:"title" validate:"required" binding:"required"`
	Description string        `json:"description" db:"description" validate:"required"`
	Status      AquaHubStatus `json:"status" db:"status" validate:"omitempty,oneof=active archived" enums:"active,archived" swaggertype:"string" example:"active"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at" truss:"api-read"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at" truss:"api-read"`
	ArchivedAt  *pq.NullTime  `json:"archived_at,omitempty" db:"archived_at" truss:"api-hide" swaggertype:"string"`
}

// Checklists a list of Checklists.
type AquahubLists []*AquahubList

// ChecklistFindRequest defines the possible options to search for checklists. By default
// archived checklist will be excluded from response.
type AquahubListFindRequest struct {
	Where           string        `json:"where" example:"name = ? and status = ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
}

//=========================================================================================

type NameOfDeviceSensor struct {
	NameDevice string
	NameSensor string
}

type SensorDataSet struct {
	ID         int       `json:"id" db:"id"`
	Account_id int       `json:"account_id" db:"account_id"`
	Aquahub_id int       `json:"aquahub_id" db:"aquahub_id"`
	Device_id  int       `json:"device_id" db:"device_id"`
	Sensor_id  int       `json:"sensor_id" db:"sensor_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	User_time  time.Time `json:"user_time" db:"user_time"`

	// Список связка:
	//    ID_шников пользователя, аквахаба, устройств, сенсоров БД
	//    с локальными ID_шниками устройств и сенсоров
	// user_id & aquahub_id & device_id & sensor_id <===> local_device_id & local_sensor_id

	// По запросу api_v1 приходят данные: Local_device_id, Local_sensor_id, Value
	Local_device_id int `json:"local_device_id" db:"local_device_id"`
	Local_sensor_id int `json:"local_sensor_id" db:"local_sensor_id"`

	Value string `json:"value" db:"value"`
}

//--------------------------------------------------------------------------------------

type AquahubItem struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}

type AquahubListsItem struct {
	Id     int
	ListId int
	ItemId int
}

type UpdateAquahubListInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdateAquahubListInput) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

type UpdateAquahubItemInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (i UpdateAquahubItemInput) Validate() error {
	if i.Title == nil && i.Description == nil && i.Done == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

// AquaHubStatus represents the status of checklist.
type AquaHubStatus string

// AquaHubStatus values define the status field of checklist.
const (
	// AquaHubStatus_Active defines the status of active for checklist.
	AquaHubStatus_Active AquaHubStatus = "active"
	// AquaHubStatus_Archived defines the status of disabled for checklist.
	AquaHubStatus_Archived AquaHubStatus = "archived"
)

// AquaHubStatus_Values provides list of valid AquaHubStatus values.
var AquaHubStatus_Values = []AquaHubStatus{
	AquaHubStatus_Active,
	AquaHubStatus_Archived,
}

// String converts the AquaHubStatus value to a string.
func (s AquaHubStatus) String() string {
	return string(s)
}

// AquaHubStatus_ValuesInterface returns the AquaHubStatus options as a slice interface.
func AquaHubStatus_ValuesInterface() []interface{} {
	var l []interface{}
	for _, v := range AquaHubStatus_Values {
		l = append(l, v.String())
	}
	return l
}

/*
// ====================================================================================================
// ====================================   Интерфейсы   ================================================
// ====================================================================================================

type IAquahubs interface {
	DB_find(ctx *gin.Context, claims auth.Claims, req AquahubListFindRequest) (AquahubLists, error)

	//-----------------------------------

	GetAquahubs_OfUser(userId int) ([]AquahubList, error)
	GetDevices_OfAquahub(aquahubId int) ([]AquahubList, error)
	GetSensors_OfDevice(deviceId int) ([]AquahubList, error)
	GetDataSet_OfSensor(sensorId int) ([]SensorDataSet, error)

	AppendData_OfSensor(list []SensorDataSet) error
	Device_CreateOrUpdate(aquahub_id, device_local_id int, value string) error
	Sensor_CreateOrUpdate(aquahub_id, device_local_id, sensor_local_id int, value string) error

	GetName_DeviceSensor(sensor_id int) (NameOfDeviceSensor, error)
}
*/
