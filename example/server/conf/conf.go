package conf

// "github.com/go-yaml/yaml"

const (
	ChaRecvLen = 100
	FileConfig = "server.yaml"
)

type ConfServer struct {
	Port int `yaml:"port"`
}

func (t ConfServer) New() *ConfServer {
	return &t
}

func (t *ConfServer) Init() {
	// f, err := ioutil.ReadFile(FileConfig)
	// if err != nil {
	// 	log.Fatal("ioutil.ReadFile failed,error:", err)
	// }

	// err = yaml.Unmarshal(f, t)
	// if err != nil {
	// 	log.Fatal("yaml.Unmarshal failed,error:", err)
	// }
}

func (t *ConfServer) GetPort() int { return t.Port }
