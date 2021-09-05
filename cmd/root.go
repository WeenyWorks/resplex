package cmd
import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/WeenyWorks/resplex/cmd/resplex"
	"github.com/WeenyWorks/resplex/cmd/resplexc"
	"github.com/WeenyWorks/resplex/cmd/visit"

	"fmt"
	"os"
)

var cfgFile string

var rootCmd = &cobra.Command{
  Use:   "resplex",
  Short: "ResPLEX is a reverse proxy designed for edge environment",
  Long: `Inspired by frp, ResPLEX let you visit services on edge easily.
         Unlike frp, ResPLEX doesn't consume your port resource.
                Complete documentation is available at http://edgefusion.io`,
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(resplex.ServeCMD, resplexc.RegisterCMD, visit.VisitCMD)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.resplex.toml)")
}


func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".resplex")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
