package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

// OpenTelemetry configuration
type OpenTelemetry struct {
	ServiceName    string `mapstructure:"ServiceName"`
	JaegerEndpoint string `mapstructure:"JaegerEndpoint"`
	OTLPEndpoint   string `mapstructure:"OtlpEndpoint"`
	EnableStdout   bool   `mapstructure:"EnableStdout"`
}

// WebUISection defines beamline section configuration
type WebUISection struct {
	Section    string   `mapstructure:"section"`
	Attributes []string `mapstructure:"attributes"`
}

// BeamlineSection defines beamline section configuration
type BeamlineSection struct {
	Schema   string         `mapstructure:"schema"`
	Sections []WebUISection `mapstructure:"sections"`
}

// BeamlineSections defines beamline section configuration on web UI
type BeamlineSections struct {
	Sections []BeamlineSection `mapstructure:"sections"`
}

// AIChat configuration
type AIChat struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Model   string `mapstructure:"model"`
	Group   string `mapstructure:"group"`   // aichat group to use
	Client  string `mapstructure:"client"`  // aichat client, e.g. ollaman or tichy
	Timeout int    `mapstructure:"timeout"` // ai response timeout
}

// FoxdenUser represents foxden user interface to use
type FoxdenUser struct {
	User string `mapstructure:"User"`
}

// TrustedUser represents single trusted user information
type TrustedUser struct {
	User string `mapstructure:"user"`
	IP   string `mapstructure:"ip"`
	MAC  string `mapstructure:"mac"`
}

// Globug represents globus information
type Globus struct {
	ClientID       string `mapstructure:"client_id"`       // client id
	ClientSecret   string `mapstructure:"client_secret"`   // client secret
	TransferURL    string `mapstructure:"transfer_url"`    // globus transfer url
	AuthURL        string `mapstructure:"auth_url"`        // globus auth url
	OriginID       string `mapstructure:"origin_id"`       // globus origin ID for CHESS Raw data collection
	CollectionPath string `mapstructure:"collection_path"` // globus collection path
}

// LDAP attributes
type LDAP struct {
	URL            string `mapstructure:"url"`           // ldap url
	BaseDN         string `mapstructure:"baseDN`         // ldap baseDN
	Login          string `mapstructure:"login"`         // LDAP login to use
	Password       string `mapstructure:"password"`      // LDAP password to use
	Expire         int    `mapstructure:"expire"`        // LDAP cache record expire (in seconds)
	RecursionLevel int    `mapstructure:recursion_level` // LDAP look-up recursion level
}

// DOI attributes
type DOI struct {
	Provider         string           `mapstructure:"Provider"`    // doi provider, e.g. Zenodo or MaterialsCommons
	ProjectName      string           `mapstructure:"ProjectName"` // name of the project (only valid for MaterialsCommons)
	Zenodo           Zenodo           `mapstructure:"Zenodo"`
	Datacite         Datacite         `mapstructure:"Datacite"`
	MaterialsCommons MaterialsCommons `mapstructure:"MaterialsCommons"`
	WebServer
}

// DID structure
type DID struct {
	Attributes string `mapstructure:"attributes"` // did attributes, comma separated, default beamline,btr,cycle,sample
	Separator  string `mapstructure:"separator"`  // did separator, default "/"
	Divider    string `mapstructure:"divider"`    // did key-value divider, default "="
}

// Embed structure
type Embed struct {
	DocDb string `mapstructure:"DocDb"`
	SqlDb string `mapstructure:"SqlDb"`
}

// QL structure
type QL struct {
	ServiceMapFile string `mapstructure:"ServiceMapFile"` // service map file name
	Separator      string `mapstructure:"separator"`      // ql separator, default ":"
	Verbose        int    `mapstructure:"verbose"`        // verbosity level
}

// OAuthRecord defines OAuth provider's credentials
type OAuthRecord struct {
	Provider     string `mapstructure:"Provider"`     // name of the provider
	ClientID     string `mapstructure:"ClientId"`     // client id
	ClientSecret string `mapstructure:"ClientSecret"` // client secret
	RedirectURL  string `mapstructure:"RedirectUrl"`  // redirect url
}

// MLAPI defines ML API structure
type MLApi struct {
	Name     string `mapstructure:"name"`
	Method   string `mapstructure:"method"`
	Endpoint string `mapstructure:"endpoint"`
	Accept   string `mapstructure:"accept"`
}

