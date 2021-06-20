module dbforum

go 1.16

require (
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/fasthttp/router v1.3.14
	github.com/jackc/pgtype v1.7.0
	github.com/jackc/pgx/v4 v4.11.0
	github.com/joho/godotenv v1.3.0
	github.com/klauspost/compress v1.13.0 // indirect
	github.com/mailru/easyjson v0.7.7
	github.com/valyala/fasthttp v1.26.0
	go.uber.org/zap v1.13.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.38.0
