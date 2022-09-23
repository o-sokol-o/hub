// handlers.user.go

package handler_api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/o-sokol-o/hub/internal/domain"

	"github.com/gin-gonic/gin"
)

type sensor_data_in map[int]string

type devices_sensors_data_in map[int]sensor_data_in

type RequestInput_API_v1 struct {
	ID              int
	DevData         devices_sensors_data_in
	API_Key         string
	Hard_token      string
	User_token      string
	Request_Payload string
}

func append_devices_sensors_data(m devices_sensors_data_in, devices_id int, sensors_id int, data string) {

	if m == nil {
		fmt.Printf("error append_devices_sensors_data: No m!!!\n")
		return
	}

	if elem, ok := m[devices_id]; ok {
		elem[sensors_id] = data
	} else {
		sd := make(sensor_data_in)
		sd[sensors_id] = data
		m[devices_id] = sd
	}

	// fmt.Printf("append_devices_sensors_data: %v   %v   %v\n", m, m[devices_id], m[devices_id][sensors_id])
}

func parce_devices_sensors_id(s string) (int, int, bool) {

	if before, after, found := strings.Cut(s, "f"); before == "" || found {

		// fmt.Printf("Cut(%q, %q) = %q, %q, %v\n", s, "f", before, after, found)

		// И обернем эту функцию в Atoi, из стандартной библиотеки strconv, для преобразования строки в число.
		dId, err := strconv.Atoi(after)
		if err != nil {
			return 0, 0, false
		}

		sId := dId % 100
		dId = dId / 100

		return dId, sId, true
	}

	return 0, 0, false
}

//##############################################################################################################################

func parce_SensorsDataset(array_pair []string, DevData devices_sensors_data_in) error {

	for _, pair := range array_pair {

		if key := strings.Split(pair, "="); len(key) == 2 && len(key[0]) > 1 && len(key[1]) > 0 {

			// fmt.Printf("len=%d cap=%d %v\n", len(key), cap(key), key)

			if devices_id, sensors_id, err := parce_devices_sensors_id(key[0]); err {

				// fmt.Printf("parce_devices_sensors_id ( %s ) = %d, %d\n", key[0], devices_id, sensors_id)

				append_devices_sensors_data(DevData, devices_id, sensors_id, key[1])
				// fmt.Printf("devices_sensors_data: %v\n", input.DevData)
			} else {

				fmt.Printf("error parceRoute_SensorsDataset: parce_devices_sensors_id ( %s ) = %d, %d", key[0], devices_id, sensors_id)
				return errors.New("error parceRoute_SensorsDataset: parce_devices_sensors_id")
			}

		} else {

			fmt.Printf("error parceRoute_SensorsDataset: Separate with '=' error, pair = %v", key)
			return errors.New("error parceRoute_SensorsDataset: Separate with '=' error")
		}
	}

	return nil
}

//_____________________________________________________________________________________________________________________________________

func parce_DeviceAdd(array_pair []string, DevData devices_sensors_data_in) error {

	devices_id := 1

	for _, pair := range array_pair {

		if key := strings.Split(pair, "="); len(key) == 2 && len(key[0]) > 1 && len(key[1]) > 0 {

			fmt.Printf("parce_DeviceAdd: len=%d cap=%d %v\n", len(key), cap(key), key)

			if key[0] == "id" {

				dId, err := strconv.Atoi(key[1])
				if err != nil {
					return errors.New("error parce_DeviceAdd: devices_id not integer")
				}
				devices_id = dId
				fmt.Printf("parce_DeviceAdd: Detected devices_id = %d\n", devices_id)

			} else {

				fmt.Printf("parce_DeviceAdd: Append data: devices_id = %d  sensors_id = %d  data = %s\n", devices_id, 1, key[1])
				append_devices_sensors_data(DevData, devices_id, 1, key[1])
			}

		} else {

			fmt.Printf("error parce_DeviceAdd: Separate with '=' error, pair = %v", key)
			return errors.New("error parce_DeviceAdd: Separate with '=' error")
		}
	}

	return nil
}

//_____________________________________________________________________________________________________________________________________

