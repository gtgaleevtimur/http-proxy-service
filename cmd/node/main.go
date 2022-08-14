package main

import (
	"module31/internal/app"
	"os"
)

//Запуск ядра приложения по умолчанию слушающего порт :8080
//Для запуска второй реплики приложения в аргументах по умолчанию писать :8081

func main() {
	var addr string
	addr = ":8080"

	if len(os.Args[:]) > 1 {
		addr = os.Args[1]
	}
	app.Run(addr)
}
