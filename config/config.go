package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

// TrustedUser represents single trusted user information
type TrustedUser struct {
	User string `json:"user"`
	IP   string `json:"ip"`
	MAC  string `json:"mac"`
}

// Globug represents globus information
type Globus struct {
	ClientID     string `mapstructure:"client_id"`     // client id
	ClientSecret string `mapstructure:"client_secret"` // client secret
	TransferURL  string `mapstructure:"transfer_url"`  // globus transfer url
	AuthURL      string `mapstructure:"auth_url"`      // globus auth url
	OriginID     string `mapstructure:"origin_id"`     // globus origin ID for CHESS Raw data collection
}

// LDAP attributes
type LDAP struct {
	URL      string `mapstructure:"url"`      // ldap url
	BaseDN   string `mapstructure:"baseDN`    // ldap baseDN
	Login    string `mapstructure:"login"`    // LDAP login to use
	Password string `mapstructure:"password"` // LDAP password to use
	Expire   int    `mapstructure:"expire"`   // LDAP cache record expire (in seconds)
}

// DOI attributes
type DOI struct {
	DocumentDir      string           `mapstructure:"DocumentDir"` // doi document directory
	Provider         string           `mapstructure:"Provider"`    // doi provider, e.g. Zenodo or MaterialsCommons
	ProjectName      string           `mapstructure:"ProjectName"` // name of the project (only valid for MaterialsCommons)
	Zenodo           Zenodo           `mapstructure:"Zenodo"`
	Datacite         Datacite         `mapstructure:"Datacite"`
	MaterialsCommons MaterialsCommons `mapstructure:"MaterialsCommons"`
	WebServer
}

// DID structure
type DID struct {
	Attributes string `mapstructure:attributes` // did attributes, comma separated, default beamline,btr,cycle,sample
	Separator  string `mapstructure:separator`  // did separator, default "/"
	Divider    string `mapstructure:divider`    // did key-value divider, default "="
}

// Embed structure
type Embed struct {
	DocDb string `mapstructure:"DocDb"`
	SqlDb string `mapstructure:"SqlDb"`
}

// QL structure
type QL struct {
	ServiceMapFile string `mapstructure:ServiceMapFile` // service map file name
	Separator      string `mapstructure:separator`      // ql separator, default ":"
	Verbose        int    `mapstructure:verbose`        // verbosity level
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
	Name     string `json:"name"`
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
	Accept   string `json:"accept"`
}

// MLBackend represents ML backend engine
type MLBackend struct {
	Name string  `json:"name"` // ML backend name, e.g. TFaaS
	Type string  `json:"type"` // ML backebd type, e.g. TensorFlow
	URI  string  `json:"uri"`  // ML backend URI, e.g. http://localhost:port
	Apis []MLApi // ML APIs
}

// ML defines ML configuration options
type ML struct {
	MLBackends []MLBackend `json:"MLBackends"` // ML backends
	StorageDir string      `json:"StorageDir"`
}

