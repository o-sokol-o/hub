package domain

// !!! Изменён пакет "github.com/gin-contrib/sessions/memstore"
import (
	"time"
)

//=====================================================================================================================
//=====================================================================================================================

type CfgAws struct {
	AccessKeyID                string `envconfig:"AWS_ACCESS_KEY_ID"`              // WEB_API_AWS_AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY_ID
	SecretAccessKey            string `envconfig:"AWS_SECRET_ACCESS_KEY" json:"-"` // don't print
	Region                     string `default:"us-west-2" envconfig:"AWS_DEFAULT_REGION"`
	S3BucketPrivate            string `envconfig:"S3_BUCKET_PRIVATE"`
	S3BucketPublic             string `envconfig:"S3_BUCKET_PUBLIC"`
	SecretsManagerConfigPrefix string `default:"" envconfig:"SECRETS_MANAGER_CONFIG_PREFIX"`

	// Get an AWS session from an implicit source if no explicit
	// configuration is provided. This is useful for taking advantage of
	// EC2/ECS instance roles.
	UseRole bool `envconfig:"AWS_USE_ROLE"`
}
type CfgAuth struct {
	UseAwsSecretManager bool          `default:"false" envconfig:"USE_AWS_SECRET_MANAGER"`
	KeyExpiration       time.Duration `default:"3600s" envconfig:"KEY_EXPIRATION"`
}

type CfgProject struct {
	Name              string `default:"" envconfig:"PROJECT_NAME"`
	SharedTemplateDir string `default:"../../resources/templates/shared" envconfig:"SHARED_TEMPLATE_DIR"`
	SharedSecretKey   string `default:"" envconfig:"SHARED_SECRET_KEY"`
	EmailSender       string `default:"test@example.com" envconfig:"EMAIL_SENDER"`
	WebApiBaseUrl     string `default:"http://127.0.0.1:3001" envconfig:"WEB_API_BASE_URL"  example:"http://api.example.com"`
}

type CfgDB struct {
	Host       string `default:"http://127.0.0.1:5432" envconfig:"HOST"`
	User       string `default:"postgres" envconfig:"USER"`
	Pass       string `default:"qwerty" envconfig:"PASS" json:"-"` // don't print
	Database   string `default:"postgres" envconfig:"DATABASE"`
	Driver     string `default:"postgres" envconfig:"DRIVER"`
	Timezone   string `default:"utc" envconfig:"TIMEZONE"`
	DisableTLS bool   `default:"true" envconfig:"DISABLE_TLS"`
}

type CfgService struct {
	Name        string   `default:"web-app" envconfig:"SERVICE_NAME"`
	BaseUrl     string   `default:"" envconfig:"BASE_URL"  example:"http://example.com"`
	HostNames   []string `envconfig:"HOST_NAMES" example:"www.example.com"`
	EnableHTTPS bool     `default:"false" envconfig:"ENABLE_HTTPS"`
	TemplateDir string   `default:"./templates" envconfig:"TEMPLATE_DIR"`
	StaticFiles struct {
		Dir               string `default:"./static" envconfig:"STATIC_DIR"`
		S3Enabled         bool   `envconfig:"S3_ENABLED"`
		S3Prefix          string `default:"public/web_app/static" envconfig:"S3_PREFIX"`
		CloudFrontEnabled bool   `envconfig:"CLOUDFRONT_ENABLED"`
		ImgResizeEnabled  bool   `envconfig:"IMG_RESIZE_ENABLED"`
	}
	Minify          bool          `envconfig:"MINIFY"`
	SessionName     string        `default:"" envconfig:"SESSION_NAME"`
	DebugHost       string        `default:"0.0.0.0:4000" envconfig:"DEBUG_HOST"`
	ShutdownTimeout time.Duration `default:"5s" envconfig:"SHUTDOWN_TIMEOUT"`
	ScaleToZero     time.Duration `envconfig:"SCALE_TO_ZERO"`
}

type CfgHTTP struct {
	Host         string        `default:"127.0.0.1:3000" envconfig:"HOST"`
	Port         string        `default:"3000" envconfig:"PORT"`
	ReadTimeout  time.Duration `default:"10s" envconfig:"READ_TIMEOUT"`
	WriteTimeout time.Duration `default:"10s" envconfig:"WRITE_TIMEOUT"`
}

type CfgHTTPS struct {
	Host         string        `default:"" envconfig:"HOST"`
	Port         string        `default:"" envconfig:"PORT"`
	ReadTimeout  time.Duration `default:"5s" envconfig:"READ_TIMEOUT"`
	WriteTimeout time.Duration `default:"5s" envconfig:"WRITE_TIMEOUT"`
}

type CfgRedis struct {
	Host            string        `default:"127.0.0.1:6379" envconfig:"HOST"`
	DB              int           `default:"1" envconfig:"DB"`
	DialTimeout     time.Duration `default:"5s" envconfig:"DIAL_TIMEOUT"`
	MaxmemoryPolicy string        `default:"allkeys-lru" envconfig:"MAXMEMORY_POLICY"`
}

type CfgTrace struct {
	Host          string  `default:"127.0.0.1" envconfig:"DD_TRACE_AGENT_HOSTNAME"`
	Port          int     `default:"8126" envconfig:"DD_TRACE_AGENT_PORT"`
	AnalyticsRate float64 `default:"0.10" envconfig:"ANALYTICS_RATE"`
}

type CfgBuildInfo struct {
	CiCommitRefName  string `envconfig:"CI_COMMIT_REF_NAME"`
	CiCommitShortSha string `envconfig:"CI_COMMIT_SHORT_SHA"`
	CiCommitSha      string `envconfig:"CI_COMMIT_SHA"`
	CiCommitTag      string `envconfig:"CI_COMMIT_TAG"`
	CiJobId          string `envconfig:"CI_JOB_ID"`
	CiJobUrl         string `envconfig:"CI_JOB_URL"`
	CiPipelineId     string `envconfig:"CI_PIPELINE_ID"`
	CiPipelineUrl    string `envconfig:"CI_PIPELINE_URL"`
}
type CfgApplication struct {
	Env       string `default:"dev" envconfig:"ENV"`
	HTTP      CfgHTTP
	HTTPS     CfgHTTPS
	Service   CfgService
	Project   CfgProject
	Redis     CfgRedis
	DB        CfgDB
	Trace     CfgTrace
	Aws       CfgAws
	Auth      CfgAuth
	Build     string
	BuildInfo CfgBuildInfo
}

//=====================================================================================================================
//=====================================================================================================================
//=====================================================================================================================
//=====================================================================================================================
