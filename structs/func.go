package structs

// Func represents an executable function and its metadata
type Func struct {
	ID         []byte
	Runtime    Runtime
	OSPackages []string
	Packages   []string
	Code       string
	Author     string
	Name       string
	Version    string
}

// FuncExecution represents a single funciton invocation
type FuncExecution struct {
	ID       []byte
	FuncID   []byte
	Ctx      string // JSON encoded string
	StdOut   string
	StdErr   string
	Start    uint32
	End      uint32
	Duration uint32
}

// Runtime is an enum indicating supported func runtimes
type Runtime uint8

var Runtimes = map[string]Runtime{
	"py": RuntimePython,
	"js": RuntimeJavascript,
	"r":  RuntimeR,
	"go": RuntimeGo,
}

// Runtime enumeration
const (
	RuntimePython = iota
	RuntimeJavascript
	RuntimeR
	RuntimeGo
)

// String returns Runtime string representation
func (r Runtime) String() string {
	for k, v := range Runtimes {
		if v == r {
			return k
		}
	}
	return ""
}

// RuntimeFromString converts a string into a matching Runtime
func RuntimeFromString(s string) Runtime {
	return runtimes[s]
}
