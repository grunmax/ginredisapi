package dom

type TodoItem struct {
	Id        string `					redis:"id" 			json:"id"` //`json:"-"`
	Title     string `validate:"nonzero"	redis:"title" 		json:"title"`
	Completed bool   `					redis:"completed"	json:"completed"`
	Order     int    `validate:"min=21" 	redis:"order"		json:"order"`
}
