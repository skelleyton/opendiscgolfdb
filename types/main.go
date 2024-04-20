package types

type Config struct {
	DbUser     string `dotenv:"DB_USER"`
	DbPassword string `dotenv:"DB_PASSWORD"`
	ConnStr    string `dotenv:"CONN_STR"`
}
