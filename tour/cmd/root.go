package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().Int8VarP(&z, "zap", "z", 0, "if z then printToTxt")
	rootCmd.AddCommand(wordCmd)
	//go run main.go sql struct --username JACMES --password
	//hfutIE2019 --db JACMES --host localhost:1521 --serviceName LocalCummins --table AREA --type oci8 -z 1

	//go run main.go sql struct --username root --password 123456 --db bigdata --host localhost:3306 --table dept --type mysql -z 1
	rootCmd.AddCommand(sqlCmd)

}
