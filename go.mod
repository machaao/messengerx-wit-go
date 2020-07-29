module github.com/abhishekraj272/golang_machaao

go 1.23

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/wit-ai/wit-go v1.0.12
	machaao-go/extras v0.0.0-00010101000000-000000000000
	machaao-go/machaao v0.0.0-00010101000000-000000000000
)

replace (
	machaao-go/extras => ./extras
	machaao-go/machaao => ./machaao
)
