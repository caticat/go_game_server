package petcd

type ServicePrefix string

const (
	Config  ServicePrefix = "/config"
	Service ServicePrefix = "/service"
)
