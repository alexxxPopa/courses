package conf

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"fmt"
	"github.com/sirupsen/logrus"
)

var ConfigFile string

type Config struct {
	DB struct {
		Driver string `json:"driver"`
		Url    string    `json:"url"`
	}

	SERVER struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	//TODO add Stripe configuration
}

func LoadConfig(cmd *cobra.Command) (*Config, error) {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return nil, err
	}

	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath("./") // adding home directory as first search path
		viper.SetConfigName("config")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	config := new(Config)

	if err := viper.Unmarshal(&config); err != nil {
		logrus.Fatalf("Unable to decode into struct")
	}

	fmt.Println(config)

	return config, nil


}
