package cmd

import (
	"github.com/spf13/cobra"
	"github.com/alexxxPopa/courses/conf"
	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"

	"github.com/alexxxPopa/courses/storage/sql"
)

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Long: "Migrate database strucutures. This will create new tables and add missing collumns and indexes.",
	Run: func (cmd *cobra.Command, args[] string) {
		config, _ := conf.LoadConfig(cmd)

		if err := viper.Unmarshal(&config); err != nil {
			logrus.Fatalf("Unable to decode into struct")
		}
		conn, err := sql.Connect(config);
		defer conn.Close()

		if err != nil {
			logrus.Fatal(err)
		}

		if err = conn.Migrate(); err != nil {
			logrus.Fatal(err)
		}
	},

}

func init() {
	RootCmd.AddCommand(migrateCmd)
}