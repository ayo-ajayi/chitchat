package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/ayo-ajayi/chitchat/cmd/chat"
	"github.com/ayo-ajayi/chitchat/cmd/history"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var (
	version = "1.0.0"
)
var rootCmd = &cobra.Command{
	Use:     "chitchat",
	Short:   "ChatGPT CLI",
	Long:    ``,
	Version: version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(chat.ChatCmd)
	rootCmd.AddCommand(history.HistoryCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.chitchat.yaml)")
}
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err == nil {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".chitchat")
		viper.AutomaticEnv()
		viper.ReadInConfig()
	}

}
