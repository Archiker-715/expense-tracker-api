package entity

type ID struct {
	ID int `json:"id"`
}

type AvailableExpCategories struct {
	AvailableCategories []string
}

func (a *AvailableExpCategories) SetAvailableExpCategories() {
	a.AvailableCategories = make([]string, 0, 7)
	a.AvailableCategories = append(a.AvailableCategories, "groceries", "leisure", "electronics", "utilities", "clothing", "health", "others")
}
