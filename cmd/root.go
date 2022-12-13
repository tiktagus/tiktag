/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tiktag [file to upload]",
	Short: "A command-line tool for preparing images for blog post or sharing.",
	Long:  `Upload a photo and get its S3 URL back as a response, for use in Markdown for publishing.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fn := args[0]

		fHash, contentType := getFileHash(fn)
		fmt.Println(fHash)

		id := getFileId()
		fmt.Println(id)

		info := publishFile(fn, contentType)
		fmt.Println(info)
		// https://s3.tikoly.com/tiktag/myfilename.jpg

		saveAsset()
	},
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tiktag.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	viper.SetConfigName("config")        // name of config file (without extension)
	viper.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/tiktag/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.tiktag") // call multiple times to add many search paths
	viper.AddConfigPath(".")             // optionally look for config in the working directory
	viper.ReadInConfig()
}
