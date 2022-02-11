package config

import (
	"github.com/kelseyhightower/envconfig"
	gonanoid "github.com/matoous/go-nanoid"
	"strings"
	"time"
)

type (
	Config struct {
		Database Database
		HTTP     HTTP
		S3       S3
		Storage  Storage
		Server   Server
		Session  Session
		Guest    Guest
	}

	Database struct {
		Driver     string `envconfig:"FILEX_DATABASE_DRIVER" default:"sqlite3"`
		DataSource string `envconfig:"FILEX_DATABASE_DATASOURCE" default:""`
	}

	Storage struct {
		LocalPath string `envconfig:"FILEX_STORAGE_LOCAL_PATH"`
	}

	S3 struct {
		Bucket    string `envconfig:"FILEX_S3_BUCKET"`
		Prefix    string `envconfig:"FILEX_S3_PREFIX"`
		Endpoint  string `envconfig:"FILEX_S3_ENDPOINT"`
		PathStyle bool   `envconfig:"FILEX_S3_PATH_STYLE"`
	}

	Server struct {
		Addr  string `envconfig:"-"`
		Host  string `envconfig:"FILEX_SERVER_HOST" default:"localhost:80"`
		Port  string `envconfig:"FILEX_SERVER_PORT" default:":80"`
		Proto string `envconfig:"FILEX_SERVER_PROTO" default:"http"`
		Cert  string `envconfig:"FILEX_TLS_CERT"`
		Key   string `envconfig:"FILEX_TLS_KEY"`
	}

	HTTP struct {
		AllowedHosts          []string          `envconfig:"FILEX_HTTP_ALLOWED_HOSTS"`
		HostsProxyHeaders     []string          `envconfig:"FILEX_HTTP_PROXY_HEADERS"`
		SSLRedirect           bool              `envconfig:"FILEX_HTTP_SSL_REDIRECT"`
		SSLTemporaryRedirect  bool              `envconfig:"FILEX_HTTP_SSL_TEMPORARY_REDIRECT"`
		SSLHost               string            `envconfig:"FILEX_HTTP_SSL_HOST"`
		SSLProxyHeaders       map[string]string `envconfig:"FILEX_HTTP_SSL_PROXY_HEADERS"`
		STSSeconds            int64             `envconfig:"FILEX_HTTP_STS_SECONDS"`
		STSIncludeSubdomains  bool              `envconfig:"FILEX_HTTP_STS_INCLUDE_SUBDOMAINS"`
		STSPreload            bool              `envconfig:"FILEX_HTTP_STS_PRELOAD"`
		ForceSTSHeader        bool              `envconfig:"FILEX_HTTP_STS_FORCE_HEADER"`
		BrowserXSSFilter      bool              `envconfig:"FILEX_HTTP_BROWSER_XSS_FILTER"    default:"true"`
		FrameDeny             bool              `envconfig:"FILEX_HTTP_FRAME_DENY"            default:"true"`
		ContentTypeNosniff    bool              `envconfig:"FILEX_HTTP_CONTENT_TYPE_NO_SNIFF"`
		ContentSecurityPolicy string            `envconfig:"FILEX_HTTP_CONTENT_SECURITY_POLICY"`
		ReferrerPolicy        string            `envconfig:"FILEX_HTTP_REFERRER_POLICY"`
	}

	Session struct {
		Timeout time.Duration `envconfig:"FILEX_COOKIE_TIMEOUT" default:"720h"`
		Secret  string        `envconfig:"FILEX_COOKIE_SECRET"`
		Secure  bool          `envconfig:"FILEX_COOKIE_SECURE"`
	}

	Guest struct {
		AllowUpload     bool          `envconfig:"FILEX_GUEST_ALLOW_UPLOAD" default:"true"`
		MaxUploadSize   int64         `envconfig:"FILEX_GUEST_UPLOAD_MAX_SIZE" default:"2097152"`
		MaxFileDuration time.Duration `envconfig:"FILEX_GUEST_UPLOAD_MAX_FILE_DURATION" default:"168h"`
	}
)

func Environ() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	defaultAddress(&cfg)
	defaultSession(&cfg)
	return cfg, err
}

func cleanHostname(hostname string) string {
	hostname = strings.ToLower(hostname)
	hostname = strings.TrimPrefix(hostname, "http://")
	hostname = strings.TrimPrefix(hostname, "https://")

	return hostname
}

func defaultAddress(c *Config) {
	if c.Server.Key != "" || c.Server.Cert != "" {
		c.Server.Port = ":443"
		c.Server.Proto = "https"
	}
	c.Server.Host = cleanHostname(c.Server.Host)
	c.Server.Addr = c.Server.Proto + "://" + c.Server.Host
}

func defaultSession(c *Config) {
	if c.Session.Secret == "" {
		c.Session.Secret = gonanoid.MustID(32)
	}
}
