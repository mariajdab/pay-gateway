package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mariajdab/pay-gateway/internal/cards"
	cardsDBStorage "github.com/mariajdab/pay-gateway/internal/cards/repository/postgres"
	cardUseCase "github.com/mariajdab/pay-gateway/internal/cards/usecases"
	"github.com/mariajdab/pay-gateway/internal/customers"
	custDBStorage "github.com/mariajdab/pay-gateway/internal/customers/repository/postgres"
	customerUseCase "github.com/mariajdab/pay-gateway/internal/customers/usecase"
	"github.com/mariajdab/pay-gateway/internal/merchant"
	merchantDBStorage "github.com/mariajdab/pay-gateway/internal/merchant/repository/postgres"
	merchantUseCase "github.com/mariajdab/pay-gateway/internal/merchant/usecase"
	"github.com/mariajdab/pay-gateway/internal/payments"
	paymentHttp "github.com/mariajdab/pay-gateway/internal/payments/delivery"
	paymentDBStorage "github.com/mariajdab/pay-gateway/internal/payments/repository/postgres"
	paymentUseCase "github.com/mariajdab/pay-gateway/internal/payments/usecase"
)

type App struct {
	server     *http.Server
	merchantUC merchant.UseCaseMerchant
	customerUC customers.UseCaseCustomer
	paymentUC  payments.UseCasePayment
	cardUC     cards.UseCaseCard
}

func NewApp() *App {
	dbStr, exists := os.LookupEnv("DB_SOURCE")
	if !exists {
		log.Fatalf("Failed to read database connection")
	}

	dbpool, err := pgxpool.New(context.Background(), dbStr)
	if err != nil {
		log.Fatalf("Unable to connect to db: %v", err)
	}

	merchantRepo := merchantDBStorage.NewMerchantRepo(dbpool)
	merchantUC := merchantUseCase.NewMerchantUC(merchantRepo)

	custRepo := custDBStorage.NewCustomerRepo(dbpool)
	customerUC := customerUseCase.NewcustomerUC(custRepo)

	cardRepo := cardsDBStorage.NewCardRepo(dbpool)
	cardUC := cardUseCase.NewCardUC(cardRepo)

	paymentRepo := paymentDBStorage.NewPaymentRepo(dbpool)
	paymentUC := paymentUseCase.NewPaymentUC(paymentRepo, cardRepo, merchantRepo)

	return &App{
		merchantUC: merchantUC,
		customerUC: customerUC,
		cardUC:     cardUC,
		paymentUC:  paymentUC,
	}
}

func (app *App) Run(port string) error {
	router := gin.Default()

	router.Use(gin.Recovery())

	paymentHttp.RegisterHandler(router, app.paymentUC, app.cardUC, app.merchantUC, app.customerUC)

	app.server = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return app.server.ListenAndServe()
}
