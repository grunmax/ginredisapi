package dom

type TodoItem struct {
	Id        string `redis:"id"`        //`json:"id"` //`json:"-"`
	Title     string `redis:"title"`     //`json:"title"`
	Completed bool   `redis:"completed"` //`json:"completed"`
	Order     int    `redis:"order"`     //`json:"order"`
}