// MLBackend represents ML backend engine
type MLBackend struct {
	Name string  `mapstructure:"name"` // ML backend name, e.g. TFaaS
	Type string  `mapstructure:"type"` // ML backebd type, e.g. TensorFlow
	URI  string  `mapstructure:"uri"`  // ML backend URI, e.g. http://localhost:port
	Apis []MLApi // ML APIs
}

// ML defines ML configuration options
type ML struct {
	MLBackends []MLBackend `mapstructure:"MLBackends"` // ML backends
	StorageDir string      `mapstructure:"StorageDir"`
}

// Kerberos defines kerberos optinos
type Kerberos struct {
	Krb5Conf string `mapstructure:"Krb5Conf"`
	Keytab   string `mapstructure:"Keytab"`
	Realm    string `mapstructure:"Realm"`
}

// GinOptions controls go-gin specific options
type GinOptions struct {
	DisableConsoleColor bool   `mapstructure:"DisableConsoleColor"` // gin console color mode
	Production          bool   `mapstructure:"Production"`          // production mode
	Mode                string `mapstructure:"Mode"`                // gin mode: test, debug, release
}

// WebServer represents common web server configuration
type WebServer struct {
	// git server options
	GinOptions `mapstructure:"GinOptions"`

	// basic options
	Port        int    `mapstructure:"Port"`        // server port number
	Verbose     int    `mapstructure:"Verbose"`     // verbose output
	Base        string `mapstructure:"Base"`        // base URL
	StaticDir   string `mapstructure:"StaticDir"`   // specify static dir location
	LogFile     string `mapstructure:"LogFile"`     // server log file
	LogLongFile bool   `mapstructure:"LogLongFile"` // server log structure

	// middleware server parts
	LimiterPeriod   string   `mapstructure:"Rate"`              // limiter rate value
	LimiterHeader   string   `mapstructure:"limiter_header"`    // limiter header to use
	LimiterSkipList []string `mapstructure:"limiter_skip_list"` // limiter skip list
	MetricsPrefix   string   `mapstructure:"metrics_prefix"`    // metrics prefix used for Prometheus

	// etag options
	Etag         string `mapstructure:"etag"`          // etag value to use for ETag generation
	CacheControl string `mapstructure:"cache_control"` // Cache-Control value, e.g. max-age=300

	// proxy parts
	XForwardedHost      string `mapstructure:"X-Forwarded-Host"`       // X-Forwarded-Host field of HTTP request
	XContentTypeOptions string `mapstructure:"X-Content-Type-Options"` // X-Content-Type-Options option

	// TLS server parts
	RootCAs     string   `mapstructure:"RootCAs"`     // server Root CAs path
	ServerCrt   string   `mapstructure:"ServerCert"`  // server certificate
	ServerKey   string   `mapstructure:"ServerKey"`   // server certificate
	DomainNames []string `mapstructure:"DomainNames"` // LetsEncrypt domain names
}

// String provides string representation of WebServer structure
func (w *WebServer) String() string {
	data, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return fmt.Sprintf("ERROR: unable to parse web server configuration, error %v", err)
	}
	return string(data)
}

// Frontend stores frontend configuration parameters
type Frontend struct {
	FoxdenUser `mapstructure:"FoxdenUser"`
	WebServer  `mapstructure:"WebServer"`

	// Aggregate results from all FOXDEN services
	AggregateResults bool `mapstructure:"AggregateResults`

	// OAuth parts
	OAuth []OAuthRecord `mapstructure:"OAuth"` // oauth configurations

	// captcha parts
	CaptchaSecretKey string `mapstructure:"CaptchaSecretKey"` // re-captcha secret key
	CaptchaPublicKey string `mapstructure:"CaptchaPublicKey"` // re-captcha public key
	CaptchaVerifyUrl string `mapstructure:"CaptchaVerifyUrl"` // re-captcha verify url

	// cookies parts
	UserCookieExpires int64 `mapstructure:"UserCookieExpires"` // expiration of user cookie

	// other options
	TestMode        bool   `mapstructure:"TestMode"`        // test mode
	CheckBtrs       bool   `mapstructure:"CheckBtrs"`       // enable check for CHESS btrs
	CheckAdmins     bool   `mapstructure:"CheckAdmins"`     // enable check for FOXDEN admins
	AllowAllRecords bool   `mapstructure:"AllowAllRecords"` // allow all records to be seen
	DefaultEndPoint string `mapstructure:"DefaultEndPoint"`
	DocUrl          string `mapstructure:"DocUrl"` // documentation url
}

// Encryption represents encryption configuration parameters
type Encryption struct {
	Secret string `mapstructure:"Secret"`
	Cipher string `mapstructure:"Cipher"`
}

