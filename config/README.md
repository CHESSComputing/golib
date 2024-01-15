# Configuration module
The configuration module for FOXDEN/CHESS services.

### Client configuration
In order to run FOXDEN/CHESS client you should have the following
configuration in your `$HOME/.foxden.yaml` file or supply your file to the
client `--config option`:

```
Services:
  FrontendUrl: https://foxden.url
  DiscoveryUrl: https://foxden-discover.url
  DataManagementUrl: https://foxden-s3.url
  DataBookkeepingUrl: https://foxden-dbs.url
  MetaDataUrl: https://foxden-meta.url
  AuthzUrl: https://foxden-authz.url
```
Please note, you may change the above example according to FOXDEN
setup, e.g. change urls or provide additional port to FOXDEN
services.

### Server configuration
The server configuration is more complex and defined according
to [config.go](config.go) implementation. Here is an example of
FOXDEN local setup:
```
Services:
  FrontendUrl: http://localhost:8344
  DiscoveryUrl: http://localhost:8320
  DataManagementUrl: http://localhost:8340
  DataBookkeepingUrl: http://localhost:8310
  MetaDataUrl: http://localhost:8300
  AuthzUrl: http://localhost:8380
CHESSMetaData:
  SchemaFiles:
    - "schemas/lite.json"
    - "schemas/ID4B.json"
    - "schemas/ID3A.json"
    - "schemas/ID1A3.json"
  SchemaSections: ["User", "Alignment", "DataLocations", "Beam", "Experiment", "Sample"]
  WebSectionKeys:
    User: ["Facility", "Cycle", "PI", "BTR", "Experimenters", "Beamline", "StaffScientist", "BeamlineFundingPartner"]
    Alignment: []
    DataLocations: ["DataLocationRaw", "DataLocationMeta", "DataLocationReduced", "DataLocationScratch", "DataLocationBeamtimeNotes"]
    Beam: ["CESRCondtions"]
    Experiment: ["Detectors","ExperimentType","Technique"]
    Sample: ["SampleType","SampleName","Calibration"] 
  MongoDB:
    DBUri: mongodb://localhost:8230
    DBName: chess
    DBColl: meta
  WebServer:
    Port: 8300
    Verbose: 1
    LogLongFile: true
    GinOptions:
      DisableConsoleColor: true
Kerberos:
  Krb5Conf:  /etc/krb5.conf
  Realm: CLASSE.CORNELL.EDU
Authz:
  DBUri: ./auth.db
  ClientId: xxx
  ClientSecret: xyz
  WebServer:
    Port: 8380
    Verbose: 1
    LogLongFile: true
    GinOptions:
      DisableConsoleColor: true
DataBookkeeping:
  ApiParametersFile: /Users/vk/Work/CHESS/FOXDEN/DataBookkeeping/static/parameters.json
  DBFile: /Users/vk/Work/CHESS/FOXDEN/DataBookkeeping/dbfile
  LexiconFile: /Users/vk/Work/CHESS/FOXDEN/DataBookkeeping/static/lexicon_reader.json
  WebServer:
    Port: 8310
    StaticDir: /Users/vk/Work/CHESS/FOXDEN/DataBookkeeping/static
    Verbose: 1
    LogLongFile: true
    GinOptions:
      DisableConsoleColor: true
DataManagement:
  WebServer:
    Port: 8340
    Verbose: 1
    LogLongFile: true
    GinOptions:
      DisableConsoleColor: true
Discovery:
  MongoDB:
    DBUri: mongodb://localhost:8230
    DBName: chess
    DBColl: meta
  WebServer:
    Port: 8320
    Verbose: 2
    LogLongFile: true
    GinOptions:
      DisableConsoleColor: true
Encryption:
  Cipher: aes
  Secret: bla
Frontend:
  WebServer:
    Port: 8344
    StaticDir: Static
    Verbose: 1
    LogLongFile: true
    GinOptions:
      DisableConsoleColor: true
      Production: true
  OAuth:
    -
      Provider: github
      ClientID: clientid
      ClientSecret: secret
    -
      Provider: google
      ClientID: cleintid
      ClientSecret: secret
      RedirectURL: http://localhost:8344/google/callback
```
