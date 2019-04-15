package main

import (
	"os"
	"github.com/spf13/cobra"
	"fmt"
	"time"
)

type json map[string]interface {}

var Provider string

var installCmd = &cobra.Command{
	Use:   "install [--provider=]",
	Short: "Install package",
	Long:  `Install package based on info in the entry point (depends on provider)`,
	Run: func(cmd *cobra.Command, args []string) {
		err := InstallProject(".")
		if(err != nil) {
			fmt.Println(err)
		} 
	},
}
  
var rootCmd = &cobra.Command{
	Use:   "GuPM",
	Short: "GuPM is the Global Universal Project Manager",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test")
	},
}
  
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	start := time.Now()

	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "Provider plugin")

	Execute()
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}