package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
	"fmt"
	"os"
	"github.com/alexxxPopa/courses/conf"
	"github.com/alexxxPopa/courses/api"
	"github.com/alexxxPopa/courses/storage/sql"
	"github.com/alexxxPopa/courses/services/impl"
)

var RootCmd = &cobra.Command{
	Use:   "courses",
	Short: "A brief description of your application",
	Long:  `A longer description.`,

	Run: func(cmd *cobra.Command, args [] string) {
		config, err := conf.LoadConfig(cmd)
		if err != nil {
			logrus.Fatal(err)
		}
		conn, err := sql.Connect(config);
		if err != nil {
			logrus.WithError(err).Fatal("connection to database failed")
		}
		stripe := impl.Setup(config)
		api := api.Create(config, conn, stripe)
		if err := api.ListenAndServe(fmt.Sprintf("%v:%v", config.SERVER.Host, config.SERVER.Port)); err != nil {
			logrus.WithError(err).Fatal("http server listen failed")
		}
	},

}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&conf.ConfigFile, "config", "", "config file (default is $HOME/.auth.yaml)")

	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
