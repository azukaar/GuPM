package main

import (
	"os"
	// "github.com/robertkrimen/otto"
	"github.com/spf13/cobra"
	"fmt"
	"time"
)

var Provider string

var installCmd = &cobra.Command{
	Use:   "install [--provider=]",
	Short: "Install package",
	Long:  `Install package based on info in the entry point (depends on provider)`,
	Run: func(cmd *cobra.Command, args []string) {
		if(Provider != "") {
			fmt.Println("Reading provider config for", Provider);
			providerJson := readJSON("plugins/provider-" + Provider + "/gupm.json")
			fmt.Println("Initialisation OK for", Provider);
			fmt.Println(providerJson["config"])
		}
	},
}
  
  var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
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