package config

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// JWTProvider stores pair iss and aud
type JWTProvider struct {
	Issuer   string `mapstructure:"iss"`
	JWKURL   string `mapstructure:"jwks_url"`
	Audience string `mapstructure:"aud"`
}

// Config stores app config
type Config struct {
	Port                  int
	BackendURL            *url.URL
	JWTProviders          []JWTProvider
	UnauthenticatedRoutes []*regexp.Regexp
}

// Init - initializes configuration
func Init() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	if configEnvPrefix, ok := os.LookupEnv("AIRBAG_CONFIG_ENV_PREFIX"); ok {
		viper.SetEnvPrefix(configEnvPrefix)
	}

	if configName, ok := os.LookupEnv("AIRBAG_CONFIG_NAME"); ok {
		viper.SetConfigName(configName)
	}

	if configPath, ok := os.LookupEnv("AIRBAG_CONFIG_PATH"); ok {
		viper.AddConfigPath(configPath)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// BackendHostName
	viper.SetDefault("BackendHostName", "localhost")
	viper.BindEnv("BackendHostName", "AIRBAG_BACKEND_HOST_NAME")
	backendHostName := viper.GetString("BackendHostName")

	// BackendServicePort
	viper.SetDefault("BackendServicePort", 80)
	viper.BindEnv("BackendServicePort", "AIRBAG_BACKEND_SERVICE_PORT")
	backendServicePort := viper.GetInt("BackendServicePort")

	backendURL, err := url.Parse(fmt.Sprintf("http://%s:%v", backendHostName, backendServicePort))
	if err != nil {
		return nil, fmt.Errorf("Error while parsing backend url: %v", err)
	}

	viper.BindEnv("UnauthenticatedRoutes", "AIRBAG_UNAUTHENTICATED_ROUTES")
	unAuthPathStrs := viper.GetStringSlice("UnauthenticatedRoutes")
	unAuthPaths := make([]*regexp.Regexp, len(unAuthPathStrs))
	for i, str := range unAuthPathStrs {
		unAuthPaths[i] = regexp.MustCompile(str)
	}

	var providers []JWTProvider
	err = viper.UnmarshalKey("JWTProviders", &providers)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing JWTProviders : %v", err)
	}

	port := 80
	if customPort, ok := os.LookupEnv("AIRBAG_PORT"); ok {
		if port, err = strconv.Atoi(customPort); err != nil {
			return nil, fmt.Errorf("Error while getting port from AIRBAG_PORT env variable: %v", err)
		}
	}

	cfg := &Config{
		Port:                  port,
		BackendURL:            backendURL,
		JWTProviders:          providers,
		UnauthenticatedRoutes: unAuthPaths,
	}
	return cfg, nil
}
