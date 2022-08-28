package webserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

//=====================================================================================================================

type Server struct {
	log        *logrus.Logger
	HttpServer *http.Server
}

func NewServer(log *logrus.Logger) *Server {
	return &Server{
		log: log,
		HttpServer: &http.Server{
			// Инкапсуляция параметров сервера
			MaxHeaderBytes: 1 << 20, // 1 MB
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
		},
	}
}

// Неблокирующее чтение канала
func readChanNoBlock(ch chan error) (result error, err error) {
	select {
	case result = <-ch:
		return result, nil
	default:
		return nil, errors.New("channel has no data")
	}
}

func (s *Server) Run(port string, handlers http.Handler) error {

	if handlers == nil {
		fmt.Printf("Handler to invoke = http.DefaultServeMux (no custom)")
	}
	s.HttpServer.Addr = ":" + port
	s.HttpServer.Handler = handlers

	// С помощью канала типа os.Signal в дальнейшем будем ожидать завершения вебсервера.
	// Запись канала будет происходить, когда процесс, в котором выполняется наше приложение,
	// получит сигнал от системы типа SIGINT или SIGTERM - это сигналы юникс системах.
	ch_err := make(chan error, 1)
	ch_quit := make(chan os.Signal, 1)
	signal.Notify(ch_quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)

	go func(chan os.Signal, chan error) {

		if s.log != nil {
			s.log.Printf("Server starting... PORT = %s", port)
		}

		// err := errors.New("http: Server closed")
		err := s.HttpServer.ListenAndServe() // Бесконечный цикл прослушивания входящих запросов

		if err != nil {
			if err.Error() == "http: Server closed" {
				ch_err <- nil
			} else {
				str := fmt.Sprintf("error occured while running http server: %s", err.Error())
				ch_err <- errors.New(str)
			}
			ch_quit <- syscall.SIGABRT
		}

	}(ch_quit, ch_err)

	// Даём время запуститься серверу
	time.Sleep(2000 * time.Millisecond)

	// Не вернул ли сервер ошибок при запуске
	// Неблокирующее чтение канала, и если сервер не стартанул - возвращаем ошибку
	if err, e := readChanNoBlock(ch_err); e == nil {
		return err
	}

	if s.log != nil {
		s.log.Printf("Server started. PORT = %s", port)
	}

	// Чтение os.Signal из канала. Будет блокировать выполнение главной горутины приложения.
	<-ch_quit

	if err, e := readChanNoBlock(ch_err); e == nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown() error {
	return s.HttpServer.Shutdown(context.Background())
}
