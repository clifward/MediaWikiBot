package config

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var configFile string

// Config holds the configuration file
type Config struct {
	Mode                 string               `yaml:"mode" mapstructure:"mode"`
	ChanID               int                  `yaml:"chan_id" mapstructure:"chan_id"`
	WikiClientConfig     WikiClientConfig     `yaml:"wiki" mapstructure:"wiki"`
	TelegramClientConfig TelegramClientConfig `yaml:"telegram" mapstructure:"telegram"`
}

// WikiClientConfig holds the configuration for the wiki client
type WikiClientConfig struct {
	URL          string        `yaml:"url" mapstructure:"url"`
	UserAgent    string        `yaml:"user_agent" mapstructure:"user_agent"`
	Username     string        `yaml:"username" mapstructure:"username"`
	Password     string        `yaml:"password" mapstructure:"password"`
	RefreshDelay time.Duration `yaml:"refresh_delay" mapstructure:"refresh_delay"`
}

// TelegramClientConfig holds the configuration for the Telegram client
type TelegramClientConfig struct {
	ChanID   int    `yaml:"chan_id" mapstructure:"chan_id"`
	Admin    int    `yaml:"admin" mapstructure:"admin"`
	Name     string `yaml:"bot_name" mapstructure:"bot_name"`
	BotToken string `yaml:"bot_token" mapstructure:"bot_token"`
}

// RegisterFlags overwrite the configuration with parameter passed with flags
func (cfg *Config) RegisterFlags(flags *pflag.FlagSet) {
	flags.StringVar(&configFile, "config", "config.yaml", "configuration file to use")

	flags.StringVar(&cfg.Mode, "mode", "job", "'job' : 1 execution, 'loop_job' : loop and execute the job every refresh_delay")
}

// RegisterConfigFile loads the new config file is another than default is specified
func (cfg *Config) RegisterConfigFile() error {
	v := viper.New()
	if configFile != "" {
		v.SetConfigFile(configFile)
		err := v.ReadInConfig()
		if err != nil {
			return err
		}
		err = v.Unmarshal(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}
