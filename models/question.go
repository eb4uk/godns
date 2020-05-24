package models

type Question struct {
	Name  string
	Type  string
	Class string
}

func (q *Question) String() string {
	return q.Name + " " + q.Class + " " + q.Type
}
