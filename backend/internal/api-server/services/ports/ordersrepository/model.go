package ordersrepository

type Order struct {
	OrderUid          string        `json:"order_uid,omitempty"`
	TrackNumber       string        `json:"track_number"`
	Entry             string        `json:"entry"`
	Delivery          OrderDelivery `json:"delivery"`
	Payment           OrderPayment  `json:"payment"`
	Items             []Item        `json:"items"`
	Locale            string        `json:"locale"`
	InternalSignature string        `json:"internal_signature,omitempty"`
	CustomerId        string        `json:"customer_id"`
	Meest             string        `json:"meest"`
	Shardkey          string        `json:"shardkey"`
	SmId              uint          `json:"sm_id"`
	DateCreated       string        `json:"date_created"`
	OofShard          string        `json:"oof_shard"`
}

type OrderDelivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type OrderPayment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id,omitempty"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       uint   `json:"amount"`
	PaymentDt    uint   `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost uint   `json:"delivery_cost"`
	GoodsTotal   uint   `json:"goods_total"`
	CustomFee    uint   `json:"custom_fee"`
}

type Item struct {
	ChrtId      uint   `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       uint   `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        uint   `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  uint   `json:"total_price"`
	NmId        uint   `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      uint   `json:"status"`
}
