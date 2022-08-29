package main

import (
	"github.com/AquaEngineering/AquaHub/internal/app"
	"github.com/sirupsen/logrus"
)

// Обозначить тэг для сваггера pq.NullTime ===>>> `swaggertype:"string"`

// @title                      AquaHub API
// @version                    1.0
// @description                API Server for AquaHub
// @host                       localhost:8000
// @BasePath                   /
// securityDefinitions.apikey ApiKeyAuth
// in                         header
// name                       Authorization
func main() {

	var log = logrus.New()

	// Задаём формат лога JSON
	// log.Formatter = new(logrus.JSONFormatter)
	// log.Formatter = new(logrus.TextFormatter)                     // default
	// log.Formatter.(*logrus.TextFormatter).DisableColors = true    // remove colors
	// log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	// log.Level = logrus.TraceLevel
	// log.SetReportCaller(true) // печать номеров строк

	// log.Out = os.Stdout

	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err == nil {
	// 	log.Out = file
	// } else {
	// 	log.Info("Failed to log to file, using default stderr")
	// }

	App, err := app.NewApplication(log)
	if err != nil {
		log.Fatalf("--------------- App init error: %+v", err)
	}

	/*
		// Print the config for our logs. It's important to any credentials in the config
		// that could expose a security risk are excluded from being json encoded by
		// applying the tag `json:"-"` to the struct var.

		// Печатаем конфиг.
		// Важно, чтобы любые учетные данные в конфигурации,
		// которые могут представлять угрозу безопасности,
		// были исключены из json путем применения тега `json:"-"` к структуре cfg.
		cfgJSON, err := json.MarshalIndent(App.Cfg, "", "    ")
		if err != nil {
			log.Fatalf("--------------- App: Marshalling Config to JSON : %+v", err)
		}
		fmt.Printf("main: Config : %v\n", string(cfgJSON))
	*/

	err = App.Run()
	if err != nil {
		log.Printf("--------------- App: Server Run : error : %+v", err)
	} else {
		log.Print("App: Server closed successful")
	}

	App.Shutdown()
}
