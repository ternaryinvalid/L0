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

// создает новый заказ
func GetOrder() *models.OrderJSON {
	orderId := hash32()
	var orderCount = 1 + rand.Intn(2)
	items := make([]models.Item, orderCount)

	for i := 0; i < orderCount; i++ {
		items[i] = models.Item{
			ChrtId:       0,
			Track_number: "",
			Price:        0,
			Rid:          "",
			Name:         "",
			Sale:         0,
			Size:         "",
			TotalPrice:   0,
			NmId:         0,
			Brand:        "",
			Status:       0,
		}
	}

	order := &models.OrderJSON{
		OrderUid:          orderId[:len(orderId)-15],
		TrackNumber:       "SOMETRACK",
		Entry:             "SOMEIL",
		Delivery:          models.Delivery{Name: "test", Phone: "+79999999999", Zip: "0", City: "moscow", Adress: "adress", Region: "msk", Email: "email"},
		Payments:          models.Payment{Transaction: "", Request_id: "", Currency: "", Provider: "", Amount: 0, PaymentDt: 0, Bank: "", DeliveryCost: 0, GoodsTotal: 0, CustomFee: 0},
		Items:             items,
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "some service",
		Shardkey:          "9",
		Sm_id:             0,
		Date_created:      time.Now().Format(time.RFC3339),
		OOF_shard:         "0",
	}

	getOrderItems(order)
	getOrderDelivery(order)
	getOrderPayment(order)

	return order
}

func getOrderPayment(order *models.OrderJSON) {
	currency := []string{"USD", "RUB", "EUR"}
	banks := []string{"sber", "alpha", "tinkoff"}
	var amount float32 = 0

	for _, item := range order.Items {
		amount += item.TotalPrice
	}

	deliveryCost := float32(rand.Intn(1500))
	order.Payments.Transaction = order.OrderUid + order.CustomerId
	order.Payments.Currency = currency[rand.Intn(len(currency))]
	order.Payments.Provider = "wbpay"
	order.Payments.Amount = amount + deliveryCost
	order.Payments.PaymentDt = uint32(1000000000 + rand.Intn(1000000000))
	order.Payments.Bank = banks[rand.Intn(len(banks))]
	order.Payments.DeliveryCost = uint32(deliveryCost)
	order.Payments.GoodsTotal = amount
	order.Payments.CustomFee = 0
}

func getOrderDelivery(order *models.OrderJSON) {
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

func getOrderItems(order *models.OrderJSON) {
	for i := 0; i < len(order.Items); i++ {
		amount := float32(1500 + rand.Intn(10000))
		sale := float32(rand.Intn(50))
		total := ((100 - sale) / 100.0) * amount
		totalPrice := math.Round(float64(total)*10) / 10
		order.Items[i] = models.Item{
			ChrtId:       uint32(rand.Intn(1000000)),
			Track_number: order.TrackNumber,
			Price:        uint16(amount),
			Rid:          hash32()[:len(hash32())-15] + order.CustomerId,
			Name:         "Bread",
			Sale:         uint16(sale),
			Size:         "0",
			TotalPrice:   float32(totalPrice),
			NmId:         uint32(rand.Intn(1000000)),
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