func parce_DeviceMeta(array_pair []string, DevData devices_sensors_data_in) error {

	devices_id := 1

	for _, pair := range array_pair {

		if key := strings.Split(pair, "="); len(key) == 2 && len(key[0]) > 1 && len(key[1]) > 0 {

			fmt.Printf("parce_Meta: len=%d cap=%d %v\n", len(key), cap(key), key)

			if key[0] == "id" || key[0] == "device_id" {

				dId, err := strconv.Atoi(key[1])
				if err != nil {
					return errors.New("error parce_Meta: devices_id not integer")
				}
				devices_id = dId
				fmt.Printf("parce_Meta: Detected devices_id = %d\n", devices_id)

			}

			if key[0] == "devdata" {

				if list := strings.Split(key[1], ";"); len(list) > 0 {

					fmt.Printf("parce_Meta: Meta devdata split list: %v:\n", list)

					for _, item := range list {

						if before, after, found := strings.Cut(item, ":"); found && len(before) > 0 && len(after) > 0 {

							fmt.Printf("parce_Meta:  before: %s   after: %s\n\n", before, after)
							meta_type, err := strconv.Atoi(before)
							if err != nil {
								return errors.New("error parce_Meta: meta type not integer")
							}

							fmt.Printf("parce_Meta: Append data: devices_id = %d  sensors_id = %d  data = %s\n", devices_id, meta_type, after)
							append_devices_sensors_data(DevData, devices_id, meta_type, after)

						} else {
							fmt.Printf("parce_Meta:   meta data empty item:  before: %s   after: %s\n\n", before, after)
							// return errors.New("error parce_Meta: meta data empty item")
						}
					}
				} else {
					return errors.New("error parce_Meta: meta data apsent")
				}
			}

		} else {

			fmt.Printf("error parce_Meta: Separate with '=' error, pair = %v", key)
			return errors.New("error parce_Meta: Separate with '=' error")
		}
	}

	return nil
}

//_____________________________________________________________________________________________________________________________________

func parce_SensorAdd(array_pair []string, DevData devices_sensors_data_in) error {

	devices_id := 1

	for _, pair := range array_pair {

		if key := strings.Split(pair, "="); len(key) == 2 && len(key[0]) > 1 && len(key[1]) > 0 {

			fmt.Printf("parce_SensorAdd: len=%d cap=%d %v\n", len(key), cap(key), key)

			if key[0] == "device_id" {

				dId, err := strconv.Atoi(key[1])
				if err != nil {
					return errors.New("error parce_SensorAdd: devices_id not integer")
				}
				devices_id = dId
				fmt.Printf("parce_SensorAdd: Detected devices_id = %d\n", devices_id)

			} else {

				fmt.Printf("parce_SensorAdd: Append IDs: devices_id = %d  sensors_ids = %s\n", devices_id, key[1])

				if list := strings.Split(key[1], ","); len(list) > 0 {

					fmt.Printf("parce_SensorAdd: sensor IDs split list: %v:\n", list)

					for _, item := range list {

						if before, after, found := strings.Cut(item, "f"); found && len(before) == 0 && len(after) > 0 {

							fmt.Printf("parce_SensorAdd:  before: %s   after: %s\n\n", before, after)
							sensors_id, err := strconv.Atoi(after)
							if err != nil {
								return errors.New("error parce_SensorAdd: sensor id not integer")
							}
							sensors_id = sensors_id % 100

							fmt.Printf("parce_SensorAdd: Append data: devices_id = %d  sensors_id = %d\n", devices_id, sensors_id)
							append_devices_sensors_data(DevData, devices_id, sensors_id, "")

						} else {
							fmt.Printf("parce_SensorAdd:   sensor id item is empty:  before: %s   after: %s\n\n", before, after)
							// return errors.New("error parce_SensorAdd: sensor id item is empty")
						}
					}
				} else {
					return errors.New("error parce_SensorAdd: sensor IDs apsent")
				}
			}

		} else {

			fmt.Printf("error parce_SensorAdd: Separate with '=' error, pair = %v", key)
			return errors.New("error parce_SensorAdd: Separate with '=' error")
		}
	}

	return nil
}

//_____________________________________________________________________________________________________________________________________

