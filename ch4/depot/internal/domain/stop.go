package domain

type Stops map[string]*Stop

type Stop struct {
	StoreName     string
	StoreLocation string
	Items         Items
}

func (s *Stop) AddItem(product *Product, quantity int) error {
	if _, exists := s.Items[product.ID]; !exists {
		s.Items[product.ID] = &Item{
			ProductName: product.Name,
			Quantity:    quantity,
		}
	}

	return nil
}
