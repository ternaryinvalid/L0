package app

import (
	"L0/config"
	"L0/internal/cache"
	"L0/internal/controller"
	"L0/internal/database"
	"L0/internal/generate"
	"L0/internal/kafka"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
	port string
}

func NewServer() *Server {
	return &Server{
		echo: echo.New(),
		port: ":8000",
	}
}

func (s *Server) ListenAndServe(orderHand echo.HandlerFunc, allOrdersHand echo.HandlerFunc) error {
	s.echo.GET("/orders/:order", orderHand)
	s.echo.GET("/orders/get", allOrdersHand)
	return s.echo.Start(s.port)
}

func Run(cfg *config.Config) {
	// Создаем продюсера Kafka
	producer, err := kafka.NewProducer(&cfg.Kafka)
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}

	// Создаем консьюмера Kafka
	consumer, err := kafka.NewConsumer(&cfg.Kafka)
	if err != nil {
		log.Fatalf("failed to create Kafka consumer: %v", err)
	}

	// Подключаемся к PostgreSQL
	ctx := context.Background()
	conn, err := database.Connect(&cfg.DB, ctx)
	if err != nil {
		log.Fatalf("error connecting to PostgreSQL: %v", err)
	}

	defer conn.Close()
	db := database.NewDB(conn, cfg.DB.Schema)

	// Создаем таблицу и схему, если их нет
	err = db.CreateSchemaAndTable(ctx)
	if err != nil {
		log.Fatalf("error creating table: %v", err)
	}

	// Инициализируем кэш
	cache := cache.NewCache(db)
	cache.Preload(ctx)

	// Запускаем горутину для публикации заказов
	go func() {
		for {
			order := generate.GetOrder()
			fmt.Println("Order is sent")
			err := producer.Publish(*order) // Публикуем заказ
			if err != nil {
				log.Printf("error publishing: %v\n", err)
			}

			time.Sleep(30 * time.Second)
		}
	}()

	// Запускаем горутину для подписки на заказы
	go func() {
		for {
			order, err := consumer.Subscribe() // Подписываемся на заказы
			if err != nil {
				log.Printf("error subscribing: %v\n", err)
				continue
			}

			fmt.Println("Order received")
			cache.AddCache(*order, ctx) // Добавляем заказ в кэш
			time.Sleep(30 * time.Second)
		}
	}()

	// Создаем HTTP-сервер
	httpServer := NewServer()
	apiController := controller.NewOrderController(cache)

	// Запускаем HTTP-сервер
	serverErr := httpServer.ListenAndServe(apiController.GetOrder, apiController.GetAllOrder)
	if serverErr != nil {
		log.Fatalf("error starting server: %v", serverErr)
	}

}
