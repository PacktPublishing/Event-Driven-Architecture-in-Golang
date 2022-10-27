package domain

type Items map[string]*Item

type Item struct {
	ProductName string
	Quantity    int
}
