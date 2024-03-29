package util

import (
	"encoding/json"
	"os"
)

type DBConfig struct {
	DBDialect     string `json:"db_dialect"`
	DBPath        string `json:"db_path"`
	Password      string `json:"password"`
	Username      string `json:"username"`
	MaxIdleConns  int    `json:"max_idle_conns"`
	MaxOpenConns  int    `json:"max_open_conns"`
	AWSRegion     string `json:"aws_region"`
	AWSSecretName string `json:"aws_secret_name"`
}

type LogConfig struct {
	Level                        string `json:"level"`
	Filename                     string `json:"filename"`
	MaxFileSizeInMB              int    `json:"max_file_size_in_mb"`
	MaxBackupsOfLogFiles         int    `json:"max_backups_of_log_files"`
	MaxAgeToRetainLogFilesInDays int    `json:"max_age_to_retain_log_files_in_days"`
	UseConsoleLogger             bool   `json:"use_console_logger"`
	UseFileLogger                bool   `json:"use_file_logger"`
	Compress                     bool   `json:"compress"`
}

type APIConfig struct {
	EnableCache bool `json:"enable_cache"`
}

type ServerConfig struct {
	Env       string     `json:"env"`
	DBConfig  *DBConfig  `json:"db_config"`
	APIConfig *APIConfig `json:"api_config"`
	LogConfig *LogConfig `json:"log_config"`
}

func ParseServerConfigFromFile(filePath string) *ServerConfig {
	bz, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var config ServerConfig
	if err := json.Unmarshal(bz, &config); err != nil {
		panic(err)
	}

	if config.DBConfig.Username == "" || config.DBConfig.Password == "" { // read password from ENV
		config.DBConfig.Username, config.DBConfig.Password = GetDBUsernamePasswordFromEnv()
	}
	if config.DBConfig.Username == "" || config.DBConfig.Password == "" { // read password from AWS secret
		config.DBConfig.Username, config.DBConfig.Password = GetDBUsernamePasswordFromSM(config.DBConfig) // get from env
	}

	return &config
}

type MonitorConfig struct {
	Env string `json:"env"`

	BscRpcAddrs            []string `json:"bsc_rpc_addrs"`
	BscBlocksForFinality   int      `json:"bsc_blocks_for_finality"`
	BscMarketplaceContract string   `json:"bsc_marketplace_contract"`
	BscStartHeight         uint64   `json:"bsc_start_height"`

	GnfdRpcAddrs    []string `json:"gnfd_rpc_addrs"`
	GnfdChainId     string   `json:"gnfd_chain_id"`
	GnfdStartHeight uint64   `json:"gnfd_start_height"`

	GroupBucketRegex  string `json:"group_bucket_regex"`  // example "dm_b_.*"
	GroupBucketPrefix string `json:"group_bucket_prefix"` // example "dm_b_"
	GroupObjectRegex  string `json:"group_object_regex"`  // example "dm_o_.*"
	GroupObjectPrefix string `json:"group_object_prefix"` // example "dm_o_"

	DBConfig  *DBConfig  `json:"db_config"`
	LogConfig *LogConfig `json:"log_config"`
}

func ParseMonitorConfigFromFile(filePath string) *MonitorConfig {
	bz, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	var config MonitorConfig
	if err := json.Unmarshal(bz, &config); err != nil {
		panic(err)
	}

	if config.DBConfig.Username == "" || config.DBConfig.Password == "" { // read password from ENV
		config.DBConfig.Username, config.DBConfig.Password = GetDBUsernamePasswordFromEnv()
	}
	if config.DBConfig.Username == "" || config.DBConfig.Password == "" { // read password from AWS secret
		config.DBConfig.Username, config.DBConfig.Password = GetDBUsernamePasswordFromSM(config.DBConfig) // get from env
	}

	return &config
}

func GetDBUsernamePasswordFromEnv() (string, string) {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	return username, password
}

func GetDBUsernamePasswordFromSM(cfg *DBConfig) (string, string) {
	result, err := GetSecret(cfg.AWSSecretName, cfg.AWSRegion)
	if err != nil {
		panic(err)
	}
	type DBPass struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var dbPassword DBPass
	err = json.Unmarshal([]byte(result), &dbPassword)
	if err != nil {
		panic(err)
	}
	return dbPassword.Username, dbPassword.Password
}
