module github.com/SinKingCloud/sinking-go/sinking-consul

go 1.11

require (
	github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20211008065154-eb6948aeb02a
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/imroc/req v0.3.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/spf13/viper v1.9.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	gorm.io/driver/mysql v1.1.2 // indirect
	gorm.io/driver/sqlite v1.1.6 // indirect
	gorm.io/gorm v1.21.16 // indirect
)

replace github.com/SinKingCloud/sinking-go/sinking-web v0.0.0-20211008065154-eb6948aeb02a => ../sinking-web
