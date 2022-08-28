package domain

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

// Поля полностью совпадают с БД кроме пароля - в БД хеш пароля

// Добавлены JSON теги чтобы корректно принимать и выводить данные в HTTP-запросах:
// Тег binding:"required" - валидирует наличие данного поля в теле запроса (является реализацией фремворка gin)
/*
type User_old struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
	Token    string `json:"token" binding:"required"`
}
*/

type User struct {
	ID            int             `json:"-" validate:"required" db:"id" example:"1"`
	FirstName     string          `json:"first_name" validate:"required" db:"first_name" example:"Andy"`
	LastName      string          `json:"last_name" validate:"required" db:"last_name" example:"Sokol"`
	Email         string          `json:"email" validate:"required,email,unique" db:"email" example:"ae@ae.ae"`
	PasswordSalt  string          `json:"-" validate:"required" db:"password_salt"`
	PasswordHash  string          `json:"password" validate:"required" db:"password_hash"  example:"jyWtbKg76by"`
	PasswordReset *sql.NullString `json:"-" db:"password_reset"`
	Timezone      *string         `json:"-" validate:"omitempty" db:"timezone" example:"Europe/Kiev"`
	CreatedAt     time.Time       `json:"-" db:"created_at"`
	UpdatedAt     time.Time       `json:"-" db:"updated_at"`
	ArchivedAt    *pq.NullTime    `json:"-,omitempty" db:"archived_at" swaggertype:"string"`
}

func (v User) Validate() error {
	if v.FirstName == "" || v.LastName == "" || v.Email == "" || v.PasswordHash == "" {
		return errors.New("not all values are specified")
	}

	return nil
}

/*
CREATE  TABLE x_users (
	id            		serial   NOT NULL unique,
	email                varchar(200)  NOT NULL  ,
	first_name           varchar(200) DEFAULT ''::character varying NOT NULL  ,
	last_name            varchar(200) DEFAULT ''::character varying NOT NULL  ,
	timezone             varchar(128) DEFAULT 'Europe/Kiev'::character varying NOT NULL,
	password_hash        varchar(256)  NOT NULL  ,
	password_salt        varchar(36)  NOT NULL  ,
	password_reset       varchar(36) DEFAULT NULL::character varying   ,
	created_at	timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at      timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	archived_at     timestamptz,
	CONSTRAINT users_pkey PRIMARY KEY ( id ),
	CONSTRAINT email UNIQUE ( email )
 );
*/
