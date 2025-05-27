package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"strconv"
	"strings"
)

type GigapiConfiguration struct {
	// Root folder for all the data files
	Root string `json:"root" mapstructure:"root" default:""`
	// Base timeout between merges
	MergeTimeoutS int `json:"merge_timeout_s" mapstructure:"merge_timeout_s" default:"10"`
	// Timeout before saving the new data to the disk
	SaveTimeoutS float64 `json:"save_timeout_s" mapstructure:"save_timeout_s" default:"1"`
	// Disable merging
	NoMerges bool `json:"no_merges" mapstructure:"no_merges" default:"false"`
	// Enable UI for querier
	UI bool `json:"ui" mapstructure:"ui" default:"true"`
	// Execution mode (readonly, writeonly, compaction, aio)
	Mode string `json:"mode" mapstructure:"mode" default:"aio"`
	// Index configuration for the data storage
	Metadata MetadataConfiguration `json:"metadata" mapstructure:"metadata" default:""`
}

type MetadataConfiguration struct {
	// Type of metadata storage (json or redis)
	Type string `json:"type" mapstructure:"type" default:"json"`
	// URL: Redis url in case of Redis metadata storage
	//
	// Example:
	//  - redis://localhost:6379/0 - for no authentication
	//  - redis://username:password@localhost:6379/0 - for password authentication
	//  - rediss://username:password@localhost:6379/0 - for SSL
	URL string `json:"url" mapstructure:"url" default:""`
}

type BasicAuthConfiguration struct {
	Username string `json:"username" mapstructure:"username" default:""`
	Password string `json:"password" mapstructure:"password" default:""`
}

type FlightSqlConfiguration struct {
	// Port to run flightSQL server
	Port int `json:"port" mapstructure:"port" default:"8082"`
	// Enable FlightSQL server
	Enable bool `json:"enable" mapstructure:"enable" default:"true"`
}

type HTTPConfiguration struct {
	// Port to listen on
	Port int `json:"port" mapstructure:"port" default:"7971"`
	// Host to bind to (0.0.0.0 for all interfaces)
	Host string `json:"host" mapstructure:"host" default:"0.0.0.0"`
	// Basic authentication configuration
	BasicAuth BasicAuthConfiguration `json:"basic_auth" mapstructure:"basic_auth" default:""`
}

type Configuration struct {
	Gigapi GigapiConfiguration `json:"gigapi" mapstructure:"gigapi" default:""`
	// HTTP server configuration (reader and writer)
	HTTP HTTPConfiguration `json:"http" mapstructure:"http" default:""`
	// FlightSQL server configuration
	FlightSql FlightSqlConfiguration `json:"flightsql" mapstructure:"flightsql" default:""`
	// Log level (debug, info, warn, error, fatal)
	Loglevel string `json:"loglevel" mapstructure:"loglevel" default:"info"`
}

var Config *Configuration

func InitConfig(file string) {
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	// If a file is provided, use it as the config file
	if file != "" {
		viper.SetConfigFile(file)
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("error reading config file: %s", err))
		}
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Using environment variables for configuration")
	}

	Config = &Configuration{}
	err := viper.Unmarshal(Config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %s", err))
	}
	if Config.Gigapi.SaveTimeoutS == 0 {
		Config.Gigapi.SaveTimeoutS = 1
	}
	setDefaults(Config)
	fmt.Printf("Loaded configuration: %+v\n", Config)
}

func setDefaults(config any) {
	configValue := reflect.ValueOf(config).Elem()
	configType := configValue.Type()

	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Field(i)
		fieldType := configType.Field(i)

		if field.Kind() == reflect.Struct {
			setDefaults(field.Addr().Interface())
			continue
		}

		defaultTag := fieldType.Tag.Get("default")
		if defaultTag == "" {
			continue
		}

		if field.IsZero() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(defaultTag)
			case reflect.Int:
				if intValue, err := strconv.Atoi(defaultTag); err == nil {
					field.SetInt(int64(intValue))
				}
			case reflect.Float64:
				if floatValue, err := strconv.ParseFloat(defaultTag, 64); err == nil {
					field.SetFloat(floatValue)
				}
			case reflect.Bool:
				if boolValue, err := strconv.ParseBool(defaultTag); err == nil {
					field.SetBool(boolValue)
				}
			}
		}
	}
}
