package entitis

//Сущности для реализации обработчиков запросов

type Account struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"`
	Id      int    `json:"id"`
}

type Friends struct {
	TargetId int `json:"target_id"`
}

type NewAge struct {
	NewAge int `json:"age"`
}

type Counter struct {
	Id           string `json:"_id"`
	CounterValue int    `json:"counterValue"`
}
