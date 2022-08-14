# httpProxyService

# Golang http.ProxyService
В сервисе реализовано:
- Поднятие двух реплик приложения
- Запуск реверсивной прокси для получения запросов и дилигирования их двум запущенным репликам приложения.
![img.png](img.png)
- Запись,изменение и удаление данных для MongoDB.

## Начало работы

### Требования
Установите компилятор `Go` (если еще не установлен): https://golang.org/doc/tutorial/getting-started
Необходим доступ к рабочей MongoDB: https://docs.mongodb.com/manual/

### Установка
Склонируйте репозиторий ,запустите две реплики node/main.go и запустите прокси proxy/main.go,
сервис будет доступен локально по URL http://localhost:8082.

### Взаимодействие с сервисом
- При запуске node/main.go по умолчанию запускает приложение локально на http://localhost:8080 ,
  для запуска второй реплики приложения передайте в аргументов адрес локального порта ':8081'.
- Далее для работы сервиса запускаете proxy/main.go , прокси слушает запросы локально по адресу http://localhost:8082
  и передает репликам по очереди , возвращает ответы.
- В сервисе реализован обработчик создания пользователя. У пользователя должны быть следующие поля: имя, возраст и
  массив друзей. Пример запроса :

POST http://localhost:8082/users HTTP/1.1
content-type: application/json

{"name":"name","age":age,"friends":[]}

Данный запрос должен возвращать id и статус 201.

-В сервисе реализован обработчик , который делает друзей из двух пользователей. Например, если мы создали двух
пользователей и нам вернулись их ID, то в запросе мы можем указать ID пользователя, который инициировал запрос на
дружбу, и ID пользователя, который примет инициатора в друзья. Пример запроса:

PUT  http://localhost:8082/users/1/friends HTTP/1.1
content-type: application/json

{"target_id": 2}

Данный запрос должен возвращать статус 200 и сообщение «username_1 и username_2 теперь друзья».

-В сервисе реализован обработчик ,который удаляет пользователя. Данный обработчик принимает ID пользователя и
удаляет его из хранилища, а также стирает его из массива friends у всех его друзей. Пример запроса:

DELETE http://localhost:8082/users/1 HTTP/1.1
content-type: application/json

Данный запрос должен возвращать 200 и имя удалённого пользователя с id = 1.

-В сервисе реализован обработчик, который возвращает всех друзей пользователя. Пример запроса:

GET  http://localhost:8082/users/1/friends HTTP/1.1
content-type: application/json

Данный запрос должен возвращать 200 и имя друзей пользователя с id = 1.

-В сервисе реализован обработчик, который обновляет возраст пользователя. Пример запроса:

PATCH http://localhost:8082/users/1 HTTP/1.1
content-type: application/json

{"age": 33}

Запрос должен возвращать 200 и сообщение «возраст «пользователя» успешно обновлён».



### Завершение работы с сервисом
- Для завершения работы нажмите `Ctrl+C` в его консоли (graceful shutdown).