package generate

import (
	"L0/internal/models"
	"crypto/md5"
	"encoding/hex"
	"math"
	"math/rand"
	"strconv"
	"time"
)

// GenerateOrder создает новый заказ
func GenerateOrder() *models.OrderJSON {
	orderId := hash32()
	var orderCount = 1 + rand.Intn(2)
	items := make([]models.Item, orderCount)

	for i := 0; i < orderCount; i++ {
		items[i] = models.Item{
			Chrt_id:      0,
			Track_number: "",
			Price:        0,
			Rid:          "",
			Name:         "",
			Sale:         0,
			Size:         "",
			Total_price:  0,
			Nm_id:        0,
			Brand:        "",
			Status:       0,
		}
	}

	order := &models.OrderJSON{
		Order_uid:          orderId[:len(orderId)-15],
		Track_number:       "SOMETRACK",
		Entry:              "SOMEIL",
		Delivery:           models.Delivery{Name: "test", Phone: "+79999999999", Zip: "0", City: "moscow", Adress: "adress", Region: "msk", Email: "email"},
		Payments:           models.Payment{Transaction: "", Request_id: "", Currency: "", Provider: "", Amount: 0, Payment_dt: 0, Bank: "", Delivery_cost: 0, Goods_total: 0, Custom_fee: 0},
		Items:              items,
		Locale:             "en",
		Internal_signature: "",
		Customer_id:        "test",
		Delivery_service:   "some service",
		Shardkey:           "9",
		Sm_id:              0,
		Date_created:       time.Now().Format(time.RFC3339),
		OOF_shard:          "0",
	}

	generateOrderItems(order)
	generateOrderDelivery(order)
	generateOrderPayment(order)

	return order
}

func generateOrderPayment(order *models.OrderJSON) {
	currency := []string{"USD", "RUB", "EUR"}
	banks := []string{"sber", "alpha", "tinkoff"}
	var amount float32 = 0

	for _, item := range order.Items {
		amount += item.Total_price
	}

	deliveryCost := float32(rand.Intn(1500))
	order.Payments.Transaction = order.Order_uid + order.Customer_id
	order.Payments.Currency = currency[rand.Intn(len(currency))]
	order.Payments.Provider = "wbpay"
	order.Payments.Amount = amount + deliveryCost
	order.Payments.Payment_dt = uint32(1000000000 + rand.Intn(1000000000))
	order.Payments.Bank = banks[rand.Intn(len(banks))]
	order.Payments.Delivery_cost = uint32(deliveryCost)
	order.Payments.Goods_total = amount
	order.Payments.Custom_fee = 0
}

func generateOrderDelivery(order *models.OrderJSON) {
	names := []string{"Donald Trump", "Skufotka massonabornay", "Albus Dumbeldore", "Elon Musk"}
	addresses := []string{"Orshanskay 3", "Dedstreet -101", "Chelkastreet 8", "Lenina 11"}

	order.Delivery.Name = names[rand.Intn(len(names))]
	order.Delivery.Phone = "+" + strconv.Itoa(1000000000+rand.Intn(8000000000))
	order.Delivery.Zip = strconv.Itoa(100000 + rand.Intn(150000))
	order.Delivery.City = "Moscow"
	order.Delivery.Adress = addresses[rand.Intn(len(addresses))]
	order.Delivery.Region = "Moscow"
	order.Delivery.Email = "example@gmail.com"

}

func generateOrderItems(order *models.OrderJSON) {
	for i := 0; i < len(order.Items); i++ {
		amount := float32(1500 + rand.Intn(10000))
		sale := float32(rand.Intn(50))
		total := ((100 - sale) / 100.0) * amount
		totalPrice := math.Round(float64(total)*10) / 10
		order.Items[i] = models.Item{
			Chrt_id:      uint32(rand.Intn(1000000)),
			Track_number: order.Track_number,
			Price:        uint16(amount),
			Rid:          hash32()[:len(hash32())-15] + order.Customer_id,
			Name:         "Bread",
			Sale:         uint16(sale),
			Size:         "0",
			Total_price:  float32(totalPrice),
			Nm_id:        uint32(rand.Intn(1000000)),
			Brand:        "Lenta",
			Status:       202,
		}
	}
}

// Генерация хэша
func hash32() string {
	sum := md5.Sum([]byte(strconv.Itoa(rand.Intn(150000))))
	return hex.EncodeToString(sum[:])
}
