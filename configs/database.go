package configs

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

type Config struct {
	C *viper.Viper
}

func NewConfig() *Config {
	c := viper.New()
	c.SetConfigFile(".env")
	c.AddConfigPath(".")
	_ = c.ReadInConfig()
	return &Config{C: c}
}

func ConnectDB(config *Config) *sql.DB {
	DbUser := config.C.GetString("DB_USER")
	DbPassword := config.C.GetString("DB_PASSWORD")
	DbName := config.C.GetString("DB_NAME")
	DbHost := config.C.GetString("DB_HOST")
	DbPort := config.C.GetString("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DbUser, DbPassword, DbHost, DbPort, DbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(10 * time.Minute)
	log.Println("Connected to database")
	return db
}
