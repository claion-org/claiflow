package verbose

type verbosity = int

const (
	TRACE = verbosity(2)
	DEBUG = verbosity(1)
	INFO  = verbosity(0)
	WARN  = verbosity(-1)
	ERROR = verbosity(-2)
)
