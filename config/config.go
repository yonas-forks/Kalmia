package config

import (
	"encoding/json"
	"os"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

type S3 struct {
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Bucket          string `json:"bucket"`
	UsePathStyle    bool   `json:"usePathStyle"`
	PublicUrlFormat string `json:"publicUrlFormat"`
}

type GithubOAuth struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURL  string `json:"callbackUrl"`
}

type MicrosoftOAuth struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	DirectoryID  string `json:"directoryId"`
	RedirectURL  string `json:"callbackUrl"`
}

type GoogleOAuth struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RedirectURL  string `json:"callbackUrl"`
}

type Config struct {
	Environment    string         `json:"environment"`
	Port           int            `json:"port"`
	Database       string         `json:"database"`
	LogLevel       string         `json:"logLevel"`
	AssetStorage   string         `json:"assetStorage"`
	MaxFileSize    int64          `json:"maxFileSize"` // in MB
	SessionSecret  string         `json:"sessionSecret"`
	Admins         []User         `json:"users"`
	DataPath       string         `json:"dataPath"`
	S3             S3             `json:"s3"`
	GithubOAuth    GithubOAuth    `json:"githubOAuth"`
	MicrosoftOAuth MicrosoftOAuth `json:"microsoftOAuth"`
	GoogleOAuth    GoogleOAuth    `json:"googleOAuth"`
}

var ParsedConfig *Config

func ParseConfig(path string) *Config {
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	ParsedConfig = &Config{}
	err = decoder.Decode(ParsedConfig)

	if err != nil {
		panic(err)
	}

	err = SetupDataPath()

	if err != nil {
		panic(err)
	}
	// INFO: Adds the default max file size of 10
	// Added for backwards compatibility
	if ParsedConfig.MaxFileSize == 0 {
		ParsedConfig.MaxFileSize = 10
	}

	return ParsedConfig
}

func SetupDataPath() error {
	if ParsedConfig.DataPath == "" {
		ParsedConfig.DataPath = "./data"
	}

	if _, err := os.Stat(ParsedConfig.DataPath); os.IsNotExist(err) {
		err := os.Mkdir(ParsedConfig.DataPath, 0755)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(ParsedConfig.DataPath + "/rspress_data"); os.IsNotExist(err) {
		err := os.Mkdir(ParsedConfig.DataPath+"/rspress_data", 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
