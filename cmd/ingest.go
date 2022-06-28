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
	"geolocation/internal/store/database"
	"geolocation/internal/utils"
	"geolocation/pkg/service"
	"io/fs"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ingestCmd represents the ingest command
var ingestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "reads a csv file named '*.csv' from the mounted location.",
	Long: `
	reads a csv file named '*.csv' from the mounted location.
	(for local debugging the file has to be present @ root)
	After the read is complete, the data will be loaded in database and 
	a detailed output will be presented with following details:
	
	#1. total time taken to parse & load the data in millisecond,
	#2. number of entries accepted, and
	#3. number of entries discarded.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// initialise a common logger
		logger := utils.GetLogger().WithFields(logrus.Fields{
			"command": "ingest",
		})

		// make sure the provided file (or default) is a csv file, or fail fast
		file := cmd.Flag("file").Value.String()
		if filepath.Ext(file) != ".csv" {
			logger.Error("invalid file format, only csv is supported")
			return
		}

		// open file or fail fast
		r, err := os.OpenFile(file, os.O_RDONLY, fs.FileMode(os.O_RDONLY))
		if err != nil {
			logger.WithFields(logrus.Fields{"err": err}).Error("OpenFile() failed")
			return
		}
		defer r.Close()

		// load db config or fail fast
		dbCfg, err := utils.GetDBCfg()
		if err != nil {
			logger.WithFields(logrus.Fields{"err": err}).Error("GetDBCfg() failed")
			return
		}

		// create database connection or fail fast
		conn, err := database.New(*dbCfg)
		if err != nil {
			logger.WithFields(logrus.Fields{"err": err}).Error("NewConnection() failed")
			return
		}
		defer conn.Close()

		// run migration
		if err := conn.Migrate(cmd.Context()); err != nil {
			logger.WithFields(logrus.Fields{"err": err}).Error("Migrate() failed")
			return
		}

		// initialise ingestor service
		ingestorSrvc := service.NewCSVIngestor(conn, r)

		// read, sanitise & ingest all valid locations
		logger.Debug("ingestion in progress ...")
		stat, _, err := ingestorSrvc.Ingest(cmd.Context())
		if err != nil {
			logger.WithFields(logrus.Fields{"err": err}).Error("Ingest() failed")
			return
		}

		// display stats
		logger.WithFields(logrus.Fields{
			"accepted":        stat.Accepted,
			"discarded":       stat.Discarded,
			"spent (ms)":      stat.TimeSpent.Milliseconds(),
			"total record(s)": stat.Accepted + stat.Discarded,
		}).Info("ingestion complete.")
	},
}

func init() {
	rootCmd.AddCommand(ingestCmd)
	ingestCmd.Flags().StringP("file", "f", "data_dump.csv", "csv file name to ingest data from")
}
