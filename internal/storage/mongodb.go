package storage

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"module31/internal/entitis"
)

//Реализация создания и управления хранилищем.

type MongoStorage struct {
	users *mongo.Database
}

//Создание указателя Database в MongoDB.

func NewCollection(client *mongo.Client) *MongoStorage {
	return &MongoStorage{
		users: client.Database("Persons"),
	}
}

//Функция для создания коллекции счетчика пользователей для присвоения порядкового _id в MongoDB.

func CounterPersonsCol(client *mongo.Client) {
	counter := bson.D{
		{"_id", "counterForPersons"},
		{"counterValue", 1},
	}

	collection := client.Database("Persons").Collection("Counter")
	_, err := collection.InsertOne(context.TODO(), counter)
	if err != nil {
		log.Println("Обнаружен дубликат автоинкременции id в Persons.Counter")
	}
}

//Метод создания нового пользовтеля внутри обработчика.

func (r *MongoStorage) NewUser(user *entitis.Account) (int, error) {
	var id int
	result := entitis.Counter{}

	filter := bson.D{{"_id", "counterForPersons"}}

	err := r.users.Collection("Counter").FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return 0, errors.New("NotImplementedException")
	}
	id = result.CounterValue

	update := bson.D{
		{"$inc", bson.D{
			{"counterValue", 1},
		}},
	}

	_, err = r.users.Collection("Counter").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return 0, errors.New("NotImplementedException")
	}

	newUser := bson.D{
		{"_id", id},
		{"name", user.Name},
		{"age", user.Age},
		{"friends", []int{}},
		{"flag", "on"},
	}
	_, err = r.users.Collection("Persons").InsertOne(context.TODO(), newUser)
	if err != nil {
		return 0, err
	}

	return id, nil
}

//Метод создания дружбы двух пользовтелей внутри обработчика.

func (r *MongoStorage) MakeFriends(source, target int) (sourceName, targetName string, err error) {
	result := entitis.Account{}

	firstFilter := bson.D{{"_id", source}}
	secondFilter := bson.D{{"_id", target}}

	err = r.users.Collection("Persons").FindOne(context.TODO(), firstFilter).Decode(&result)
	if err != nil {
		return "", "", errors.New("NotImplementedException")
	}
	firstName := result.Name

	err = r.users.Collection("Persons").FindOne(context.TODO(), secondFilter).Decode(&result)
	if err != nil {
		return "", "", errors.New("NotImplementedException")
	}
	secondName := result.Name

	firstUpdate := bson.D{
		{"$push", bson.D{
			{"friends", target},
		}},
	}
	_, err = r.users.Collection("Persons").UpdateOne(context.TODO(), firstFilter, firstUpdate)
	if err != nil {
		return "", "", errors.New("NotImplementedException")
	}

	secondUpdate := bson.D{
		{"$push", bson.D{
			{"friends", source},
		}},
	}

	_, err = r.users.Collection("Persons").UpdateOne(context.TODO(), secondFilter, secondUpdate)
	if err != nil {
		return "", "", errors.New("NotImplementedException")
	}

	return firstName, secondName, nil
}

//Метод создания удаления пользовтеля внутри обработчика.

func (r *MongoStorage) DeleteUser(source int) (string, error) {
	result := entitis.Account{}

	firstFilter := bson.D{{"_id", source}}

	err := r.users.Collection("Persons").FindOne(context.TODO(), firstFilter).Decode(&result)
	if err != nil {
		return "", errors.New("NotImplementedException")
	}
	name := result.Name

	_, err = r.users.Collection("Persons").DeleteOne(context.TODO(), firstFilter)

	firstUpdate := bson.D{
		{"$pull", bson.D{
			{"friends", source},
		}},
	}
	secondFilter := bson.D{{"flag", "on"}}

	_, err = r.users.Collection("Persons").UpdateMany(context.TODO(), secondFilter, firstUpdate)
	if err != nil {
		return "", errors.New("NotImplementedException")
	}

	return name, nil
}

//Метод возврата списка друзей пользователя внутри обработчика.

func (r *MongoStorage) GetFriends(source int) (string, error) {
	result := entitis.Account{}
	firstFilter := bson.D{{"_id", source}}
	err := r.users.Collection("Persons").FindOne(context.TODO(), firstFilter).Decode(&result)
	if err != nil {
		return "", errors.New("NotImplementedException")
	}
	sourceName := result.Name
	friendsList := result.Friends
	res := ""

	for _, value := range friendsList {
		arrayFilter := bson.D{{"_id", value}}
		err = r.users.Collection("Persons").FindOne(context.TODO(), arrayFilter).Decode(&result)
		if err != nil {
			return "", errors.New("NotImplementedException")
		}
		name := result.Name
		res += name + ","
	}
	answer := fmt.Sprintf("Друзья %v : %v", sourceName, res)
	return answer, nil
}

//Метод обновления возраста пользовтеля внутри обработчика.

func (r *MongoStorage) RefreshAge(source, age int) (int, error) {
	result := entitis.Account{}
	firstFilter := bson.D{{"_id", source}}
	err := r.users.Collection("Persons").FindOne(context.TODO(), firstFilter).Decode(&result)
	if err != nil {
		return 0, errors.New("NotImplementedException")
	}

	firstUpdate := bson.D{
		{"$set", bson.D{
			{"age", age},
		}},
	}
	_, err = r.users.Collection("Persons").UpdateOne(context.TODO(), firstFilter, firstUpdate)
	if err != nil {
		return 0, errors.New("NotImplementedException")
	}

	return age, nil
}
