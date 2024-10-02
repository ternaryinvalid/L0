package app

import (
	"L0/config"
	"L0/internal/cache"
	"L0/internal/controller"
	"L0/internal/database"
	"L0/internal/generate"
	"L0/internal/kafka"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		port: ":8080",
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
	conn, err := database.Connect(&cfg.DB)
	if err != nil {
		log.Fatalf("error connecting to PostgreSQL: %v", err)
	}
	defer conn.Close()
	db := database.NewDB(conn)
	// Создаем таблицу, если ее нет
	err = db.CreateTable()
	if err != nil {
		log.Fatalf("error creating table: %v", err)
	}
	// Инициализируем кэш
	cache := cache.NewCache(db)
	cache.Preload()
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
			cache.AddCache(*order) // Добавляем заказ в кэш
			time.Sleep(30 * time.Second)
		}
	}()
	// Создаем HTTP-сервер
	httpServer := NewServer()
	apiController := controller.NewOrderController(cache)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		err := httpServer.ListenAndServe(apiController.GetOrder, apiController.GetAllOrder)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Print("Graceful shutdown")

	// Корректное завершение работы сервера
	if err := httpServer.echo.Shutdown(ctx); err != nil {
		httpServer.echo.Logger.Fatal(err)
	}

}
