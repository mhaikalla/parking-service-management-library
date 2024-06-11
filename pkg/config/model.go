package config

// IConfig interface to provide configuration
type IConfig interface {
	// GetConfig arg is string uri to lookout a configuration return an error
	GetConfig(string) error
	Config() map[string]map[string]interface{}
}
