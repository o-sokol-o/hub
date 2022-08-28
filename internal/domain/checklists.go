package domain

import (
	"errors"
	"time"

	"github.com/lib/pq"
)

// Поля полностью совпадают с БД
// Добавлены JSON теги чтобы корректно принимать и выводить данные в HTTP-запросах
// Тег binding:"required" - валидирует наличие данного поля в теле запроса (является реализацией фремворка gin)

/*
CREATE  TABLE checklists (
	id            		serial   NOT NULL unique,
	account_id           	integer  NOT NULL,
	name                 	varchar(255)  NOT NULL  ,
	status               	project_status_t DEFAULT 'active'::project_status_t NOT NULL  ,
	created_at	timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at      timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	archived_at     timestamptz,
	CONSTRAINT projects_pkey PRIMARY KEY ( id )
 );

ALTER TABLE checklists ADD CONSTRAINT projects_account_id_fkey FOREIGN KEY ( account_id ) REFERENCES accounts( id ) ON DELETE SET NULL;
*/

type Checklist struct {
	ID          int             `json:"id" db:"id" validate:"required" example:"1"`
	AccountID   int             `json:"-" db:"account_id" validate:"required" truss:"api-create"`
	Title       string          `json:"title" db:"title" validate:"required" example:"Rocket Launch"`
	Description string          `json:"description" db:"description" example:"Rocket Launch Description"`
	Status      ChecklistStatus `json:"status" db:"status" validate:"omitempty,oneof=active disabled" enums:"active,disabled" swaggertype:"string" example:"active"`
	CreatedAt   time.Time       `json:"-" db:"created_at" truss:"api-read"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at" truss:"api-read"`
	ArchivedAt  *pq.NullTime    `json:"-" db:"archived_at" truss:"api-hide" swaggertype:"string"`
}

// ChecklistStatus represents the status of checklist.
type ChecklistStatus string

func (i Checklist) Validate() error {
	if i.Title == "" || i.Description == "" {
		return errors.New("the input data does not have the desired values")
	}

	return nil
}

type UpdateChecklist struct {
	ID          int              `json:"-"`
	Title       *string          `json:"title" example:"Title Checklist"`
	Description *string          `json:"description" example:"Description Checklist"`
	Status      *ChecklistStatus `json:"status,omitempty" validate:"omitempty,oneof=active disabled" enums:"active,disabled" swaggertype:"string" example:"active"`
}

func (i UpdateChecklist) Validate() error {
	if (i.Title == nil && i.Description == nil && i.Status == nil) || i.ID == 0 {
		return errors.New("update has no values")
	}

	return nil
}

type ChecklistItem struct {
	ID          int          `json:"-" db:"id"`
	ChecklistID int          `json:"-" db:"checklist_id" validate:"required" truss:"api-create"`
	Title       string       `json:"title,omitempty" db:"title" validate:"required" example:"Rocket Launch"`
	Description string       `json:"description,omitempty" db:"description"`
	Done        bool         `json:"done,omitempty" db:"done"`
	CreatedAt   time.Time    `json:"-" db:"created_at" truss:"api-read"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty" db:"updated_at" truss:"api-read"`
	ArchivedAt  *pq.NullTime `json:"-,omitempty" db:"archived_at" truss:"api-hide" swaggertype:"string"`
}

type UpdateChecklistItem struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (i UpdateChecklistItem) Validate() error {
	if i.Title == nil && i.Description == nil && i.Done == nil {
		return errors.New("update has no values")
	}

	return nil
}
