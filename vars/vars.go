package vars

const (
	DefaultPollInterval          = 30
	DefaultMaxConcurrentIndexers = 2
	DefaultPollEnabled           = true
	DefaultVcs                   = "git"
	DefaultBaseUrl               = "{url}/blob/master/{path}{anchor}"
	DefaultAnchor                = "#L{line}"
	PageStep                     = 5
	SearchNum                    = 25
	Source                       = "gshark"
)

var (
	REPO_PATH    string
	MAX_INDEXERS int

	HTTP_HOST string
	HTTP_PORT int

	MAX_Concurrency_REPOS int

	DEBUG_MODE bool

	PAGE_SIZE = 10
	API_TOKEN string
)
