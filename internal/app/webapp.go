package app

import (
	"context"
	"errors"
	"expvar"
	"os"

	"github.com/jmoiron/sqlx"

	cachememory "github.com/o-sokol-o/cache-memory"
	cmd_line "github.com/o-sokol-o/hub/pkg/cmd_line"
	database "github.com/o-sokol-o/hub/pkg/database/postgres"
	webserver "github.com/o-sokol-o/hub/pkg/server"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	repositories "github.com/o-sokol-o/hub/internal/repositories/pgsql"
	services "github.com/o-sokol-o/hub/internal/services"
	handlers "github.com/o-sokol-o/hub/internal/transports"

	config "github.com/o-sokol-o/hub/internal/domain"
)

//=====================================================================================================================
//=====================================================================================================================

// build is the git version of this program. It is set using build flags in the makefile.

// сборка — это git-версия этой программы. Он устанавливается с помощью флагов сборки в make-файле.
var build = "develop"

// service is the name of the program used for logging, tracing and the
// the prefix used for loading env variables
// ie: export WEB_APP_ENV=dev

// служба — это имя программы, используемой для логирования,
// трассировки и префикса, используемого для загрузки переменных env
// то есть: экспорт WEB_APP_ENV=dev
var serviceTagName = "WEB_APP"

//=====================================================================================================================
//=====================================================================================================================

type AppContext struct {
	cfg         config.CfgApplication
	log         *logrus.Logger
	server      *webserver.Server
	cacheMemory cachememory.Cache
	masterDB    *sqlx.DB
	handlers    *handlers.Handler

	// Env          webcontext.Env
	// MasterDbHost string
	// Redis             *redis.Client
	// UserRepo *user.Repository
	// UserAccountRepo   *user_account.Repository
	// AccountRepo       *account.Repository
	// AccountPrefRepo   *account_preference.Repository
	// AuthRepo 		 *user_auth.Repository
	// SignupRepo *signup.Repository
	// InviteRepo        *invite.Repository
	// ChecklistRepo     *checklist.Repository
	// GeoRepo *geonames.Repository
	// Authenticator     *auth.Authenticator
	// StaticDir   string
	// TemplateDir string
	// Renderer    web.Renderer
	// WebRoute          webroute.WebRoute
	// PreAppMiddleware  []web.Middleware
	// PostAppMiddleware []web.Middleware
	// AwsSession *session.Session
}

//=====================================================================================================================
//=====================================================================================================================

