package config

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.SetConfigName("microservice")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}

	return nil
}

func GetDB() string {
	mysqlconf := mysql.Config{
		User:      viper.GetString("connections.db.user"),
		Passwd:    viper.GetString("connections.db.password"),
		Addr:      viper.GetString("connections.db.addr"),
		DBName:    viper.GetString("connections.db.name"),
		ParseTime: true,
	}

	return mysqlconf.FormatDSN()
}
