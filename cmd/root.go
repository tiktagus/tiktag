/*
Copyright © 2022 Atman An <twinsant@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

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
		// fmt.Println(fHash)
		url := searchAsset(fHash)
		if strings.Compare(url, "") == 1 {
			fmt.Println("We found this asset at,")
			fmt.Println(url)
			return
		}

		id := NextID()
		// fmt.Println(id)

		fmt.Printf("Tik...Tag...")
		url, ext := publishFile(id, fn, contentType)
		fmt.Printf("your asset is successfully hosted at,\n%s\n", url)
		// https://s3.tikoly.com/tiktag/myfilename.jpg

		ttasset := TTAsset{
			ttid:     id,
			hash:     fHash,
			filename: fn,
			fileext:  ext,
			url:      url,
		}
		saveAsset(ttasset)
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
	os.Setenv("LOG_LEVEL", "warn") // Clear immudb log

	home := os.Getenv("HOME")
	name := "config"
	ftype := "yaml"
	dir := ".tiktag"
	configFile := fmt.Sprintf("%s.%s", path.Join(home, dir, name), ftype)
	fmt.Println(configFile)
	if _, err := os.Stat(configFile); err != nil {
		fmt.Println("The first time you run tiktag, will create a config file step by step:")

		// Create .tiktag if not exist

		// Set configs

		// Save
	}

	viper.AutomaticEnv()

	viper.SetConfigName(name)                         // name of config file (without extension)
	viper.SetConfigType(ftype)                        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(fmt.Sprintf("$HOME/%s", dir)) // call multiple times to add many search paths
	viper.ReadInConfig()
}
