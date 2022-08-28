package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AquaEngineering/AquaHub/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ChecklistItemPostgres struct {
	db *sqlx.DB
}

func NewChecklistItemPostgres(db *sqlx.DB) *ChecklistItemPostgres {
	return &ChecklistItemPostgres{db: db}
}

func (r *ChecklistItemPostgres) Create(listId int, item domain.ChecklistItem) (int, error) {

	createItemQuery := fmt.Sprintf("INSERT INTO %s (checklist_id, title, description) values ($1, $2, $3) RETURNING id",
		checklistItemsTable)

	row := r.db.QueryRow(createItemQuery, listId, item.Title, item.Description)

	var itemId int
	err := row.Scan(&itemId)
	if err != nil {
		return 0, err
	}

	return itemId, nil
}

func (r *ChecklistItemPostgres) GetAll(userId, listId int) ([]domain.ChecklistItem, error) {

	// query := fmt.Sprintf(`SELECT it.id, it.title, it.description, it.done FROM %s ti INNER JOIN %s clt on clt.checklist_id = it.id
	// 								INNER JOIN %s a on a.id = clt.account_id WHERE clt.id = $1 AND a.signup_user_id = $2`,

	query := fmt.Sprintf(`SELECT it.id, it.title, it.description, it.done, it.updated_at FROM %s it 
		INNER JOIN %s clt on clt.id = it.checklist_id
		INNER JOIN %s a on a.id = clt.account_id 
		WHERE clt.id = $1 AND a.signup_user_id = $2`,
		checklistItemsTable, checklistsTable, accountTable)

	var items []domain.ChecklistItem
	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ChecklistItemPostgres) GetById(userId, listId, itemId int) (domain.ChecklistItem, error) {
	query := fmt.Sprintf(`SELECT it.id, it.title, it.description, it.done, it.updated_at FROM %s it INNER JOIN %s clt on clt.id = it.checklist_id
									INNER JOIN %s a on a.id = clt.account_id WHERE it.id = $1 AND clt.id = $2 AND a.signup_user_id = $3`,
		checklistItemsTable, checklistsTable, accountTable)

	var item domain.ChecklistItem
	if err := r.db.Get(&item, query, itemId, listId, userId); err != nil {
		return item, err
	}

	return item, nil
}

func (r *ChecklistItemPostgres) Update(userId, listId, itemId int, input domain.UpdateChecklistItem) error {

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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s it SET %s FROM %s clt, %s a
						  WHERE it.checklist_id = clt.id AND clt.account_id = a.id 
								AND a.signup_user_id = $%d AND clt.id = $%d AND it.id = $%d RETURNING it.id`,
		checklistItemsTable, setQuery, checklistsTable, accountTable,
		argId, argId+1, argId+2)
	args = append(args, userId, listId, itemId)

	var item domain.ChecklistItem
	if err := r.db.Get(&item, query, args...); err != nil {
		return err
	}

	if item.ID == 0 {
		return errors.New("update error")
	}

	return nil
}

func (r *ChecklistItemPostgres) Delete(userId, listId, itemId int) error {

	query := fmt.Sprintf(`DELETE FROM %s it USING %s clt, %s a 
							WHERE it.checklist_id = clt.id AND clt.account_id = a.id 
							AND a.signup_user_id = $1 AND clt.id = $2 AND it.id = $3`,
		checklistItemsTable, checklistsTable, accountTable)

	_, err := r.db.Exec(query, userId, listId, itemId)

	return err
}