// MongoDB represents MongoDB parameters
type MongoDB struct {
	DBName string `mapstructure:"DBName"` // database name
	DBColl string `mapstructure:"DBColl"` // database collection
	DBUri  string `mapstructure:"DBUri"`  // database URI
}

// MLHub represents ML service configuration
type MLHub struct {
	WebServer  `mapstructure:"WebServer"`
	MongoDB    `mapstructure:"MongoDB"`
	Encryption `mapstructure:"Encryption"`
	ML         `mapstructure:"ML"`
}

// DataHub represents DataHub service configuration
type DataHub struct {
	WebServer  `mapstructure:"WebServer"`
	Encryption `mapstructure:"Encryption"`
	StorageDir string `mapstructure:"StorageDir"`
}

// Discovery represents discovery service configuration
type Discovery struct {
	WebServer  `mapstructure:"WebServer"`
	MongoDB    `mapstructure:"MongoDB"`
	Encryption `mapstructure:"Encryption"`
}

// Sync represents Sync service configuration
type Sync struct {
	WebServer     `mapstructure:"WebServer"`
	MongoDB       `mapstructure:"MongoDB"`
	SleepInterval int `mapstructure:"SleepInterval"`
	NWorkers      int `mapstructure:"NumberOfWorkers"`
}

// MetaData represents metadata service configuration
type MetaData struct {
	FoxdenUser             `mapstructure:"FoxdenUser"`
	WebServer              `mapstructure:"WebServer"`
	MongoDB                `mapstructure:"MongoDB"`
	LexiconFile            string   `mapstructure:"LexiconFile"`            // lexicon file
	DataLocationAttributes []string `mapstructure:"DataLocationAttributes"` // data location attributes to use
}

// UserMetaData represents User MetaData configuration
type UserMetaData struct {
	FoxdenUser `mapstructure:"FoxdenUser"`
	WebServer  `mapstructure:"WebServer"`
	MongoDB    `mapstructure:"MongoDB"`
}

// CHESSMetaData represents CHESS MetaData configuration
type CHESSMetaData struct {
	FoxdenUser             `mapstructure:"FoxdenUser"`
	WebServer              `mapstructure:"WebServer"`
	MongoDB                `mapstructure:"MongoDB"`
	SchemaRenewInterval    int      `mapstructure:"SchemaRenewInterval"`    // schema renew interval
	WebSectionsFile        string   `mapstructure:"WebSectionsFile"`        // file for web form sections
	LexiconFile            string   `mapstructure:"LexiconFile"`            // lexicon file
	TestMode               bool     `mapstructure:"TestMode"`               // test mode
	DataLocationAttributes []string `mapstructure:"DataLocationAttributes"` // data location attributes to use
	SchemaFiles            []string `mapstructure:"SchemaFiles"`            // schema files
	OrderedSections        []string `mapstructure:"OrderedSections"`        // ordered sections for web UI
	SkipKeys               []string `mapstructure:"SkipKeys"`               // keys to skip for web forms
	SpecScanBeamlines      []string `mapstructure:"SpecScanBeamlines"`      // list of beamlines that uses spec scan service
}

// OreCastMetaData represents OreCast MetaData configuration
type OreCastMetaData struct {
	WebServer `mapstructure:"WebServer"`
	MongoDB   `mapstructure:"MongoDB"`
}

// Publication represents Publication service configuration
type Publication struct {
	WebServer `mapstructure:"WebServer"`
}

// Zenodo represents Zenodo service configuration
type Zenodo struct {
	Url         string `mapstructure:"Url"`
	AccessToken string `mapstructure:"AccessToken"`
}

// MaterialsCommons represents MaterialsCommons service configuration
type MaterialsCommons struct {
	Url                string `mapstructure:"Url"`
	AccessToken        string `mapstructure:"AccessToken"`
	ProjectName        string `mapstructure:"ProjectName"`        // name of the project (only valid for MaterialsCommons)
	ProductionInstance bool   `mapstructure:"ProductionInstance"` // specify to use production instance
}

// Datacite represents Datacite service configuration
type Datacite struct {
	Url          string `mapstructure:"Url"`
	AccessToken  string `mapstructure:"AccessToken"`
	AccessKey    string `mapstructure:"AccessKey"`
	AccessSecret string `mapstructure:"AccessSecret"`
	Username     string `mapstructure:"Username"`
	Password     string `mapstructure:"Password"`
	Prefix       string `mapstructure:"Prefix"`
}

