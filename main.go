package main

import (
    "fmt"
    "log"
    "mime"
    "net/http"
    "github.com/spf13/viper"
)

type Config struct {
    AuthSecret      string  `mapstructure:"secret"`
    CloudflareToken string  `mapstructure:"cloudflareToken"`
    CloudflareZone  string  `mapstructure:"cloudflareZone"`
    HcloudToken     string  `mapstructure:"hcloudToken"`
    HcloudSSHKeys []string  `mapstructure:"hcloudSSHKeys"`
    ServerISO       string  `mapstructure:"serverISO"`
    ServerID        string  `mapstructure:"serverID"`
    ServerType      string  `mapstructure:"serverType"`
    ServerLocation  string  `mapstructure:"serverLocation"`
    VolumeName      string  `mapstructure:"volumeName"`
    VolumeSize      string  `mapstructure:"volumeSize"`
}

func getConfig() (*Config, error) {
    var config *Config

    // basic viper conf
    v := viper.New()
    v.AddConfigPath("./")
    v.SetConfigName("config")
    v.SetConfigType("yaml")

    // overwrite if env variables exists
    v.AutomaticEnv()

    err := v.ReadInConfig()
    
    if err != nil {
        return nil, err
    }

    err = v.Unmarshal(&config)
    log.Print(config)    
    if err != nil {
        return nil, err
    }

    return config, nil
}

func enforceJSONHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        contentType := r.Header.Get("Content-Type")

        if contentType == "" {
            http.Error(w, "Missing Content-Type header", http.StatusBadRequest)
            return
        }

        mt, _, err := mime.ParseMediaType(contentType)
		
        if err != nil {
			http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
            return
		}

        if mt != "application/json" {
            http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func (config *Config) authHandler(next http.Handler) http.Handler {
    log.Print("Using secret: ", config.AuthSecret)
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")

        if authHeader == "" {
            http.Error(w, "Missing Authorization header", http.StatusBadRequest)
            return
        }

        if authHeader != fmt.Sprintf("Bearer %s", config.AuthSecret) {
            http.Error(w, "Invalid auth token", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func (config *Config) provisionServer(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
    w.Write([]byte(config.ServerID))
}

func main() {
    // Get the config from env variables or use defaults
    config, err := getConfig()
    
    if err != nil {
	log.Fatal("Could not load config file: ", err)
    }

    log.Print("Using configuration:", config)

    mux := http.NewServeMux()
    provisionHandler := http.HandlerFunc(config.provisionServer)
    mux.Handle("/", enforceJSONHandler(config.authHandler(provisionHandler)))

    log.Print("Listening on :3000...")
    httpErr := http.ListenAndServe(":3000", mux)
    log.Fatal(httpErr)
}		
