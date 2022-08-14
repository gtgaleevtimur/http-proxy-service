package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"module31/internal/entitis"
	"module31/internal/storage"
	"net/http"
	"strconv"
)

//Реализация роутера для обработки запросов сервиса через контроллер

func Build(router *chi.Mux, store *storage.MongoStorage) {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	controller := NewController(store)

	router.Post("/users", controller.CreateUser)             //запрос создания пользователя
	router.Put("/users/{id}/friends", controller.MakeFriend) //запрос создания дружбы двух пользователей
	router.Delete("/users/{id}", controller.DeleteUser)      //запрос удаления пользователя
	router.Get("/users/{id}/friends", controller.GetFriends) //запрос возврата всех друзей пользователя
	router.Patch("/users/{id}", controller.UpdateAge)        //запрос сизменения возраста пользователя
}

type Controller struct {
	storage *storage.MongoStorage
}

func NewController(storage *storage.MongoStorage) *Controller {
	return &Controller{
		storage: storage,
	}
}

//Реализация обработчика запроса создания пользователя

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &entitis.Account{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	id, err := c.storage.NewUser(user)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	result := map[string]int{"id": id}
	response, err := json.Marshal(result)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	buildResponse(w, http.StatusCreated, response)
}

//Реализация обработчика запроса создания дружбы двух пользовтелей

func (c *Controller) MakeFriend(w http.ResponseWriter, r *http.Request) {
	friendList := &entitis.Friends{}
	err := json.NewDecoder(r.Body).Decode(friendList)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	id := chi.URLParam(r, "id")
	sourceId, err := strconv.Atoi(id)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}

	sourceName, targetName, err := c.storage.MakeFriends(sourceId, friendList.TargetId)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	res := fmt.Sprintf("%v и %v теперь друзья", sourceName, targetName)
	result := map[string]string{"message": res}
	response, err := json.Marshal(result)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	buildResponse(w, http.StatusOK, response)
}

//Реализация обработчика запроса удаления пользовтеля

func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sourceId, err := strconv.Atoi(id)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}

	res, err := c.storage.DeleteUser(sourceId)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	res = fmt.Sprintf("%v удален.", res)
	result := map[string]string{"message": res}
	response, err := json.Marshal(result)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	buildResponse(w, http.StatusOK, response)
}

//Реализация обработчика запроса списка друзей пользовтеля

func (c *Controller) GetFriends(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sourceId, err := strconv.Atoi(id)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	res, err := c.storage.GetFriends(sourceId)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	result := map[string]string{"message": res}
	response, err := json.Marshal(result)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	buildResponse(w, http.StatusOK, response)
}

//Реализация обработчика запроса обновления возраста пользовтеля

func (c *Controller) UpdateAge(w http.ResponseWriter, r *http.Request) {
	newAge := &entitis.NewAge{}
	err := json.NewDecoder(r.Body).Decode(newAge)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	id := chi.URLParam(r, "id")
	sourceId, err := strconv.Atoi(id)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	res, err := c.storage.RefreshAge(sourceId, newAge.NewAge)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	result := map[string]int{"newAge": res}
	response, err := json.Marshal(result)
	if err != nil {
		buildResponse(w, http.StatusNotFound, nil)
		return
	}
	buildResponse(w, http.StatusOK, response)
}

//Функция формирования ответа обработчика

func buildResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Set("Content-Type", "application/json ; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(body)
}