func NewApplication(log *logrus.Logger) (*AppContext, error) {

	var app = AppContext{
		log:      log,
		masterDB: nil, // *sqlx.DB,
		// Env:                 webcontext.Env,
		// MasterDbHost:        "",
		// Redis:        		*redis.Client,
		// UserRepo:          	*user.Repository,
		// UserAccountRepo:   	*user_account.Repository,
		// AccountRepo:       	*account.Repository,
		// AccountPrefRepo:   	*account_preference.Repository,
		// AuthRepo:          	*user_auth.Repository,
		// SignupRepo:        	*signup.Repository,
		// InviteRepo:        	*invite.Repository,
		// ChecklistRepo:     	*checklist.Repository,
		// GeoRepo:           	*geonames.Repository,
		// Authenticator:     	*auth.Authenticator,
		// StaticDir:   "",
		// TemplateDir: "",
		// Renderer:    		web.RendereInit,
		// WebRoute          	webroute.WebRoute,
		// PreAppMiddleware  	[]web.Middleware,
		// PostAppMiddleware 	[]web.Middleware,
		// AwsSession 			*session.Session,
	}

	// =============   Работа с конфигом и окружением   =============
	// TODO: Надо ли рефакторить?
	{
		// Пакет envconfig реализует декодирование переменных среды на основе заданной
		// пользователем спецификации. Типичное использование — использование переменных
		// среды для параметров конфигурации.
		// For additional details refer to https://github.com/kelseyhightower/envconfig
		if err := envconfig.Process(serviceTagName, &app.cfg); err != nil {
			log.Printf("main : Parsing Config : %+v", err)
		}

		// Process compares the specified command line arguments against the provided
		// struct value and updates the fields that are identified.

		// Процесс сравнивает указанные аргументы командной строки с
		// предоставленным значением структуры и обновляет идентифицированные поля.
		if err := cmd_line.Process(&app.cfg); err != nil {
			if err != cmd_line.ErrHelp {
				log.Printf("main : Parsing Command Line : %+v", err)
			}
			return nil, err // We displayed help.
		}

		// =========================================================================
		// configValidationAndDefaults

		// If base URL is empty, set the default value from the HTTP Host
		app.cfg.Service.BaseUrl = baseURLtoDefault(app.cfg.Service.BaseUrl, app.cfg.HTTP.Host)

		// When HTTPS is not specifically enabled, but an HTTP host is set, enable HTTPS.
		if !app.cfg.Service.EnableHTTPS && app.cfg.HTTPS.Host != "" {
			app.cfg.Service.EnableHTTPS = true
		}

		// =========================================================================
		// Log Service Info

		// Print the build version for our logs. Also expose it under /debug/vars.
		expvar.NewString("build").Set(build)
		log.Printf("main : Started : Service Initializing version %q", build)
		defer log.Println("main : Completed")

		// =========================================================================
		// Инициализируем и читаем конфиг

		viper.AddConfigPath("configs")
		viper.SetConfigName("config")
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("error initializing configs: %s", err.Error())
			return nil, err
		}

		// Инициализируем и читаем окружение
		if err := godotenv.Load(); err != nil {
			log.Printf("error loading env variables: %s", err.Error())
			return nil, err
		}

		app.cfg.HTTP.Port = os.Getenv("PORT")
		if app.cfg.HTTP.Port == "" {
			app.cfg.HTTP.Port = os.Getenv("LOCALPORT")
		}
	}

	// =============   Инициализируем кэш и подключаемся к БД   =============
	{
		var err error

		// Кэш
		app.cacheMemory = cachememory.New(3600) // 1 hour

		// Постгрес
		app.masterDB, err = database.NewDB(context.TODO(),
			viper.GetString("db.host"), // из конфига
			viper.GetString("db.port"),
			viper.GetString("db.username"),
			viper.GetString("db.password"), //os.Getenv("DB_PASSWORD"), // из окружения
			viper.GetString("db.dbname"),
			viper.GetString("db.sslmode"),
		)
		if err != nil {
			return nil, errors.New("Failed to initialize db: " + err.Error())
		}

		// TODO: Добавить чтение кэша из базы
		// AppCtx.MasterDB.LoadCacheData(AppCtx.CacheMemory)

		// TODO: Добавить в базу кэширование запросов
		// AppCtx.MasterDB.SetCache(AppCtx.CacheMemory)
	}

	// =============   Внедряем зависимости согласно "Чистой архитектуре"   =============

	// Внутренний слой - БД   - не зависит ни от чего
	// Средний слой - Service - бизнес логика общается только со слоем БД
	// Внешний слой - Handler - обработчики общаются только со слоем бизнес логики

	app.handlers =
		handlers.NewHandler(
			services.NewServices(
				repositories.NewRepositories(app.log, app.cacheMemory, app.masterDB)))

	err := app.handlers.InitRoutes()
	if err != nil {
		log.Printf("InitRoutes error: %s", err.Error())
		return nil, err
	}

	return &app, nil
}

//_____________________________________________________________________________________________________________________

func (app *AppContext) Run() error {

	app.server = webserver.NewServer(app.log)

	return app.server.Run(app.cfg.HTTP.Port, app.handlers.Router)
}

//_____________________________________________________________________________________________________________________

func (app *AppContext) Shutdown() {

	app.log.Print("Application shutting down ...")

	// Плавное завершение работы должно гарантировать, что при выходе из приложения
	// мы перестанем принимать все входящие и запросы, но при этом закончим
	// обработку всех текущих запросов и операций в базе данных.

	// Вызовем 2 метода остановки сервера и закрытие всех соединений с базой данных.
	// Это гарантирует нам, что мы закончим выполнение всех текущих операций перед выходом из приложения.
	err := app.server.Shutdown()
	if err != nil {
		app.log.Printf("error occured on server shutting down: %s", err.Error())
	}

	// TODO: Добавить запись кэша в базу
	// AppCtx.MasterDB.SaveCacheData(AppCtx.CacheMemory)

	// Закрыть базу
	if err = app.masterDB.Close(); err != nil {
		app.log.Printf("Error occured on db connection close: %s", err.Error())
	} else {
		app.log.Printf("DB disconnection successful")
	}

	app.log.Print("Application shutting down ended")
}
