package tasks

// Task - структура задачи
type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
