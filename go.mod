module neotype-backend

// +heroku goVersion go1.17
go 1.17

// +heroku install ./cmd/...

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.8.2
	github.com/go-redis/redis/v8 v8.11.4
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/joho/godotenv v1.4.0
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.8.1
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/datatypes v1.0.4
	gorm.io/driver/mysql v1.2.1
	gorm.io/gorm v1.22.4
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/goccy/go-json v0.9.11 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