// SpecScans represents SpecScansService configuration
type SpecScans struct {
	WebServer  `mapstructure:"WebServer"`
	MongoDB    `mapstructure:"MongoDB"`
	DBFile     string `mapstructure:"DBFile"`
	SchemaFile string `mapstructure:"SchemaFile"`
}

// S3 defines s3 structure
type S3 struct {
	Name         string `mapstructure:"Name"`
	AccessKey    string `mapstructure:"AccessKey"`
	AccessSecret string `mapstructure:"AccessSecret"`
	Endpoint     string `mapstructure:"Endpoint"`
	UseSSL       bool   `mapstructure:"UseSSL"`
	Region       string `mapstructure:"region"`
}

// FS defines FileSystem backend
type FS struct {
	Name    string `mapstructure:"Name"`
	Kind    string `mapstructure:"Kind"`
	Storage string `mapstructure:"Storage"`
}

// DataManagement represents data-management service configuration
type DataManagement struct {
	S3
	FS
	FileExtensions []string `mapstructure:"FileExtensions"`
	WebServer      `mapstructure:"WebServer"`
}

// DataBookkeeping represents data-bookkeeping service configuration
type DataBookkeeping struct {
	WebServer `mapstructure:"WebServer"`

	DBFile             string `mapstructure:"DBFile"`             // dbs db file with secrets
	LexiconFile        string `mapstructure:"LexiconFile"`        // dbs lexicon file
	MaxDBConnections   int    `mapstructure:"MaxDbConnections"`   // maximum number of DB connections
	MaxIdleConnections int    `mapstructure:"MaxIdleConnections"` // maximum number of idle connections
}

// Authz represents authz service configuration
type Authz struct {
	WebServer  `mapstructure:"WebServer"`
	Encryption `mapstructure:"Encryption"`

	TestMode     bool   `mapstructure:"TestMode"`  // test mode
	CheckLDAP    bool   `mapstructure:"CheckLDAP"` // check users scope in LDAP
	DBFile       string `mapstructure:"DBFile"`
	ClientID     string `mapstructure:"ClientId"`
	ClientSecret string `mapstructure:"ClientSecret"`
	Domain       string `mapstructure:"Domain"`
	TokenExpires int64  `mapstructure:"TokenExpires"` // token expiration in seconds
}

// Services represents services structure
type Services struct {
	FrontendURL        string `mapstructure:"FrontendUrl"`
	DiscoveryURL       string `mapstructure:"DiscoveryUrl"`
	MetaDataURL        string `mapstructure:"MetaDataUrl"`
	MLHubURL           string `mapstructure:"MLHubUrl"`
	DataHubURL         string `mapstructure:"DataHubUrl"`
	DataManagementURL  string `mapstructure:"DataManagementUrl"`
	DataBookkeepingURL string `mapstructure:"DataBookkeepingUrl"`
	AuthzURL           string `mapstructure:"AuthzUrl"`
	SpecScansURL       string `mapstructure:"SpecScansUrl"`
	PublicationURL     string `mapstructure:"PublicationUrl"`
	CHAPBookURL        string `mapstructure:"CHAPBookUrl"`
	DOIServiceURL      string `mapstructure:"DOIServiceUrl"`
	SyncServiceURL     string `mapstructure:"SyncServiceUrl"`
	UserMetaDataURL    string `mapstructure:"UserMetaDataUrl"`
}

// SrvConfig represents configuration structure
type SrvConfig struct {
	AIChat          `mapstructure:"AIChat"`
	S3              `mapstructure:"S3"`
	QL              `mapstructure:"QL"`
	DID             `mapstructure:"DID"`
	LDAP            `mapstructure:"LDAP"`
	Frontend        `mapstructure:"Frontend"`
	Discovery       `mapstructure:"Discovery"`
	MetaData        `mapstructure:"MetaData"`
	MLHub           `mapstructure:"MLHub"`
	DataHub         `mapstructure:"DataHub"`
	DataManagement  `mapstructure:"DataManagement"`
	DataBookkeeping `mapstructure:"DataBookkeeping"`
	Authz           `mapstructure:"Authz"`
	Kerberos        `mapstructure:"Kerberos"`
	Services        `mapstructure:"Services"`
	Encryption      `mapstructure:"Encryption"`
	CHESSMetaData   `mapstructure:"CHESSMetaData"`
	UserMetaData    `mapstructure:"UserMetaData"`
	OreCastMetaData `mapstructure:"OreCastMetaData"`
	SpecScans       `mapstructure:"SpecScansService"`
	Publication     `mapstructure:"PublicationService"`
	Globus          `mapstructure:"Globus"`
	DOI             `mapstructure:"DOI"`
	Sync            `mapstructure:"Sync"`
	Embed           `mapstructure:"Embed"`
	OpenTelemetry   `mapstructure:"OpenTelemetry"`

	TrustedUsers     []TrustedUser     `mapstructure:"TrustedUsers"`
	BeamlineSections []BeamlineSection `mapstructure:"BeamlineSections"`
}

