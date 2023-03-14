package setting

var (
	DatabaseSettings *DatabaseSetting
	JWTSettings      *JWTSetting
)

type JWTSetting struct {
	Secret string `yaml:"Secret"`
}
type DatabaseSetting struct {
	DBType   string `yaml:"DBType"`
	UserName string `yaml:"UserName"`
	Password string `yaml:"Password"`
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	DBName   string `yaml:"DBName"`
	Charset  string `yaml:"Charset"`
}
