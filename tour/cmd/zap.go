package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"src/go-programming-tour-book/tour/zap"
)

var z int8

var zapCmd = &cobra.Command{
	Use:   "zap",
	Short: "是否写入日志",
	Long:  "是否写入日志",
	Run: func(cmd *cobra.Command, args []string) {
		//if z==1 {
		if cmd.Use == "word" {

		}
		zap.WriteLog(fmt.Sprintf("go run main.go -s %s -m %d", str, mode))
		zap.WriteLog(fmt.Sprintf("输出结果 : %s ", content))
		//}
	},
}

func init() {
	//zapCmd.Flags().Int8VarP(&z,"zap","z",0,"是否写入日志 0：否 1：是")

	//zapCmd.AddCommand(wordCmd)
	//zapCmd.AddCommand(sqlCmd)

}