// Config represents configuration instance object
var Config *SrvConfig

// String shows SrvConfig object string representation
func (c *SrvConfig) String() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Println("ERROR:", err)
		return fmt.Sprintf("%s", string(data))
	}
	return string(data)
}

// ParseConfig provides method to parse given configuration file
func ParseConfig(cfile string) (SrvConfig, error) {
	var config SrvConfig
	if cfile == "" {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("ERROR", err)
			return config, err
		}
		// setup cfile to $HOME/.foxden.yaml
		cfile = filepath.Join(home, ".foxden.yaml")
		// setup cfile to $FOXDEN_CONFIG if it is set in environment
		if _, err := os.Stat(os.Getenv("FOXDEN_CONFIG")); err == nil {
			cfile = os.Getenv("FOXDEN_CONFIG")
		}
	}

	// check if we do have configuration file
	if _, err := os.Stat(cfile); err == nil {
		viper.SetConfigFile(cfile)
		if os.Getenv("FOXDEN_DEBUG") != "" {
			fmt.Println("Parse FOXDEN config:", cfile)
		}
	} else {
		if _, err := os.Stat(os.Getenv("FOXDEN_CONFIG")); err == nil {
			if os.Getenv("FOXDEN_DEBUG") != "" {
				fmt.Println("Parse FOXDEN config read from FOXDEN_CONGIV env:", os.Getenv("FOXDEN_CONFIG"))
			}
			viper.SetConfigFile(os.Getenv("FOXDEN_CONFIG"))
		} else {
			cfile = "/nfs/chess/user/chess_chapaas/.foxden.yaml"
			if _, err := os.Stat(cfile); err == nil {
				viper.SetConfigFile(cfile)
				if os.Getenv("FOXDEN_DEBUG") != "" {
					fmt.Println("Parse FOXDEN config:", cfile)
				}
			} else {
				cfile = `\\chesssamba.classe.cornell.edu\user\chess_chapaas\.foxden.yaml`
				if _, err := os.Stat(cfile); err == nil {
					viper.SetConfigFile(cfile)
					if os.Getenv("FOXDEN_DEBUG") != "" {
						fmt.Println("Parse FOXDEN config:", cfile)
					}
				} else {
					msg := "FOXDEN configuration file is not found"
					fmt.Println(msg)
					return config, errors.New(msg)
				}
			}
		}
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		msg := err.Error()
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			msg = fmt.Sprintf("fail to read %s file, error %v", cfile, err)
		} else {
			// Config file was found but another error was produced
			msg = fmt.Sprintf("unable to parse %s, error %v", cfile, err)
		}
		return config, errors.New(msg)
	}
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}
	// set defaults
	if config.LDAP.Expire == 0 {
		config.LDAP.Expire = 3600 // LDAP cache records expire in 1 hour
	}
	if len(config.CHESSMetaData.DataLocationAttributes) == 0 {
		config.CHESSMetaData.DataLocationAttributes = []string{"data_location_raw"}
	}
	if len(config.MetaData.DataLocationAttributes) == 0 {
		config.MetaData.DataLocationAttributes = []string{"data_location_raw"}
	}
	if config.Sync.SleepInterval == 0 {
		config.Sync.SleepInterval = 600 // sync daemon interval in seconds
	}
	if len(config.CHESSMetaData.SpecScanBeamlines) == 0 {
		config.CHESSMetaData.SpecScanBeamlines = []string{"ID1A3", "ID3A", "ID3B", "ID4B", "1A3", "3A", "3B", "4B"}
	}
	return config, nil
}

func Info() string {
	goVersion := runtime.Version()
	tstamp := time.Now()
	return fmt.Sprintf("git={{VERSION}} go=%s date=%s", goVersion, tstamp)
}

/*
// Example of using this configuration module
func main() {
	cfile := "srv.json"
	config, err := ParseConfig(cfile)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	fmt.Printf("Frontend %+v\n", config.Frontend)
	fmt.Printf("Authz %+v\n", config.Authz)
}
*/