// Kerberos defines kerberos optinos
type Kerberos struct {
	Krb5Conf string `mapstructure:Krb5Conf`
	Keytab   string `mapstructure:Keytab`
	Realm    string `mapstructure:Realm`
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
	LimiterPeriod   string   `mapstructure:"Rate"`      // limiter rate value
	LimiterHeader   string   `json:"limiter_header"`    // limiter header to use
	LimiterSkipList []string `json:"limiter_skip_list"` // limiter skip list
	MetricsPrefix   string   `json:"metrics_prefix"`    // metrics prefix used for Prometheus

	// etag options
	Etag         string `json:"etag"`          // etag value to use for ETag generation
	CacheControl string `json:"cache_control"` // Cache-Control value, e.g. max-age=300

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
	WebServer `mapstructure:"WebServer"`

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

// Discovery represents discovery service configuration
type Discovery struct {
	WebServer  `mapstructure:"WebServer"`
	MongoDB    `mapstructure:"MongoDB"`
	Encryption `mapstructure:"Encryption"`
}

// MetaData represents metadata service configuration
type MetaData struct {
	WebServer              `mapstructure:"WebServer"`
	MongoDB                `mapstructure:"MongoDB"`
	LexiconFile            string   `mapstructure:"LexiconFile"`            // lexicon file
	DataLocationAttributes []string `mapstructure:"DataLocationAttributes"` // data location attributes to use
}

// CHESSMetaData represents CHESS MetaData configuration
type CHESSMetaData struct {
	WebServer              `mapstructure:"WebServer"`
	MongoDB                `mapstructure:"MongoDB"`
	SchemaRenewInterval    int      `json:"SchemaRenewInterval"`            // schema renew interval
	WebSectionsFile        string   `mapstructure:"WebSectionsFile"`        // file for web form sections
	LexiconFile            string   `mapstructure:"LexiconFile"`            // lexicon file
	TestMode               bool     `mapstructure:TestMode`                 // test mode
	DataLocationAttributes []string `mapstructure:"DataLocationAttributes"` // data location attributes to use
	SchemaFiles            []string `json:"SchemaFiles"`                    // schema files
	OrderedSections        []string `mapstructure:OrderedSections`          // ordered sections for web UI
	SkipKeys               []string `mapstructure:SkipKeys`                 // keys to skip for web forms
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
	Url         string `mapstructure:"Url"`
	AccessToken string `mapstructure:"AccessToken"`
	ProjectName string `mapstructure:"ProjectName"` // name of the project (only valid for MaterialsCommons)
}

// Datacite represents Datacite service configuration
type Datacite struct {
	Url          string `mapstructure:"Url"`
	AccessToken  string `mapstructure:"AccessToken"`
	AccessKey    string `mapstructure:"AccessKey"`
	AccessSecret string `mapstructure:"AccessSecret"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
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

	TestMode     bool   `mapstructure:TestMode`  // test mode
	CheckLDAP    bool   `mapstructure:CheckLDAP` // check users scope in LDAP
	DBFile       string `mapstructure:"DBFile"`
	ClientID     string `mapstructure:"ClientId"`
	ClientSecret string `mapstructure:"ClientSecret"`
	Domain       string `mapstructure:"Domain"`
	TokenExpires int64  `mapstructure:TokenExpires` // expiration of token
}

// Services represents services structure
type Services struct {
	FrontendURL        string `mapstructure:"FrontendUrl"`
	DiscoveryURL       string `mapstructure:"DiscoveryUrl"`
	MetaDataURL        string `mapstructure:"MetaDataUrl"`
	MLHubURL           string `mapstructure:"MLHubUrl"`
	DataManagementURL  string `mapstructure:"DataManagementUrl"`
	DataBookkeepingURL string `mapstructure:"DataBookkeepingUrl"`
	AuthzURL           string `mapstructure:"AuthzUrl"`
	SpecScansURL       string `mapstructure:"SpecScansUrl"`
	PublicationURL     string `mapstructure:"PublicationUrl"`
	CHAPBookURL        string `mapstructure:"CHAPBookUrl"`
}

// SrvConfig represents configuration structure
type SrvConfig struct {
	S3              `mapstructure:"S3"`
	QL              `mapstructure:"QL"`
	DID             `mapstructure:"DID"`
	LDAP            `mapstructure:LDAP`
	Frontend        `mapstructure:"Frontend"`
	Discovery       `mapstructure:"Discovery"`
	MetaData        `mapstructure:"MetaData"`
	MLHub           `mapstructure:"MLHub"`
	DataManagement  `mapstructure:"DataManagement"`
	DataBookkeeping `mapstructure:"DataBookkeeping"`
	Authz           `mapstructure:"Authz"`
	Kerberos        `mapstructure:"Kerberos"`
	Services        `mapstructure:"Services"`
	Encryption      `mapstructure:"Encryption"`
	CHESSMetaData   `mapstructure:"CHESSMetaData"`
	OreCastMetaData `mapstructure:"OreCastMetaData"`
	SpecScans       `mapstructure:"SpecScansService"`
	Publication     `mapstructure:"PublicationService"`
	Globus          `mapstructure:"Globus"`
	DOI             `mapstructure:"DOI"`
	Embed           `mapstructure:"Embed"`

	TrustedUsers []TrustedUser `mapstructure:"TrustedUsers"`
}

func (c *SrvConfig) String() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Println("ERROR:", err)
		return fmt.Sprintf("%s", string(data))
	}
	return string(data)
}

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
	return config, nil
}

/*
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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

// Config represnets configuration instance
var Config *SrvConfig

func Info() string {
	goVersion := runtime.Version()
	tstamp := time.Now()
	return fmt.Sprintf("git={{VERSION}} go=%s date=%s", goVersion, tstamp)
}

func Init() {
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	var config string
	flag.StringVar(&config, "config", "", "server config file")
	flag.Parse()
	if version {
		fmt.Println("server version:", Info())
		return
	}
	if config == "" {
		// check env variable
		config = os.Getenv("FOXDEN_CONFIG")
	}
	if os.Getenv("FOXDEN_DEBUG") != "" {
		log.Println("Read FOXDEN config:", config)
	}
	oConfig, err := ParseConfig(config)
	if err != nil {
		log.Fatal("ERROR", err)
	}
	Config = &oConfig
}
