package dom

type TestItem struct {
	Id   int    `json:"id"` //`json:"-"`
	Text string `json:"text"`
}
