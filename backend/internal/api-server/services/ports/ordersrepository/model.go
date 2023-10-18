package ordersrepository

type Order struct {
	OrderUid          string        `json:"order_uid" validate:"required"`
	TrackNumber       string        `json:"track_number" validate:"required"`
	Entry             string        `json:"entry" validate:"required"`
	Delivery          OrderDelivery `json:"delivery" validate:"required,structonly"`
	Payment           OrderPayment  `json:"payment" validate:"required,structonly"`
	Items             []Item        `json:"items" validate:"required,structonly"`
	Locale            string        `json:"locale" validate:"required"`
	InternalSignature string        `json:"internal_signature,omitempty" validate:"omitempty"`
	CustomerId        string        `json:"customer_id" validate:"required"`
	Meest             string        `json:"meest" validate:"required"`
	Shardkey          string        `json:"shardkey" validate:"required"`
	SmId              uint          `json:"sm_id" validate:"required"`
	DateCreated       string        `json:"date_created" validate:"required"`
	OofShard          string        `json:"oof_shard" validate:"required"`
}

type OrderDelivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required"`
}

type OrderPayment struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestId    string `json:"request_id,omitempty" validate:"omitempty"`
	Currency     string `json:"currency" validate:"required"`
	Provider     string `json:"provider" validate:"required"`
	Amount       uint   `json:"amount" validate:"required"`
	PaymentDt    uint   `json:"payment_dt" validate:"required"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost uint   `json:"delivery_cost" validate:"required"`
	GoodsTotal   uint   `json:"goods_total" validate:"required"`
	CustomFee    uint   `json:"custom_fee" validate:"required"`
}

type Item struct {
	ChrtId      uint   `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       uint   `json:"price" validate:"required"`
	Rid         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        uint   `json:"sale" validate:"required"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  uint   `json:"total_price" validate:"required"`
	NmId        uint   `json:"nm_id" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
	Status      uint   `json:"status" validate:"required"`
}
