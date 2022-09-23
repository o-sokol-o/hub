package repository

import (
	"fmt"
	// "strings"

	"github.com/jmoiron/sqlx"
	"github.com/o-sokol-o/hub/internal/domain"
	"github.com/sirupsen/logrus"
)

type AquahubListPostgres struct {
	db  *sqlx.DB
	log *logrus.Logger
}

func NewAquahubListPostgres(log *logrus.Logger, db *sqlx.DB) *AquahubListPostgres {
	return &AquahubListPostgres{log: log, db: db}
}

//__________________________________________________________________________________________________________________________________________________________________

/*
func (r *AquahubListPostgres) Create(userId int, list domain.AquahubList) (int, error) {
	// Для того, чтобы создать списки в базе данных, нам необходимо провести 2 операции вставки:
	//  - сначала в таблицу AquahubLists,
	//  - потом в таблицу usersLists, которая связывает пользователей с их списками.
	// Поэтому данные операции мы будем проводить в транзакции.
	// Транзакция это последовательность нескольких операций,
	// которая рассматривается как одна операция.
	// Транзакция либо выполняется целиком, либо не выполняется вообще.

	// Для создания транзакции в объекте db есть метод Begin.
	// Создадим новую транзакцию, в которой выполним 2 операции вставки.
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// Выполним запрос для создания записи в таблице AquahubLists.
	// Возвращаем id нового списка
	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", aquahubListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Сделаем вставку в таблицу usersLists, в которой свяжем id пользователя и id нового списка.
	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)

	// Для простого выполнения запроса, без чтения возвращаемой информации воспользуемся методом Exec.
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback() // В случае ошибок мы вызываем метод Rollback, которая откатывает все изменения БД до начала выполнения транзакции.
		return 0, err
	}

	// После выполнения транзакции вызовем метод Commit, который применит изменения к БД и закончит транзакцию.
	return id, tx.Commit()
}

//__________________________________________________________________________________________________________________________________________________________________

func (r *AquahubListPostgres) GetAll_Aquahub() ([]domain.AquahubList, error) {
	var lists []domain.AquahubList // Создадим слайс списка

	// Подготовим запрос.
	// tl.id, tl.title, tl.description - возвращаемые поля
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl", aquahubListsTable)

	// На этот раз мы используем для выборки из базы метод селект.
	// Он работает аналогично с методом Get только применяется при выборке больше одного элемента
	// и результат записывает в слайс.
	err := r.db.Select(&lists, query)

	return lists, err
}

//__________________________________________________________________________________________________________________________________________________________________



func (r *AquahubListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2",
		aquahubListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)

	return err
}

//__________________________________________________________________________________________________________________________________________________________________

func (r *AquahubListPostgres) Update(userId, listId int, input domain.UpdateAquahubListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	// title=$1
	// description=$1
	// title=$1, description=$2
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		aquahubListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, listId, userId)

	r.log.Debugf("updateQuery: %s", query)
	r.log.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
*/
//_____________________________________________________________________________________________________

