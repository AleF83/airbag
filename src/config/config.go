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
	Port                 int
	BackendURL           *url.URL
	JWTProviders         []JWTProvider
	UnauthenticatedPaths []*regexp.Regexp
}

// Init - initializes configuration
func Init() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if configEnvPrefix, ok := os.LookupEnv("AIRBAG_CONFIG_ENV_PREFIX"); ok {
		viper.SetPrefix(configEnvPrefix)
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

	backendURL, err := url.Parse(viper.GetString("backend"))
	if err != nil {
		return nil, fmt.Errorf("Error while parsing backend url: %v", err)
	}

	unAuthPathStrs := viper.GetStringSlice("UnauthenticatedPaths")
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
		Port:                 port,
		BackendURL:           backendURL,
		JWTProviders:         providers,
		UnauthenticatedPaths: unAuthPaths,
	}
	return cfg, nil
}
