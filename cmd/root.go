/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var table string
var chain string
var targetRule string
var unique bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iptables_sync",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&unique, "unique", "u", false, "unique iptables rule")
	rootCmd.PersistentFlags().StringVarP(&table, "table", "t", "filter", "iptables table name")
	rootCmd.PersistentFlags().StringVarP(&chain, "chain", "c", "INPUT", "iptables chain name")
	rootCmd.PersistentFlags().StringVarP(&targetRule, "target", "r", "-A INPUT -j REJECT --reject-with icmp-host-prohibited", "iptables target rule")
}
