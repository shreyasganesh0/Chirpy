module github.com/shreyasganesh0/Chirpy

go 1.23.4

replace (
    github.com/shreyasganesh0/chirpy/database v0.0.0 => ./internal/database
    )

require (
    github.com/shreyasganesh0/chirpy/database v0.0.0
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/dotenv v2.2.0+incompatible // indirect
	github.com/lib/pq v1.10.9 // indirect
)
