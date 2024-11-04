package main

import (
	"log/slog"
	"os"
	"os/signal"
	"product-wallet/config"
	"product-wallet/internal/delivery/http"
	"product-wallet/internal/delivery/http/route"
	"product-wallet/internal/repository"
	services "product-wallet/internal/services"
	"product-wallet/migration"
	"product-wallet/pkg/database"
	"product-wallet/pkg/logger"
	"product-wallet/pkg/server"
	"product-wallet/pkg/signature"
	"product-wallet/pkg/xvalidator"
	"strconv"
	"syscall"
)

var (
	sqlClient *database.Database
)

// @title           Pigeon
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/notificationsvc/api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	validate, _ := xvalidator.NewValidator()
	conf := config.InitAppConfig(validate)
	logger.SetupLogger(&logger.Config{
		AppENV:  conf.AppEnvConfig.AppEnv,
		LogPath: conf.AppEnvConfig.LogFilePath,
		Debug:   conf.AppEnvConfig.AppDebug,
	})
	initInfrastructure(conf)
	ginServer := server.NewGinServer(&server.GinConfig{
		HttpPort:     conf.AppEnvConfig.HttpPort,
		AllowOrigins: conf.AppEnvConfig.AllowOrigins,
		AllowMethods: conf.AppEnvConfig.AllowMethods,
		AllowHeaders: conf.AppEnvConfig.AllowHeaders,
	})
	//external
	signaturer := signature.NewSignature(conf.AuthConfig.JwtSecretAccessToken, conf.AuthConfig.HMACSecretAccessToken)

	// repository
	userRepository := repository.NewUserSQLRepository()
	walletRepository := repository.NewWalletSQLRepository()
	transactionRepository := repository.NewTransactionSQLRepository()

	// service
	userService := services.NewUserService(sqlClient.GetDB(), userRepository, signaturer, validate)
	walletService := services.NewWalletService(sqlClient.GetDB(), walletRepository, userRepository, transactionRepository, validate)
	transactionService := services.NewTransactionService(sqlClient.GetDB(), transactionRepository, walletRepository, validate)
	// Handler
	userHandler := http.NewUserHTTPHandler(userService)
	walletHandler := http.NewWalletHTTPHandler(walletService)
	transactionHandler := http.NewTransactionHTTPHandler(transactionService)

	router := route.Router{
		App:                ginServer.App,
		UserHandler:        userHandler,
		WalletHandler:      walletHandler,
		TransactionHandler: transactionHandler,
	}
	router.Setup()
	router.SwaggerRouter()
	echan := make(chan error)
	go func() {
		echan <- ginServer.Start()
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		slog.Info("signal terminated detected")
	case err := <-echan:
		slog.Error("Failed to start http server", err)
	}
}

func initInfrastructure(config *config.Config) {
	//initPostgreSQL()
	sqlClient = initSQL(config)

}

func initSQL(conf *config.Config) *database.Database {
	db := database.NewDatabase(conf.DatabaseConfig.Dbservice, &database.Config{
		DbHost:   conf.DatabaseConfig.Dbhost,
		DbUser:   conf.DatabaseConfig.Dbuser,
		DbPass:   conf.DatabaseConfig.Dbpassword,
		DbName:   conf.DatabaseConfig.Dbname,
		DbPort:   strconv.Itoa(conf.DatabaseConfig.Dbport),
		DbPrefix: conf.DatabaseConfig.DbPrefix,
	})
	if conf.IsStaging() {
		migration.AutoMigration(db)
	}
	return db
}
