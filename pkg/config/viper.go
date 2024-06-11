package config

import (
	"github.com/spf13/viper"
)

// ViperLocalProvider using viper to load config from local file
type ViperLocalProvider struct {
	config map[string]map[string]interface{}
}

// parseDir return directory path from uri
// func parseDir(uri string) string {
// 	splitted := strings.Split(uri, "/")
// 	n := len(splitted) - 1
// 	splitted[n] = ""
// 	return strings.Join(splitted, "/")
// }

// parseFile return filename from uri
// func parseFile(uri string) string {
// 	splitted := strings.Split(uri, "/")
// 	n := len(splitted) - 1
// 	return splitted[n]
// }

// explodeFilename split filename and its extension
// func explodeFilename(filename string) []string {
// 	splitted := strings.Split(filename, ".")
// 	n := len(splitted) - 1
// 	ext := splitted[n]
// 	splitted[n] = ""
// 	return []string{strings.Join(splitted, "."), ext}
// }

// GetConfig read and set config from uri as arg
func (p *ViperLocalProvider) GetConfig(uri string) error {
	viper.SetConfigFile(uri)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	for k, v := range viper.AllSettings() {
		p.config[k] = v.(map[string]interface{})
	}

	return nil
}

// Config return config readed
func (p *ViperLocalProvider) Config() map[string]map[string]interface{} {
	return p.config
}

// NewViperLocalProvider create a new viper local provider
func NewViperLocalProvider() IConfig {
	return &ViperLocalProvider{config: map[string]map[string]interface{}{}}
}
