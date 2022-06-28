/*
Copyright © 2022 Manish Sharma bhardwaz007@yahoo.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"geolocation/internal/utils"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "geolocation",
	Long: `
	geolocation is a CLI App which does the task of resolving 
	an IP Address to Country, City, Latitude & Longitude.
	It does so after ingesting location data from a '*.csv' file .
	
	For doing that it exposes 2 commands:

	#1. ingest
	#2. serve 

	For more details run --help on commands [ingest, serve]
	`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		utils.GetLogger().WithFields(logrus.Fields{"err": err}).Error("rootCmd.Execute() failed.")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in .env file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		utils.GetLogger().WithFields(logrus.Fields{"err": err}).Error("ReadInConfig() failed.")
	} else {
		utils.GetLogger().WithFields(logrus.Fields{"config": viper.ConfigFileUsed()}).Debug("config loaded.")
	}
}
