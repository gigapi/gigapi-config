package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
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
	// Layers configuration
	Layers []LayersConfiguration `json:"layers" mapstructure:"layers" default:""`
}

type LayersConfiguration struct {
	// Name of the layer
	Name string `json:"name" mapstructure:"name" default:""`
	// Type of the layer (s3, fs)
	Type string `json:"type" mapstructure:"type" default:""`
	// If the layer is local for writer or global
	Global bool `json:"global" mapstructure:"global" default:"false"`
	// URL of the layer
	//
	//   Example: s3://key:secret@localhost:8181/bucket/prefix - s3 URL
	//   Example: file:///data/folder/path - root path for filesystem
	URL string `json:"url" mapstructure:"url" default:""`
	// How long to keep data before moving to the next layer (empty for unlimited)
	//
	//   Example: 1h - keep data for 1 hour
	//   Example: 10m - keep data for 10 minutes
	TTL time.Duration `json:"ttl" mapstructure:"ttl" default:""`
	// Auth configuration for s3 layers
	Auth LayerAuthConfiguration `json:"auth" mapstructure:"auth" default:""`
}

type LayerAuthConfiguration struct {
	// Key for authentication
	Key string `json:"key" mapstructure:"key" default:""`
	// Secret for authentication
	Secret string `json:"secret" mapstructure:"secret" default:""`
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
	setLayers()
	fmt.Printf("Loaded configuration: %+v\n", Config)
}

func setLayers() {
	for i := 0; ; i += 1 {
		if os.Getenv(fmt.Sprintf("GIGAPI_LAYERS_%d_NAME", i)) == "" {
			break
		}
		var l LayersConfiguration
		l.Name = os.Getenv(fmt.Sprintf("GIGAPI_LAYERS_%d_NAME", i))
		l.Type = os.Getenv(fmt.Sprintf("GIGAPI_LAYERS_%d_TYPE", i))
		l.Global = os.Getenv(fmt.Sprintf("GIGAPI_LAYERS_%d_GLOBAL", i)) == "true"
		l.URL = os.Getenv(fmt.Sprintf("GIGAPI_LAYERS_%d_URL", i))
		l.Auth.Key = os.Getenv(fmt.Sprintf("GIGAPI_LAYERS_%d_AUTH_KEY", i))
		l.Auth.Secret = os.Getenv(fmt.Sprintf("GIGAPI_LAYERS_%d_AUTH_SECRET", i))
		if fmt.Sprintf("GIGAPI_LAYERS_%d_TTL", i) == "" {
			l.TTL = 0
		} else {
			l.TTL, _ = time.ParseDuration(os.Getenv(fmt.Sprintf("GIGAPI_LAYERS_%d_TTL", i)))
		}
		if i < len(Config.Gigapi.Layers) {
			Config.Gigapi.Layers[i] = l
		} else {
			Config.Gigapi.Layers = append(Config.Gigapi.Layers, l)
		}
	}
	if len(Config.Gigapi.Layers) > 0 {
		return
	}
	Config.Gigapi.Layers = []LayersConfiguration{
		{Name: "default", Type: "fs", Global: true, URL: Config.Gigapi.Root},
	}
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
