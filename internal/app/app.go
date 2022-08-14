package app

import (
	"context"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"module31/internal/handler"
	"module31/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//Функция сборки и запуска приложения

func Run(addr string) {
	//Инициализация клиента для управления MongoDB
	ctx := context.TODO()
	opts := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	server := &http.Server{Addr: addr, Handler: service(client)}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	//Реализация graceful-shutdown приложения.

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-sigs
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-serverCtx.Done()
}

//Инициализация хранилища MongoDB,и сборка хэндлера для ядра приложения.
func service(client *mongo.Client) http.Handler {
	newStorage := storage.NewCollection(client)
	storage.CounterPersonsCol(client)
	router := chi.NewRouter()
	handler.Build(router, newStorage)

	return router
}
