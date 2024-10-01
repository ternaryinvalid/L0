package models

type OrderJSON struct {
	OrderUid          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payments          Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_siganture"`
	CustomerId        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	Shardkey          string   `json:"shardkey"`
	Sm_id             uint32   `json:"sm_id"`
	Date_created      string   `json:"date_created"`
	OOF_shard         string   `json:"oof_shard"`
}

type Delivery struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Zip    string `json:"zip"`
	City   string `json:"city"`
	Adress string `json:"address"`
	Region string `json:"region"`
	Email  string `json:"email"`
}

type Payment struct {
	Transaction  string  `json:"transcation"`
	Request_id   string  `json:"request_id"`
	Currency     string  `json:"currency"`
	Provider     string  `json:"provider"`
	Amount       float32 `json:"amount"`
	PaymentDt    uint32  `json:"payment_dt"`
	Bank         string  `json:"bank"`
	DeliveryCost uint32  `json:"delivery_cost"`
	GoodsTotal   float32 `json:"goods_total"`
	CustomFee    float32 `json:"custom_fee"`
}

type Item struct {
	ChrtId       uint32  `json:"chrt_id"`
	Track_number string  `json:"track_number"`
	Price        uint16  `json:"price"`
	Rid          string  `json:"rid"`
	Name         string  `json:"name"`
	Sale         uint16  `json:"sale"`
	Size         string  `json:"size"`
	TotalPrice   float32 `json:"total_price"`
	NmId         uint32  `json:"nm_id"`
	Brand        string  `json:"brand"`
	Status       uint16  `json:"status"`
}
