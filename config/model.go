package config

import "time"

type Config struct {
	Retry            int `toml:"retry"`             // retry times
	Workers          int `toml:"workers"`           // worker count
	Threads          int `toml:"threads"`           // download threads for each task
	Stream           bool `toml:"stream"`            // enable stream mode
	NoCleanCache     bool `toml:"no_clean_cache"`    // don't clean cache on exit
	Lang             string `toml:"lang"`              // language
	Telegram         Telegram `toml:"telegram"`         // telegram config
	Storages         []Storage `toml:"storages"`         // storage list
	Users            []User `toml:"users"`             // user list
	Hook             Hook `toml:"hook"`               // hook config
	Temp             Temp `toml:"temp"`               // temp config
	Notification     Notification `toml:"notification"`   // notification config
	AIRename         AIRename `toml:"ai_rename"`       // AI rename config
}

type Telegram struct {
	Token   string `toml:"token"`      // bot token
	AppID   int `toml:"app_id"`     // telegram app id
	AppHash string `toml:"app_hash"`  // telegram app hash
	Proxy   Proxy `toml:"proxy"`     // proxy config
	Userbot Userbot `toml:"userbot"`  // userbot config
}

type Proxy struct {
	Enable bool `toml:"enable"`     // enable proxy
	URL    string `toml:"url"`       // proxy url
}

type Userbot struct {
	Enable    bool `toml:"enable"`      // enable userbot
	SessionDB string `toml:"session_db"` // session db path
}

type Storage struct {
	Name     string `toml:"name"`       // display name
	Type     string `toml:"type"`       // storage type
	Enable   bool `toml:"enable"`      // enable storage
	BasePath string `toml:"base_path"`  // base path

	// local specific
	Owner string `toml:"owner"`       // file owner
	Group string `toml:"group"`       // file group
	Mode  uint32 `toml:"mode"`        // file mode

	// webdav specific
	URL      string `toml:"url"`       // webdav url
	Username string `toml:"username"`  // webdav username
	Password string `toml:"password"`  // webdav password

	// s3 specific
	Endpoint        string `toml:"endpoint"`         // s3 endpoint
	AccessKey       string `toml:"access_key"`       // s3 access key
	SecretKey       string `toml:"secret_key"`       // s3 secret key
	Bucket          string `toml:"bucket"`           // s3 bucket
	Region          string `toml:"region"`           // s3 region
	NoSSL           bool `toml:"no_ssl"`            // disable ssl
	PathStyle       bool `toml:"path_style"`        // use path style
	CustomDomain    string `toml:"custom_domain"`     // custom domain for s3
	Cloudfront      bool `toml:"cloudfront"`        // use cloudfront
	UploadConcurrency int `toml:"upload_concurrency"` // upload concurrency

	// alist specific
	Token string `toml:"token"`       // alist token

	// telegram storage specific
	BotToken string `toml:"bot_token"`  // bot token for telegram storage
	ChatID   int64 `toml:"chat_id"`    // chat id for telegram storage
}

type User struct {
	ID        int64 `toml:"id"`          // telegram user id
	Storages  []string `toml:"storages"`  // storage filter list
	Blacklist bool `toml:"blacklist"`   // use blacklist mode
}

type Hook struct {
	Exec Exec `toml:"exec"` // exec hook
}

type Exec struct {
	TaskBeforeStart string `toml:"task_before_start"` // exec before task start
	TaskFail        string `toml:"task_fail"`         // exec when task failed
	TaskSuccess     string `toml:"task_success"`      // exec when task success
	TaskCancel      string `toml:"task_cancel"`       // exec when task canceled
}

type Temp struct {
	BasePath string `toml:"base_path"` // temp file base path
}

type Notification struct {
	Progress ProgressNotification `toml:"progress"` // progress notification
}

type ProgressNotification struct {
	Enable       bool `toml:"enable"`        // enable progress notification
	Silent       bool `toml:"silent"`        // silent notification
	Interval     time.Duration `toml:"interval"`  // notification interval
	EditTimeout  time.Duration `toml:"edit_timeout"` // edit timeout
	Batch        BatchProgress `toml:"batch"`     // batch progress
}

type BatchProgress struct {
	Enable        bool `toml:"enable"`         // enable batch progress
	UpdateOnce    bool `toml:"update_once"`    // update progress once
	DetailFailed  bool `toml:"detail_failed"`  // show failed detail
	DetailSuccess bool `toml:"detail_success"` // show success detail
}

type AIRename struct {
	Enable     bool `toml:"enable"`      // enable AI rename
	Endpoint   string `toml:"endpoint"`   // OpenAI compatible API endpoint
	Model      string `toml:"model"`      // model name
	APIKey     string `toml:"api_key"`    // API key
	Prompt     string `toml:"prompt"`     // custom prompt template
	Timeout    time.Duration `toml:"timeout"`   // request timeout
	MaxTokens  int `toml:"max_tokens"`   // max tokens for response
	Temperature float32 `toml:"temperature"` // temperature for AI
}

type Rule struct {
	UserID   int64 `toml:"user_id"`    // user id (0 for global)
	Name     string `toml:"name"`      // rule name
	Type     string `toml:"type"`      // rule type
	Rule     string `toml:"rule"`      // rule content
	Value    string `toml:"value"`     // rule value
	Storage  string `toml:"storage"`   // storage name
	Path     string `toml:"path"`      // save path
	Enable   bool `toml:"enable"`     // enable rule
	Priority int `toml:"priority"`    // priority (higher = more priority)
}
