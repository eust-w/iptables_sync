/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/coreos/go-iptables/iptables"
	"github.com/eust-w/iptables_sync/ctrl"

	"github.com/spf13/cobra"
)

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("insert called")
		ipt, err := iptables.New()
		if err != nil {
			fmt.Println("create iptables: %s", err)
			return
		}
		rule := []string{"-s", "198.189.1.34", "-d", "198.189.1.2", "-j", "ACCEPT"}
		if unique {
			err := ctrl.UniqueInsertIptablesRuleBeforeTargetRule(ipt, table, chain, targetRule, rule)
			if err != nil {
				fmt.Println("insert rule: %s", err)
			}
		} else {
			err := ctrl.InsertIptablesRuleBeforeTargetRule(ipt, table, chain, targetRule, rule)
			if err != nil {
				fmt.Println("insert rule: %s", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(insertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// insertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// insertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
