module github.com/my-crazy-lab/this-is-grpc/authentication

go 1.23.6

replace github.com/my-crazy-lab/this-is-grpc/proto-module => ../proto-module

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jackc/pgx/v5 v5.7.2
	github.com/joho/godotenv v1.5.1
	github.com/my-crazy-lab/this-is-grpc/proto-module v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.32.0
	google.golang.org/grpc v1.70.0
)

require (
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/oauth2 v0.25.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/grpc/examples v0.0.0-20250207091334-e0d191d8adcd // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