func parce_SensorMeta(array_pair []string, DevData devices_sensors_data_in) error {

	devices_id := 1

	for _, pair := range array_pair {

		if key := strings.Split(pair, "="); len(key) == 2 && len(key[0]) > 1 && len(key[1]) > 0 {

			fmt.Printf("parce_SensorMeta: len=%d cap=%d %v\n", len(key), cap(key), key)

			if key[0] == "device_id" {

				dId, err := strconv.Atoi(key[1])
				if err != nil {
					return errors.New("error parce_SensorMeta: devices_id not integer")
				}
				devices_id = dId
				fmt.Printf("parce_SensorMeta: Detected devices_id = %d\n", devices_id)

			} else {

				before, after, found := strings.Cut(key[0], "f")
				if !found || len(before) > 0 || len(after) == 0 {
					return errors.New("error parce_SensorMeta: sensor_id bad")
				}

				sensor_id, err := strconv.Atoi(after)
				if err != nil {
					return errors.New("error parce_SensorMeta: sensor_id not integer")
				}
				sensor_id = sensor_id % 100

				fmt.Printf("parce_SensorMeta: Append data: devices_id = %d  sensors_id = %d  data = %s\n", devices_id, sensor_id, key[1])
				append_devices_sensors_data(DevData, devices_id, sensor_id, key[1])
				/*
					if list := strings.Split(key[1], ";"); len(list) > 0 {

						fmt.Printf("parce_SensorMeta: Meta devdata split list: %v:\n", list)

						for _, item := range list {

							if before, after, found := strings.Cut(item, ":"); found && len(before) > 0 && len(after) > 0 {

								fmt.Printf("parce_SensorMeta:  before: %s   after: %s\n\n", before, after)
								sen_id, err := strconv.Atoi(before)
								if err != nil {
									return errors.New("error parce_SensorMeta: meta type not integer")
								}

								fmt.Printf("parce_SensorMeta: Append data: devices_id = %d  sensors_id = %d  data = %s\n", devices_id, sensor_id, after)
								append_devices_sensors_data(DevData, devices_id, sensor_id, after)

							} else {
								fmt.Printf("parce_SensorMeta:   meta data empty item:  before: %s   after: %s\n\n", before, after)
								// return errors.New("error parce_SensorMeta: meta data empty item")
							}

						}
					} else {
						return errors.New("error parce_SensorMeta: meta data apsent")
					}
				*/
			}

		} else {

			fmt.Printf("error parce_SensorMeta: Separate with '=' error, pair = %v", key)
			return errors.New("error parce_SensorMeta: Separate with '=' error")
		}
	}

	return nil
}

//_____________________________________________________________________________________________________________________________________

func parce_Route(route string, array_pair []string, DevData devices_sensors_data_in) error {

	switch route {
	case "/v1/device/add":
		return parce_DeviceAdd(array_pair, DevData)

	case "/v1/device/meta":
		return parce_DeviceMeta(array_pair, DevData)

	case "/v1/sensor/add":
		return parce_SensorAdd(array_pair, DevData)

	case "/v1/sensor/meta":
		return parce_SensorMeta(array_pair, DevData)

	case "/v1/sensor":
		return parce_SensorsDataset(array_pair, DevData)

		// default:
		// 	return errors.New("error parce_Route: No route")
	}

	return errors.New("error parce_Route: No route")
}

