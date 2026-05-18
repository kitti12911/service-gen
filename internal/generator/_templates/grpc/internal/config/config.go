package config

import (
	"time"

	"___LIB_PATH___/lib-monitor/profiling"
	"___LIB_PATH___/lib-monitor/tracing"
	liborm "___LIB_PATH___/lib-orm/v3"
	"___LIB_PATH___/lib-util/v3/logger"
)

type Config struct {
	Service   Service          `mapstructure:"service"   validate:"required"`
	Logging   logger.Config    `mapstructure:"logging"`
	Tracing   tracing.Config   `mapstructure:"tracing"`
	Profiling profiling.Config `mapstructure:"profiling"`
	Database  liborm.Config    `mapstructure:"database"  validate:"required"`
}

type Service struct {
	Name            string        `mapstructure:"name"             env:"SERVICE_NAME"      validate:"required"`
	Port            int           `mapstructure:"port"             env:"PORT"              validate:"required,gte=1,lte=65535"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT"`
}
