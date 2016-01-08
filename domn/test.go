package domn

type TestItem struct {
	Id   int    `json:"id"` //`json:"-"`
	Text string `json:"text"`
}

type FormFile struct {
	Id string `json:"id"`
}
