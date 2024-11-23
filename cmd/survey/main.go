package main

import (
	"context"
	"net"
	"time"

	"github.com/go-park-mail-ru/2024_2_BogoSort/config"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/survey"
	surveyProto "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/survey/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_2_BogoSort/pkg/connector"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthProto "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	defer zap.L().Sync()

	cfg, err := config.Init()
	if err != nil {
		zap.L().Fatal("Ошибка при инициализации конфигурации", zap.Error(err))
	}

	dbPool, err := connector.GetPostgresConnector(cfg.GetConnectURL())
	if err != nil {
		zap.L().Error("Failed to connect to Postgres", zap.Error(err))
		return
	}

	answerRepo, err := postgres.NewAnswerRepository(dbPool, zap.L(), context.Background(), time.Duration(cfg.PGTimeout))
	if err != nil {
		zap.L().Error("Failed to create survey repository", zap.Error(err))
		return
	}

	questionRepo, err := postgres.NewQuestionRepository(dbPool, zap.L(), context.Background(), time.Duration(cfg.PGTimeout))
	if err != nil {
		zap.L().Error("Failed to create question repository", zap.Error(err))
		return
	}

	server := grpc.NewServer()
	surveyUC := service.NewAnswerService(answerRepo, zap.L())
	surveyServer := survey.NewSurveyGrpcServer(surveyUC, questionRepo)

	healthServer := health.NewServer()
	healthProto.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", healthProto.HealthCheckResponse_SERVING)

	surveyProto.RegisterSurveyServiceServer(server, surveyServer)
	address := config.GetSurveyAddress()
	lis, err := net.Listen("tcp", address)
	if err != nil {
		zap.L().Fatal("Failed to listen on port", zap.Error(err))
	}

	zap.L().Info("Survey server started on " + address)

	if err := server.Serve(lis); err != nil {
		zap.L().Fatal("Failed to start gRPC server", zap.Error(err))
	}
}