// Это промежуточное программное обеспечение устанавливает, имеется ли ключ пользователя.
// Это middleware гарантирует, что запрос будет прерван с ошибкой,
// если отсутствует ключ пользователя
func (h *Handler) checkApiKey_middleware(c *gin.Context) {

	c.Set("api_v1_valid", false)

	// ======================= Выделяем и верифицируем формат ключей =======================

	// Отрезали левую часть URI: /v1/device/add   /v1/device/meta   /v1/sensor/add   /v1/sensor/meta   /v1/sensor
	// fmt.Printf("c.Request.RequestURI: %s\n", c.Request.RequestURI)
	before, after, found := strings.Cut(c.Request.RequestURI, "?api_key=")
	// fmt.Printf("before: %s   after: %s\n\n", before, after)
	route := before
	if !found || len(after) < 40 {

		h.log.Printf("error checkApiKey: URI not first cut: before=%s   after=%s   len(after)=%d", before, after, len(after))
		// Добавим ответ при ошибке: даём ложный успех
		c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
		return
	}

	// Отрезали слева часть с ключами
	// fmt.Printf("URI str: %s\n", after)
	before, after, found = strings.Cut(after, "&")
	// fmt.Printf("Parsing api_key: before: %s   after: %s\n\n", before, after)
	if !found || len(before) < 36 || len(after) < 2 {

		h.log.Printf("error checkApiKey: api_key or payload small: before: %s   after: %s   len(after)=%d", before, after, len(after))
		// Добавим ответ при ошибке: даём ложный успех
		c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
		return
	}
	request_payload := after
	api_key := before

	// Разделили api_key на Hard_token и User_token
	// fmt.Printf("api_key = %s\n", api_key)
	before, after, found = strings.Cut(api_key, ":")
	// fmt.Printf("before=%s & after=%s\n\n", before, after)
	if !found || len(before) < 12 || len(after) != 16 {

		h.log.Printf("error checkApiKey: Hard_token or User_token size error: Hard_token=%s   User_token=%s   len(after)=%d", before, after, len(after))
		// Добавим ответ при ошибке: даём ложный успех
		c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
		return
	}

	var input RequestInput_API_v1
	input.API_Key = api_key
	input.Hard_token = before
	input.User_token = after
	input.Request_Payload = request_payload

	// fmt.Printf("API_Key:   %s\n", input.API_Key)
	// fmt.Printf("Hard_token:  %s\n", input.Hard_token)
	// fmt.Printf("User_token:  %s\n\n", input.User_token)

	// ======================= Разделить Request_Payload на пары с сепаратором "&" =======================

	array_pair := strings.Split(request_payload, "&")
	if len(array_pair) < 1 {

		h.log.Printf("error checkApiKey: array_pair < 1: %v", array_pair)
		// Добавим ответ при ошибке: даём ложный успех
		c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
		return
	}

	fmt.Printf("route = %s\n", route)
	// fmt.Printf("Start parsing array_pair: %v\n\n", array_pair)

	input.DevData = make(devices_sensors_data_in)

	// Парсим пары согласно route

	if err := parce_Route(route, array_pair, input.DevData); err != nil || len(input.DevData) < 1 {

		fmt.Printf("error checkApiKey: parce array_pair error: %s\n", err.Error())

		// Добавим ответ при ошибке: даём ложный успех
		c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
		return
	}

	fmt.Printf("Array data: %v\n", input.DevData)

	// for dId, sd := range input.DevData {
	// 	for sId, val := range sd {
	// 		fmt.Printf("devices_id: %d, sensors_id: %d  = %s\n", dId, sId, val)
	// 	}
	// }

	// Если операция ParseToken успешна - запишем значение в контекст.
	// Это мы делаем для того чтобы иметь доступ к данным запроса
	// в последующих обработчиках, которые вызываются после данной прослойки.
	c.Set("api_v1_valid", true)
	c.Set("data_api_v1", input)
}

//____________________________________________________________________________________________________________________

// Функция, достающая ID пользователя из контекста, обрабатывает ошибки и выводит response.
func getRequestInput(c *gin.Context) (*RequestInput_API_v1, error) {

	if a, ok := c.Get("api_v1_valid"); !ok || !a.(bool) {
		return nil, errors.New("data_api_v1 is of parsing invalid")
	}

	a, ok := c.Get("data_api_v1")
	if !ok {
		return nil, errors.New("data_api_v1 id not found")
	}

	idInt, ok := a.(RequestInput_API_v1)
	if !ok {
		return nil, errors.New("data_api_v1 id is of invalid type")
	}

	return &idInt, nil
}

//____________________________________________________________________________________________________________________

// Security    ApiKeyAuth
// Accept      json

