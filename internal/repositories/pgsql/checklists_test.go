package repository

import (
	"fmt"
	"testing"

	"github.com/o-sokol-o/hub/internal/domain"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func getPointerString(s string) *string {
	return &s
}

func TestChecklistListPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type args struct {
		userId int
		list   domain.CreateChecklist
	}
	tests := []struct {
		name    string
		mock    func(userId int, list domain.CreateChecklist)
		input   args // Входящие аргументы - параметры вызова метода
		want_id int  // ожидаемое возвращаемое значение id
		wantErr bool // Ожидаем ли мы получить ошибку или нет
	}{
		{
			name: "Test case: Ok",
			// Описываем ожидаемое поведение объекта базы данных
			mock: func(userId int, list domain.CreateChecklist) {
				mock.ExpectBegin()

				//  fmt.Sprintf(`SELECT id FROM %s WHERE signup_user_id = $1`, accountTable)
				query := fmt.Sprintf(`SELECT (.+) FROM %s WHERE (.+)`, accountTable)
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1234) // результат: account_id = 1234
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)

				// fmt.Sprintf("INSERT INTO %s (account_id, title, description) VALUES ($1, $2, $3) RETURNING id", checklistsTable)
				query = fmt.Sprintf("INSERT INTO %s", checklistsTable)
				rows = sqlmock.NewRows([]string{"id"}).AddRow(4321) // результат: (checklists)id = 4321
				mock.ExpectQuery(query).
					WithArgs(1234, list.Title, list.Description). // 1234 - результат после первого селекта
					WillReturnRows(rows)

				// mock.ExpectExec("INSERT INTO users_lists").WithArgs(1, 1).
				// 	WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			input: args{
				userId: 1,
				list: domain.CreateChecklist{
					Title:       getPointerString("title"),
					Description: getPointerString("description"),
				},
			},

			want_id: 4321,
			wantErr: false, // Ошибку не ожидаем
		},
		{
			name: "Test case: Empty Fields",
			mock: func(userId int, list domain.CreateChecklist) {
				mock.ExpectBegin()

				//  fmt.Sprintf(`SELECT id FROM %s WHERE signup_user_id = $1`, accountTable)
				query := fmt.Sprintf(`SELECT (.+) FROM %s WHERE (.+)`, accountTable)
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1234) // результат: account_id = 1234
				mock.ExpectQuery(query).WithArgs(userId).WillReturnRows(rows)

				// fmt.Sprintf("INSERT INTO %s (account_id, title, description) VALUES ($1, $2, $3) RETURNING id", checklistsTable)
				query = fmt.Sprintf("INSERT INTO %s", checklistsTable)
				rows = sqlmock.NewRows([]string{"id"}) // результата нет
				mock.ExpectQuery(query).
					WithArgs(1234, list.Title, list.Description). // 1234 - результат после первого селекта
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			input: args{
				userId: 1,
				list: domain.CreateChecklist{
					Title:       getPointerString(""),
					Description: getPointerString("description"),
				},
			},
			wantErr: true, // Ошибку ожидаем
		},
	}

	r := NewChecklistPostgres(logrus.New(), db)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.input.userId, test.input.list)

			got, err := r.Create(test.input.userId, test.input.list)
			// Ожидаем ошибку ?
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want_id, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

/*
func TestChecklistPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewChecklistPostgres(logrus.New(), db)

	type args struct {
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []domain.Checklist
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "title1", "description1").
					AddRow(2, "title2", "description2").
					AddRow(3, "title3", "description3")

				mock.ExpectQuery("SELECT (.+) FROM todo_lists tl INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)
			},
			input: args{
				userId: 1,
			},
			want: []domain.Checklist{
				{1, "title1", "description1"},
				{2, "title2", "description2"},
				{3, "title3", "description3"},
			},
		},
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "title1", "description1").
					AddRow(2, "title2", "description2").
					AddRow(3, "title3", "description3")

				mock.ExpectQuery("SELECT (.+) FROM todo_lists tl INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)
			},
			input: args{
				userId: 1,
			},
			want: []domain.Checklist{
				{1, "title1", "description1"},
				{2, "title2", "description2"},
				{3, "title3", "description3"},
			},
		},
	}

	for _, test  := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.GetAll_ChecklistOfUser(test.input.userId)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestChecklistPostgres_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewChecklistPostgres(logrus.New(), db)

	type args struct {
		listId int
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    domain.Checklist
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "title1", "description1")

				mock.ExpectQuery("SELECT (.+) FROM todo_lists tl INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			input: args{
				listId: 1,
				userId: 1,
			},
			want: domain.Checklist{1, "title1", "description1"},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"})

				mock.ExpectQuery("SELECT (.+) FROM todo_lists tl INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1, 404).WillReturnRows(rows)
			},
			input: args{
				listId: 404,
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, test  := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.GetById(test.input.userId, test.input.listId)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestChecklistPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewChecklistPostgres(logrus.New(), db)

	type args struct {
		listId int
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("DELETE FROM todo_lists tl USING users_lists ul WHERE (.+)").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectExec("DELETE FROM todo_lists tl USING users_lists ul WHERE (.+)").
					WithArgs(1, 404).WillReturnError(sql.ErrNoRows)
			},
			input: args{
				listId: 404,
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, test  := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			err := r.Delete(test.input.userId, test.input.listId)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestChecklistPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewChecklistPostgres(logrus.New(), db)

	type args struct {
		listId int
		userId int
		input  domain.UpdateListInput
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "OK_AllFields",
			mock: func() {
				mock.ExpectExec("UPDATE todo_lists tl SET (.+) FROM users_lists ul WHERE (.+)").
					WithArgs("new title", "new description", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
				input: domain.UpdateListInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
				},
			},
		},
		{
			name: "OK_WithoutDescription",
			mock: func() {
				mock.ExpectExec("UPDATE todo_lists tl SET (.+) FROM users_lists ul WHERE (.+)").
					WithArgs("new title", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
				input: domain.UpdateListInput{
					Title: stringPointer("new title"),
				},
			},
		},
		{
			name: "OK_WithoutTitle",
			mock: func() {
				mock.ExpectExec("UPDATE todo_lists tl SET (.+) FROM users_lists ul WHERE (.+)").
					WithArgs("new description", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
				input: domain.UpdateListInput{
					Description: stringPointer("new description"),
				},
			},
		},
		{
			name: "OK_NoInputFields",
			mock: func() {
				mock.ExpectExec("UPDATE todo_lists tl SET FROM users_lists ul WHERE (.+)").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				listId: 1,
				userId: 1,
			},
		},
	}

	for _, test  := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			err := r.Update(test.input.userId, test.input.listId, test.input.input)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
*/
