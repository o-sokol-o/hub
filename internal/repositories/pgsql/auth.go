package repository

import (
	"database/sql"
	"fmt"
	"time"

	// account "github.com/o-sokol-o/hub/int/account"
	// "github.com/o-sokol-o/hub/int/user_account"
	"github.com/o-sokol-o/hub/internal/domain"
	"github.com/o-sokol-o/hub/pkg/randomstring"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// Принимает структуру User в качестве аргумента
// Возвращает ID созданного нового пользователя из БД
func (r *AuthPostgres) CreateUser(user domain.User) (int, error) {

	// Для создания транзакции в объекте db есть метод Begin.
	// Создадим новую транзакцию, в которой выполним операции вставки.
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// Составляем запрос вставки в таблицу usersTable
	// $1, $2, $3 - плейсхолдеры в которые подставляются значения при вызове tx.QueryRow
	// query := fmt.Sprintf("INSERT INTO %s (name, login, password_hash, u_token) values ($1, $2, $3, $4) RETURNING id", accountsTable)
	query_user := fmt.Sprintf("INSERT INTO %s (first_name, last_name, email, password_hash, password_salt) values ($1, $2, $3, $4, $5) RETURNING id", usersTable)

	// Выполнить запрос
	// Объект row хранит в себе информацию возвращаемую из БД (строка)
	// row := tx.QueryRow(query, user.Name, user.Login, user.Password, user.Token)
	row := tx.QueryRow(query_user, user.FirstName, user.LastName, user.Email, user.PasswordHash, user.PasswordSalt)

	// Достанем из ответа БД переменную id
	var user_id int
	if err := row.Scan(&user_id); err != nil {
		tx.Rollback()
		return 0, err // нет переменной id
	}

	// Создание пустого аккаунта
	a := domain.Account{
		// ID:        uuid.NewRandom().String(),
		// Name:      req.Name,
		// Address1:  req.Address1,
		// Address2:  req.Address2,
		// City:      req.City,
		// Region:    req.Region,
		// Country:   req.Country,
		// Zipcode:   req.Zipcode,
		Status: domain.AccountStatus_Active,
		// Timezone:      "America/Anchorage",
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		SignupUserID:  &sql.NullInt64{Int64: int64(user_id), Valid: true},
		BillingUserID: &sql.NullInt64{Int64: int64(1), Valid: true}, // Привязываем нового пользователя к аккаунту AirSoft
		UserToken:     randomstring.RandomBase64String(16),
	}

	// Execute account creation.
	query := sqlbuilder.NewInsertBuilder() // Build the insert SQL statement.
	query.InsertInto(accountTable)
	query.Cols("status", "signup_user_id", "billing_user_id", "created_at", "updated_at", "u_token")
	query.Values(a.Status.String(), a.SignupUserID, a.BillingUserID, a.CreatedAt, a.UpdatedAt, a.UserToken)
	// query.Cols("id", "name", "address1", "address2", "city", "region", "country", "zipcode", "status", "timezone", "signup_user_id", "billing_user_id", "created_at", "updated_at")
	// query.Values(a.ID, a.Name, a.Address1, a.Address2, a.City, a.Region, a.Country, a.Zipcode, a.Status.String(), a.Timezone, a.SignupUserID, a.BillingUserID, a.CreatedAt, a.UpdatedAt)

	// Execute the query with the provided context.
	sql, args := query.Build()               // "INSERT INTO accounts (status, signup_user_id, billing_user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	sql = r.db.Rebind(sql) + " RETURNING id" // "INSERT INTO accounts (status, signup_user_id, billing_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"

	rows := tx.QueryRow(sql, args...)

	// Достанем из ответа БД переменную id
	var account_id int
	err = rows.Scan(&account_id)
	if err != nil {
		err = errors.Wrapf(err, "query - %s", query.String())
		err = errors.WithMessage(err, "create account failed")
		tx.Rollback()
		return 0, err // нет переменной id
	}

	// Создание дефолтной роли пользователя
	// Associate the created user with the new account. The first user for the account will
	// always have the role of manager.
	ua := domain.UserAccount{
		//ID:        uuid.NewRandom().String(),
		UserID:    user_id,
		AccountID: account_id,
		Roles:     []domain.UserAccountRole{domain.UserAccountRole_User}, // {user_account.UserAccountRole_Manager},
		Status:    domain.UserAccountStatus_Active,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Build the insert SQL statement.
	query = sqlbuilder.NewInsertBuilder()
	query.InsertInto(userAccountTableName)
	query.Cols("user_id", "account_id", "roles", "status", "created_at", "updated_at")
	query.Values(ua.UserID, ua.AccountID, ua.Roles, ua.Status.String(), ua.CreatedAt, ua.UpdatedAt)

	// Execute the query with the provided context.
	sql, args = query.Build()
	sql = r.db.Rebind(sql) // "INSERT INTO users_accounts (user_id, account_id, roles, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = tx.Exec(sql, args...)
	if err != nil {
		err = errors.Wrapf(err, "query - %s", query.String())
		err = errors.WithMessagef(err, "add account %d to user %d failed", account_id, user_id)
		tx.Rollback()
		return 0, err
	}

	return user_id, tx.Commit()
}

// gets a unique user by email from the database.
func (r *AuthPostgres) GetUserByEmail(email string) (*domain.User, error) {

	var user []domain.User
	query := fmt.Sprintf("SELECT id, email, first_name, last_name, timezone, password_hash, password_salt, password_reset, created_at, updated_at, archived_at FROM %s WHERE email=$1", usersTable)
	err := r.db.Select(&user, query, email)
	if err != nil {
		return nil, err
	}

	if len(user) == 0 || user[0].ID == 0 {
		return nil, errors.New(fmt.Sprintf("user %s not found", email))
	}

	return &user[0], nil
}

// Запросить в БД пользователя
func (r *AuthPostgres) GetUser_LoginPassword(email, password string) (domain.User, error) {
	var user domain.User
	// query := fmt.Sprintf("SELECT id FROM %s WHERE login=$1 AND password_hash=$2", accountsTable)
	query := fmt.Sprintf("SELECT id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, email, password)

	return user, err
}

// Запросить по токену список связку:
//    ID_шников пользователя, аквахаба, устройств, сенсоров БД
//    с локальными ID_шниками устройств и сенсоров
// user_id & aquahub_id & device_id & sensor_id <===> local_device_id & local_sensor_id
func (r *AuthPostgres) GetUserHW_fromTokens(h_token, u_token string) ([]domain.SensorDataSet, error) { // user_id, aquahub_id, []{device_id, sensor_id}

	var lists []domain.SensorDataSet

	// Команда INNER JOIN при SELECT помогает выбрать только те записи,
	// которые имеют одинаковое значение в обеих таблицах.

	// INNER JOIN devices dlt on dlt.aquahub_id = aht.id
	// INNER JOIN sensors slt on slt.device_id = dlt.id
	// WHERE a.u_token = 'a39831d103eb4c0d' AND aht.h_token = 'aqen104Ur2zNX1Ykwv4'
	// ORDER BY aht.id, dlt.id, slt.id;

	query := fmt.Sprintf(`SELECT aht.account_id AS account_id, aht.id AS aquahub_id, dlt.id AS device_id, slt.id AS sensor_id, dlt.local_id AS local_device_id, slt.local_id AS local_sensor_id
							FROM  %s aht
							INNER JOIN %s a on aht.account_id = a.id
							INNER JOIN %s dlt on dlt.aquahub_id = aht.id
							INNER JOIN %s slt on slt.device_id = dlt.id
							WHERE a.u_token = $2 AND aht.h_token = $1
							ORDER BY aht.id, dlt.id, slt.id`,
		aquahubsTable, accountTable, devicesTable, sensorsTable)

	// fmt.Printf("\n\n%s\n\n", query)

	err := r.db.Select(&lists, query, h_token, u_token)

	return lists, err
}

// Запросить по токенам ID аквахаба
func (r *AuthPostgres) GetAquahubId_fromTokens(h_token, u_token string) (int, error) {

	// Команда INNER JOIN при SELECT помогает выбрать только те записи,
	// которые имеют одинаковое значение в обеих таблицах.
	query := fmt.Sprintf(`SELECT aht.id
							FROM  %s aht
							INNER JOIN %s ut on aht.user_id = ut.id
							WHERE aht.h_token = $1 AND ut.u_token = $2`,
		aquahubsTable, usersTable)

	// fmt.Printf("\n\n%s\n\n", query)

	var aId int
	err := r.db.Get(&aId, query, h_token, u_token)

	return aId, err
}
