package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/o-sokol-o/hub/internal/domain"
	"github.com/sirupsen/logrus"
)

// Список используемых таблиц
const (
	accountTable               = "accounts"
	userAccountTableName       = "users_accounts"
	accountPreferenceTableName = "account_preferences"

	checklistsTable     = "checklists"
	checklistItemsTable = "checklist_items"
	usersTable          = "users"

	aquahubsTable      = "aquahubs"
	devicesTable       = "devices"
	sensorsTable       = "sensors"
	sensorDataSetTable = "sensors_dataset"
)

func NewRepositories(log *logrus.Logger, cache domain.Cache, db *sqlx.DB) (
	*logrus.Logger, domain.Cache,

	*AuthPostgres,
	*ChecklistPostgres,
	*ChecklistItemPostgres,
	*AquahubListPostgres) {

	return log, cache,

		NewAuthPostgres(db),
		NewChecklistPostgres(log, db),
		NewChecklistItemPostgres(db),
		NewAquahubListPostgres(log, db)
}
