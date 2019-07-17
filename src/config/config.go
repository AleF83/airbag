package config

import (
	"fmt"
	"net/url"
	"regexp"

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

	viper.SetConfigName("airbag-config")
	viper.AddConfigPath("/app/data/")
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

	cfg := &Config{
		Port:                 80, // TODO: get from env
		BackendURL:           backendURL,
		JWTProviders:         providers,
		UnauthenticatedPaths: unAuthPaths,
	}
	return cfg, nil
}