// @Summary     AquaHub sensors data store
// @Tags        AquaHub
// @ID          aquahub-sensor-data-store
// @Description ?api_key=aqen104Ur2zNX1Ykwv4:a39831d103eb4c0d &f100=0.01&f101=28&f102=0&f103=17.51&f104=15.52&f105=1072 &f200=17.52&f201=134.06&f202=317&f203=25.7000 &f400=3.27&f401=0.39&f402=3.26&f403=0.39&f404=0.08&f405=0.00 &f4002=0.08 &f11000=504&f10001=24.31&f10004=0
// @Produce     json
// @Success     200     {object}  statusResponse
// @Failure     default {object}  statusResponse
// @Router      /v1/sensor [get]
func (h *Handler) sensorDataStore(c *gin.Context) {

	input, err := getRequestInput(c)
	if err != nil {
		return
	}

	/*
		// Проверяем, доступен ли ключ пользователя
		if !isUserKeyAvailable(input.User_Key) {

			fmt.Printf("\n\n\nError isUserKeyAvailable = %v\n\n\n", input.User_Key)

			// При ошибке даём ложный успех
			// c.JSON(http.StatusOK, a)
			render(c, gin.H{
				"title": "Successful Login"}, "login-successful.html")
			return
		}
	*/

	fmt.Printf("\nAPI_Key:   %s\n", input.API_Key)
	fmt.Printf("Hard_Key:  %s\n", input.Hard_token)
	fmt.Printf("User_Key:  %s\n\n", input.User_token)

	fmt.Printf("Array data: Len = %d\n%v\n", len(input.DevData), input.DevData)
	for dId, sd := range input.DevData {
		fmt.Printf("devices_id: %d, sensors = %v\n", dId, sd)
	}

	// GetUser_Token
	// Метод ParseToken принимает token в качестве аргумента
	// и возвращать userId пользователя при успешном парcинге
	// userId, err := h.services.Authorization.GetUserFromToken(input.Hard_Key, input.User_Key)

	// Получить по токенам список связку:
	//    ID_шников пользователя, аквахаба, устройств, сенсоров БД
	//    с локальными ID_шниками устройств и сенсоров
	// user_id & aquahub_id & device_id & sensor_id <===> local_device_id & local_sensor_id
	list_IDs, err := h.serviceAuthentications.GetUserHWfromTokens(input.Hard_token, input.User_token)

	if err != nil || len(list_IDs) < 1 {

		h.log.Infof("sensorDataStore services.Authorization.GetUserHWfromTokens: %v  Error = %s", list_IDs, err.Error())

		// При ошибке даём ложный успех
		c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
		return
	}
	// fmt.Printf("sensorDataStore.GetUserHWfromTokens:\nList = %v\n\n", list_IDs)

	// Проверка req.DevData [devices_id] [sensors_id] значений на предмет наличия в БД соответствующих device и sensor
	// Если их нет - создать по образцу относительно devices_id
	// (devices_id => прямое соответствие конкретному типу устройства)

	// var input domain.ChecklistList
	// input.Title = req.API_Key
	// for key, val := range req.DevData {
	// 	a := input.Description + fmt.Sprintf("Key: %s, Value: %s     ", key, val)
	// 	if len(a) > 255 {
	// 		break
	// 	}
	// 	input.Description = a
	// }
	// h.log.Infof("Len input.Description: %d", len(input.Description))

	// for dId, sd := range input.DevData {
	// 	for sId, val := range sd {
	// 		fmt.Printf("devices_id: %d, sensors_id: %d  = %s\n", dId, sId, val)
	// 	}
	// }

	// Вызовем метод создания списка,
	// в который передадим наши данные, полученные из токена аутентификации и тела запроса.
	// id, err := h.services.ChecklistList.Create(userId, input)

	var list_append []domain.SensorDataSet

	for dev_id, sd := range input.DevData {
		for sens_id, val := range sd {
			// fmt.Printf("devices_id[ %d ], sensors_id[ %d ] = %s\n", dev_id, sens_id, val)

			// Подготовка списока к вставке со связкой ID_шников:
			//   пользователя, аквахаба, устройств, сенсоров БД
			//   с локальными ID_шниками устройств и сенсоров
			// user_id & aquahub_id & device_id & sensor_id <===> local_device_id & local_sensor_id
			var x domain.SensorDataSet
			x.Value = val
			x.Account_id = 1
			x.Aquahub_id = 1
			x.Device_id = 1
			x.Sensor_id = 1
			x.Local_device_id = dev_id
			x.Local_sensor_id = sens_id
			x.CreatedAt = time.Now().UTC()

			// ищем совпадение пришедших локальных IDшек с записями в БД
			for _, IDs := range list_IDs {
				if IDs.Local_device_id == dev_id && IDs.Local_sensor_id == sens_id {
					x.Account_id = IDs.Account_id
					x.Aquahub_id = IDs.Aquahub_id
					x.Device_id = IDs.Device_id
					x.Sensor_id = IDs.Sensor_id
					break
				}
			}

			// Вставить в результирующий слайс
			list_append = append(list_append, x)

		}
	}
	// fmt.Printf("sensorDataStore:\nlist_append = %v\n\n", list_append)

	// Вставить результирующий слайс в БД
	if err := h.serviceAquahubList.AppendDataOfSensor(list_append); err != nil {
		h.log.Infof("sensorDataStore AquahubList.AppendDataOfSensor: %s", err.Error())
		// При ошибке даём ложный успех
		c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
}

//#################################################################################################################################

func (h *Handler) api_DeviceAdd(c *gin.Context) {

	input, err := getRequestInput(c)
	if err != nil {
		return
	}

	// fmt.Printf("\nAPI_Key:   %s\n", input.API_Key)
	// fmt.Printf("Hard_Key:  %s\n", input.Hard_token)
	// fmt.Printf("User_Key:  %s\n\n", input.User_token)

	// Авторизуемся и получаем aquahub_id
	aquahub_id, err := h.serviceAuthentications.GetAquahubIdfromTokens(input.Hard_token, input.User_token)
	if err != nil || aquahub_id < 1 {
		h.log.Infof("api_DeviceAdd GetAquahubIdfromTokens: AquahubId = %d  Error = %s", aquahub_id, err.Error())
		// При ошибке даём ложный успех
		c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
		return
	}

	// fmt.Printf("api_DeviceAdd: AquahubId = %d input.DevData = %v\n\n", aquahub_id, input.DevData)

	// ищем совпадение пришедших локальных IDшек с записями в БД
	for device_local_id, sd := range input.DevData {
		for _, value := range sd {
			// fmt.Printf("api_DeviceAdd: aquahub_id = %d  dev_id = %d value = %s\n\n", aquahub_id, device_local_id, value)

			// Вставить в БД
			if err := h.serviceAquahubList.DeviceCreateOrUpdate(aquahub_id, device_local_id, value); err != nil {
				h.log.Infof("api_DeviceAdd: DeviceCreateOrUpdate: Error = %s", err.Error())
				// При ошибке даём ложный успех
				c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
				return
			}

		}
	}

	c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
}

//____________________________________________________________________________________________________________________

func (h *Handler) api_DeviceMeta(c *gin.Context) {

	input, err := getRequestInput(c)
	if err != nil {
		return
	}

	fmt.Printf("\nAPI_Key:   %s\n", input.API_Key)
	fmt.Printf("Hard_Key:  %s\n", input.Hard_token)
	fmt.Printf("User_Key:  %s\n\n", input.User_token)

	c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
}

//____________________________________________________________________________________________________________________

func (h *Handler) api_SensorAdd(c *gin.Context) {

	input, err := getRequestInput(c)
	if err != nil {
		return
	}

	// fmt.Printf("\nAPI_Key:   %s\n", input.API_Key)
	// fmt.Printf("Hard_Key:  %s\n", input.Hard_token)
	// fmt.Printf("User_Key:  %s\n\n", input.User_token)

	// Авторизуемся и получаем aquahub_id
	aquahub_id, err := h.serviceAuthentications.GetAquahubIdfromTokens(input.Hard_token, input.User_token)
	if err != nil || aquahub_id < 1 {
		h.log.Infof("api_SensorAdd GetAquahubIdfromTokens: AquahubId = %d  Error = %s", aquahub_id, err.Error())
		// При ошибке даём ложный успех
		c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
		return
	}

	fmt.Printf("api_SensorAdd: AquahubId = %d input.DevData = %v\n\n", aquahub_id, input.DevData)

	// ищем совпадение пришедших локальных IDшек с записями в БД
	for device_local_id, sd := range input.DevData {
		for sensor_local_id, value := range sd {
			// fmt.Printf("api_SensorAdd: aquahub_id = %d  device_local_id = %d sensor_local_id = %d\n\n", aquahub_id, device_local_id, sensor_local_id)

			// Вставить в БД
			if err := h.serviceAquahubList.SensorCreateOrUpdate(aquahub_id, device_local_id, sensor_local_id, value); err != nil {
				h.log.Infof("api_SensorAdd: SensorCreateOrUpdate: Error = %s", err.Error())
				// При ошибке даём ложный успех
				c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
				return
			}

		}
	}

	c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
}

//____________________________________________________________________________________________________________________

func (h *Handler) api_SensorMeta(c *gin.Context) {

	input, err := getRequestInput(c)
	if err != nil {
		return
	}

	fmt.Printf("\nAPI_Key:   %s\n", input.API_Key)
	fmt.Printf("Hard_Key:  %s\n", input.Hard_token)
	fmt.Printf("User_Key:  %s\n\n", input.User_token)

	c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
}

//____________________________________________________________________________________________________________________
