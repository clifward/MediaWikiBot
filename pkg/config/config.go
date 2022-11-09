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
	flags.StringVar(&configFile, "config", "", "configuration file to use")

	flags.StringVar(&cfg.Mode, "mode", "job", "Use 'job' for one execution or 'loop_job' if you want a loop and and execute the job every refresh_delay")

	// WikiClientConfig
	flags.StringVar(&cfg.WikiClientConfig.URL, "wiki.url", "", "URL of the wiki API")
	flags.StringVar(&cfg.WikiClientConfig.UserAgent, "wiki.user_agent", "wikibot-golang-client", "User agent to use for the wiki API")
	flags.StringVar(&cfg.WikiClientConfig.Username, "wiki.username", "", "Username to use for the wiki API")
	flags.StringVar(&cfg.WikiClientConfig.Password, "wiki.password", "", "Password to use for the wiki API")
	flags.DurationVar(&cfg.WikiClientConfig.RefreshDelay, "refresh_delay", 60*time.Minute, "Delay between 2 refreshes of the wiki API")

	// TelegramClientConfig
	flags.StringVar(&cfg.TelegramClientConfig.Name, "telegram.bot_name", "", "Name of the Telegram bot")
	flags.StringVar(&cfg.TelegramClientConfig.BotToken, "telegram.bot_token", "", "Token of the Telegram bot")
	flags.IntVar(&cfg.TelegramClientConfig.Admin, "telegram.admin", 0, "Telegram admin ID")
	flags.IntVar(&cfg.TelegramClientConfig.ChanID, "telegram.chan_id", 0, "Telegram channel ID")

}

// RegisterConfigFile loads the new config file if another than default is specified
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
