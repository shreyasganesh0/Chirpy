module github.com/shreyasganesh0/Chirpy

go 1.23.4

replace (
	github.com/shreyasganesh0/Chirpy/auth v0.0.0 => ./internal/auth
	github.com/shreyasganesh0/Chirpy/database v0.0.0 => ./internal/database
)

require (
	github.com/shreyasganesh0/Chirpy/auth v0.0.0
	github.com/shreyasganesh0/Chirpy/database v0.0.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/dotenv v2.2.0+incompatible // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.32.0 // indirect
)
