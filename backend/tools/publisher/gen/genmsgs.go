package gen

import (
	"encoding/json"
	"fmt"

	"github.com/go-faker/faker/v4"
	"github.com/taaanechka/order-service/internal/api-server/services/ports/ordersrepository"
)

func GeneratePositiveData() ([]byte, error) {
	var order ordersrepository.Order
	_ = faker.FakeData(&order)
	fmt.Printf("POS: %+v\n\n", order)

	b, err := json.Marshal(&order)
	if err != nil {
		return nil, fmt.Errorf("GEN+: failed to marshal order")
	}
	return b, nil
}

func GenerateNegativeData() ([]byte, error) {
	var order ordersrepository.Order
	_ = faker.FakeData(&order.OrderUid)
	fmt.Printf("NEG: %+v\n\n", order)

	b, err := json.Marshal(&order)
	if err != nil {
		return nil, fmt.Errorf("GEN-: failed to marshal order")
	}
	return b, nil
}
