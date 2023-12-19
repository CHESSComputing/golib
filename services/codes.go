package server

const (
	OK                    = 0
	GenericError          = iota + 100 // generic error
	DatabaseError                      // 101 database error
	TransactionError                   // 102 transaction error
	QueryError                         // 103 query error
	RowsScanError                      // 104 row scan error
	SessionError                       // 105 db session error
	CommitError                        // 106 db commit error
	ParseError                         // 107 parser error
	LoadError                          // 108 loading error, e.g. load template
	GetIDError                         // 109 get id db error
	InsertError                        // 110 db insert error
	UpdateError                        // 111 update error
	LastInsertError                    // 112 db last insert error
	ValidateError                      // 113 validation error
	PatternError                       // 114 pattern error
	DecodeError                        // 115 decode error
	EncodeError                        // 116 encode error
	ContentTypeError                   // 117 content type error
	ParametersError                    // 118 parameters error
	NotImplementedApiCode              // 119 not implemented API error
	ReaderError                        // 120 io reader error
	WriterError                        // 121 io writer error
	UnmarshalError                     // 122 json unmarshal error
	MarshalError                       // 123 marshal error
	HttpRequestError                   // 124 HTTP request error
	RemoveError                        // 125 remove error
	BindError                          // 126 bind error
	SchemaError                        // 127 schema error
)
