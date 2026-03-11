package logcases

import (
	"log"
	"log/slog"

	"go.uber.org/zap"
)

func main() {
	log.Info("Starting server on port 8080") // want "log message should start with a lowercase letter"
	log.Info("запуск сервера") // want "log message should be in English only"
	log.Info("server started!🚀") // want "log message should not contain special symbols or emoji"
	log.Info("user password: 12345") // want "log message should not contain sensitive data"
	log.Info("starting server on port 8080") // ok
	log.Info("server started") // ok
	log.Info("user authenticated successfully") // ok
	slog.Error("Failed to connect to database") // want "log message should start with a lowercase letter"
	slog.Error("ошибка подключения к базе данных") // want "log message should be in English only"
	slog.Error("connection failed!!!") // want "log message should not contain special symbols or emoji"
	slog.Error("token: abcdef") // want "log message should not contain sensitive data"
	slog.Error("failed to connect to database") // ok
	zap.Sugar().Info("API_KEY=123") // want "log message should not contain sensitive data"
	zap.Sugar().Info("api request completed") // ok
}
