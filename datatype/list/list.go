package list

type List interface {
	Push(index int, value interface{})
	Remove(index int)
}
