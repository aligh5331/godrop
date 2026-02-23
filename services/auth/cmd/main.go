package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	authServer "github.com/aligh5331/godrop/services/auth/grpc"
	"github.com/aligh5331/godrop/services/auth/grpc/auth"
	"github.com/aligh5331/godrop/services/auth/internal/application/usecase"
	"github.com/aligh5331/godrop/services/auth/internal/domain/repository"
	cache "github.com/aligh5331/godrop/services/auth/internal/infrastructure/cache/redis"
	"github.com/aligh5331/godrop/services/auth/internal/infrastructure/hellpers"
	orm "github.com/aligh5331/godrop/services/auth/internal/infrastructure/persistence/gorm"
	"github.com/aligh5331/godrop/services/auth/internal/infrastructure/security"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	db := SetUpDB()
	userRepo := SetUpUserRepository(db)
	sessionRepo := SetUpSessionRepository(db)
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	rdis := SetUpCache(rdb)

	sessionUseCase := usecase.NewSessionUseCase(
		15*time.Minute, 7*24*time.Hour,
		sessionRepo,
		rdis,
		&security.GoogleUUIDGen{},
		&security.SHA256TokenHasher{},
		security.NewJWTTokenGenerator([]byte("ssss")),
		hellpers.Validator{},
	)

	authUseCAse := usecase.NewAuthUseCase(
		sessionUseCase,
		userRepo,
		security.NewHasher(14),
		&security.GoogleUUIDGen{},
		hellpers.Validator{},
		&hellpers.FakeEmailVerifier{},
	)
	server := authServer.NewAuthServer(authUseCAse, sessionUseCase)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, server)
	log.Println("server is running on :50051")
	if err = grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func SetUpUserRepository(db *gorm.DB) repository.UserRepository {
	return orm.NewUserRepository(db)
}
func SetUpSessionRepository(db *gorm.DB) repository.SessionRepository {
	return orm.NewSessionRepository(db)
}
func SetUpCache(client *redis.Client) repository.CacheRepository {
	return cache.New(client)
}
func SetUpDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Tehran",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