func (r *AquahubListPostgres) GetAquahubs_OfUser(userId int) ([]domain.AquahubList, error) {

	var lists []domain.AquahubList // Создадим слайс списка

	// Подготовим запрос.
	// В нем мы будем делать выборку из базы, используя конструкцию INNER JOIN.
	// Команда INNER JOIN при SELECT помогает выбрать только те записи,
	// которые имеют одинаковое значение в обеих таблицах.
	// В нашем случае нам нужно выбрать все записи из таблицы списков, которые также есть в таблице usersLists,
	// и при этом связаны id пользователя.
	// tl.id, tl.title, tl.description - возвращаемые поля
	query := fmt.Sprintf("SELECT id, title, description FROM %s WHERE user_id = $1", aquahubsTable)

	// На этот раз мы используем для выборки из базы метод селект.
	// Он работает аналогично с методом Get только применяется при выборке больше одного элемента
	// и результат записывает в слайс.
	// Также нам нужно добавить такие дебилы нашу структуру, чтобы иметь возможность делать выборки из базы.
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

//__________________________________________________________________________________________________________________________________________________________________

func (r *AquahubListPostgres) GetDevices_OfAquahub(aquahubId int) ([]domain.AquahubList, error) {
	var list []domain.AquahubList

	query := fmt.Sprintf(`SELECT id, title, description, status FROM %s WHERE aquahub_id = $1 ORDER BY title`, devicesTable)
	err := r.db.Select(&list, query, aquahubId)

	// r.log.Infof("Query: %s,   %d", query, aquahubId)
	// r.log.Info(list)

	return list, err
}

//__________________________________________________________________________________________________________________________________________________________________

func (r *AquahubListPostgres) GetSensors_OfDevice(deviceId int) ([]domain.AquahubList, error) {
	var list []domain.AquahubList

	query := fmt.Sprintf(`SELECT id, title, description FROM %s WHERE device_id = $1 ORDER BY title`, sensorsTable)
	err := r.db.Select(&list, query, deviceId)

	// r.log.Infof("Query: %s,   %d", query, deviceId)
	// r.log.Info(list)

	return list, err
}

//__________________________________________________________________________________________________________________________________________________________________

func (r *AquahubListPostgres) GetDataSet_OfSensor(sensorId int) ([]domain.SensorDataSet, error) {

	var list []domain.SensorDataSet

	// query := fmt.Sprintf(`SELECT id, user_id, aquahub_id, device_id, sensor_id, created_at, user_time, value FROM %s WHERE sensor_id = $1`, sensorDataSetTable)
	// query := `SELECT created_at, value FROM sensor_data WHERE sensor_id = 1`
	query := fmt.Sprintf(`SELECT created_at, Value FROM %s WHERE sensor_id = $1 ORDER BY created_at DESC LIMIT 10`, sensorDataSetTable)
	err := r.db.Select(&list, query, sensorId)

	// r.log.Infof("Query: %s,   %d", query, sensorId)
	// r.log.Info(list)

	return list, err
}

//__________________________________________________________________________________________________________________________________________________________________

func (r *AquahubListPostgres) AppendData_OfSensor(list []domain.SensorDataSet) error {

	// Сделаем вставку в таблицу usersLists, в которой свяжем id пользователя и id нового списка.
	query := fmt.Sprintf(`INSERT INTO %s
	(account_id, aquahub_id, device_id, sensor_id, local_device_id, local_sensor_id, value, created_at)
	VALUES (:account_id, :aquahub_id, :device_id, :sensor_id, :local_device_id, :local_sensor_id, :value, :created_at)`, sensorDataSetTable)

	_, err := r.db.NamedExec(query, list)

	// r.log.Infof("Query: %s,   %d", query, sensorId)
	// r.log.Info(list)

	//	fmt.Printf("\n\nQuery: \n%s\n\n%v\n\n", query, list)

	return err
}

//__________________________________________________________________________________________________________________________________________________________________

func (r *AquahubListPostgres) Device_CreateOrUpdate(aquahub_id int, device_local_id int, value string) error {

	// Команда INNER JOIN при SELECT помогает выбрать только те записи,
	// которые имеют одинаковое значение в обеих таблицах.
	query := fmt.Sprintf(`SELECT dt.id
							FROM  %s dt
							INNER JOIN %s aht on dt.aquahub_id = aht.id
							WHERE aht.id = $1 AND dt.local_id = $2`,
		devicesTable, aquahubsTable)

	// fmt.Printf("\n\n%s\n\n", query)
	args := make([]interface{}, 0)
	var device_id int

	err := r.db.Get(&device_id, query, aquahub_id, device_local_id)
	if err != nil {

		query := fmt.Sprintf(`INSERT INTO %s (aquahub_id, local_id, title, description) VALUES ($1, $2, $3, $4)`, devicesTable)
		fmt.Printf("\n\nQuery: \n%s\n\n", query)

		args = append(args, aquahub_id, device_local_id, value, "Description of the "+value)

		_, err = r.db.Exec(query, args...)

		return err
	}

	// Сделаем вставку в таблицу usersLists, в которой свяжем id пользователя и id нового списка.
	query = fmt.Sprintf(`UPDATE %s SET title = $3  WHERE aquahub_id = $1 AND id = $2`, devicesTable)
	// fmt.Printf("\n\nQuery: \n%s\n\n", query)

	args = append(args, aquahub_id, device_id, value)

	_, err = r.db.Exec(query, args...)

	// r.log.Infof("Query: %s,   %d", query, sensorId)

	// fmt.Printf("\n\nQuery: \n%s\n%v\n", query, err.Error())

	return err
}

//__________________________________________________________________________________________________________________________________________________________________

func (r *AquahubListPostgres) Sensor_CreateOrUpdate(aquahub_id, device_local_id, sensor_local_id int, value string) error {

	args := make([]interface{}, 0)
	var device_id int

	// Команда INNER JOIN при SELECT помогает выбрать только те записи,
	// которые имеют одинаковое значение в обеих таблицах.
	query := fmt.Sprintf(`SELECT dt.id
							FROM  %s dt
							INNER JOIN %s aht on dt.aquahub_id = aht.id
							WHERE aht.id = $1 AND dt.local_id = $2`,
		devicesTable, aquahubsTable)

	// fmt.Printf("\nSensor_CreateOrUpdate SELECT:\n%s\n\n", query)

	err := r.db.Get(&device_id, query, aquahub_id, device_local_id)
	if err == nil {

		query := fmt.Sprintf(`SELECT st.id
		FROM  %s st
		INNER JOIN %s dt on st.device_id = dt.id
		WHERE st.device_id = $1 AND st.local_id = $2`,
			sensorsTable, devicesTable)

		// fmt.Printf("\nSensor_CreateOrUpdate SELECT:\n%s\n\n", query)

		var sensor_id int

		err := r.db.Get(&sensor_id, query, device_id, sensor_local_id)
		if err != nil {

			query := fmt.Sprintf(`INSERT INTO %s (device_id, local_id, title, description) VALUES ($1, $2, $3, $4)`, sensorsTable)
			// fmt.Printf("\nSensor_CreateOrUpdate INSERT:\n%s\n\n", query)

			args = append(args, device_id, sensor_local_id, fmt.Sprintf(`Sensor %d`, sensor_local_id), fmt.Sprintf(`Description of the sensor %d`, sensor_local_id))

			_, err = r.db.Exec(query, args...)

			return err
		}
	}

	return err
}

//__________________________________________________________________________________________________________________________________________________________________

func (r *AquahubListPostgres) GetName_DeviceSensor(sensor_id int) (domain.NameOfDeviceSensor, error) {

	query := fmt.Sprintf(`SELECT d.title AS NameDevice, s.title AS NameSensor FROM %s s 
							INNER JOIN %s d ON ( s.device_id = d.id  )  
							WHERE s.id = $1`,
		sensorsTable, devicesTable)

	var nDevSen domain.NameOfDeviceSensor

	err := r.db.Get(&nDevSen, query, sensor_id)

	// fmt.Printf("\nGetName_DeviceSensor query:\n%s\nNameDevice = %s  NameSensor= %s\n", query, nDevSen.NameDevice, nDevSen.NameSensor)

	return nDevSen, err
}
