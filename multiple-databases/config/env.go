package config

import "time"

type EnvironmentVariable struct {
	App struct {
		Protocol string `mapstructure:"PROTOCOL"`
		Host     string `mapstructure:"HOST"`
		Mode     string `mapstructure:"MODE"`
		Port     int    `mapstructure:"PORT"`
		Debug    bool   `mapstructure:"DEBUG"`
	} `mapstructure:"APP"`

	Database struct {
		MongoLog struct {
			Host     string `mapstructure:"HOST"`
			Port     int    `mapstructure:"PORT"`
			Name     string `mapstructure:"NAME"`
			Username string `mapstructure:"USERNAME"`
			Password string `mapstructure:"PASSWORD"`
		} `mapstructure:"MONGO_LOG"`

		SqliteApp struct {
			Dir  string `mapstructure:"DIR"`
			Name string `mapstructure:"NAME"`
		} `mapstructure:"SQLITE_APP"`

		SqliteLog struct {
			Dir  string `mapstructure:"DIR"`
			Name string `mapstructure:"NAME"`
		} `mapstructure:"SQLITE_LOG"`

		Timeout struct {
			Ping  time.Duration `mapstructure:"PING"`
			Read  time.Duration `mapstructure:"READ"`
			Write time.Duration `mapstructure:"Write"`
		} `mapstructure:"TIMEOUT"`
	} `mapstructure:"DATABASE"`

	Logging struct {
		LogJobExecution struct {
			CheckLockCompressionEnable bool `mapstructure:"CHECK_LOCK_COMPRESSION_ENABLE"`
			FindCompressionEnable      bool `mapstructure:"FIND_COMPRESSION_ENABLE"`
			EvaluateRulesEnable        bool `mapstructure:"EVALUATE_RULES_ENABLE"`
			CreateJobEnable            bool `mapstructure:"CREATE_JOB_ENABLE"`
			PublishJobEnable           bool `mapstructure:"PUBLISH_JOB_ENABLE"`
		} `mapstructure:"LOG_JOB_EXECUTION"`
		LogTaskExecution struct {
			TestConnectionDatabaseEnable bool `mapstructure:"TEST_CONNECTION_DATABASE_ENABLE"`
			ParsePayloadEnable           bool `mapstructure:"PARSE_PAYLOAD_ENABLE"`
			DownloadEnable               bool `mapstructure:"DOWNLOAD_ENABLE"`
			CompressEnable               bool `mapstructure:"COMPRESS_ENABLE"`
			UploadEnable                 bool `mapstructure:"UPLOAD_ENABLE"`
			DeleteSourceEnable           bool `mapstructure:"DELETE_SOURCE_ENABLE"`
		} `mapstructure:"LOG_TASK_EXECUTION"`
	} `mapstructure:"LOGGING"`
}
