# HTTP Twitter Clone Server

## Introduction

A http server written in go that allows users to write and read "chirps" of 140 characters or less.
Uses the go url Mux for url parsing.
Postgres is used to database manangment of users and chirps.
JWT token for access and Refresh tokens to create JWT tokens.

## Requirements

- Go version 1.23+
- Postgres version 15+
- Goose for go
- sqlc for go

## Usage
- GET /api/chirps
    - ? query params sort chirps by ascending (asc) and descending (desc)
    - ? query params author_id sort chirps user id

- GET /api/chirps/{chirpID}
- GET /api/healthz 
- GET /admin/metrics 
- POST /admin/reset
- POST /api/users
- POST /api/login
- POST /api/revoke
- POST /api/chirps
- POST /api/polka/webhooks
- PUT /api/users
- DELETE /api/chirps/{chirpID}

