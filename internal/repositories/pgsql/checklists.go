package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AquaEngineering/AquaHub/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ChecklistPostgres struct {
	db  *sqlx.DB
	log *logrus.Logger
}

func NewChecklistPostgres(log *logrus.Logger, db *sqlx.DB) *ChecklistPostgres {
	return &ChecklistPostgres{log: log, db: db}
}

func (r *ChecklistPostgres) Create(userId int, list domain.UpdateChecklist) (int, error) {

	// Транзакция это последовательность нескольких операций,
	// которая рассматривается как одна операция.
	// Транзакция либо выполняется целиком, либо не выполняется вообще.

	// Создадим новую транзакцию, в которой выполним 2 операции.
	tx, err := r.db.Begin()
	if err != nil {
		r.log.Errorf("db: error Create Checklist: %s", err.Error())
		return 0, errors.New("db: error Create Checklist")
	}

	query := fmt.Sprintf(`SELECT id FROM %s WHERE signup_user_id = $1`, accountTable)
	var account_id int
	err = r.db.Get(&account_id, query, userId)
	if err != nil {
		tx.Rollback() // В случае ошибок мы вызываем метод Rollback, которая откатывает все изменения БД до начала выполнения транзакции.

		r.log.Errorf("db: error Create Checklist: %s", err.Error())
		return 0, errors.New("db: error Create Checklist")
	}

	// Возвращаем id созданного списка
	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (account_id, title, description) VALUES ($1, $2, $3) RETURNING id", checklistsTable)
	row := tx.QueryRow(createListQuery, account_id, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()

		r.log.Errorf("db: error Create Checklist: %s", err.Error())
		return 0, errors.New("db: error Create Checklist")
	}

	// После выполнения транзакции вызовем метод Commit,
	// который применит изменения к БД и закончит транзакцию.
	return id, tx.Commit()
}

//_____________________________________________________________________________________________________

func (r *ChecklistPostgres) GetAll_ChecklistOfUser(userId int) ([]domain.Checklist, error) {

	var lists []domain.Checklist

	// В запросе сделаем выборку из базы, используя конструкцию INNER JOIN.
	// Команда INNER JOIN при SELECT помогает выбрать только те записи,
	// которые имеют одинаковое значение в обеих таблицах.
	// В нашем случае нам нужно выбрать запись из таблицы checklistsTable,
	// которая принадлежит аккаунту из таблицы accountTable,
	// и при этом они связаны по id аккаунта, а аккаунт принадлежит пользователю.
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description, tl.status, tl.updated_at FROM %s tl INNER JOIN %s a on a.id = tl.account_id WHERE a.signup_user_id = $1",
		checklistsTable, accountTable)

	// Результат записываем в слайс.
	err := r.db.Select(&lists, query, userId)

	if err != nil {
		r.log.Errorf("db: error GetAll Checklist: %s", err.Error())
		return nil, errors.New("db: error GetAll Checklist")
	}

	return lists, err
}

func (r *ChecklistPostgres) GetById(userId, listId int) (*domain.Checklist, error) {

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description, tl.status, tl.updated_at FROM %s tl INNER JOIN %s a on a.id = tl.account_id WHERE a.signup_user_id = $1 AND tl.id = $2`,
		checklistsTable, accountTable)

	var list domain.Checklist
	err := r.db.Get(&list, query, userId, listId)

	if err != nil {
		r.log.Errorf("db: error GetById Checklist: %s", err.Error())
		return nil, errors.New("db: error GetById Checklist")
	}

	return &list, err
}

func (r *ChecklistPostgres) Update(userId int, input domain.UpdateChecklist) error {

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

	if input.Status != nil {
		setValues = append(setValues, fmt.Sprintf("status=$%d", argId))
		args = append(args, *input.Status)
		argId++
	}

	{
		setValues = append(setValues, fmt.Sprintf("updated_at=$%d", argId))
		args = append(args, time.Now())
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s a WHERE tl.account_id = a.id AND a.signup_user_id = $%d AND tl.id = $%d RETURNING tl.id",
		checklistsTable, setQuery, accountTable, argId, argId+1)
	args = append(args, userId, input.ID)

	r.log.Debugf("updateQuery: %s", query)
	r.log.Debugf("args: %s", args)

	var id int
	err := r.db.Get(&id, query, args...)
	if err != nil {
		r.log.Errorf("db: error Update Checklist: %s", err.Error())
		return errors.New("db: error Update Checklist")
	}

	if id != input.ID {
		r.log.Error("db: error Update Checklist: id != input.ID")
		return errors.New("db: error Update Checklist")
	}

	return nil
}

func (r *ChecklistPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s a WHERE tl.account_id = a.id AND a.signup_user_id = $1 AND tl.id = $2",
		checklistsTable, accountTable)
	_, err := r.db.Exec(query, userId, listId)

	if err != nil {
		r.log.Errorf("db: error Update Checklist: %s", err.Error())
		return errors.New("db: error Delete Checklist")
	}

	return nil
}
