package config

import (
	"fmt"
	"microservice/internal/cron"
	"strings"

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

func GetKafkaBrokers() []string {
	return strings.Split(viper.GetString("connections.kafka.brokers"), ",")
}

func GetKafkaConsumerGroupName() string {
	return viper.GetString("connections.kafka.group_name")
}

func GetCronJobs() cron.JobConfig {
	jobs := make(cron.JobConfig)
	for name, schedule := range viper.GetStringMapString("cron") {
		jobs[cron.JobName(name)] = cron.JobSchedule(schedule)
	}

	return jobs
}