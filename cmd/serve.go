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
	"fmt"
	"geolocation/internal/store/database"
	"geolocation/internal/utils"
	"geolocation/pkg/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run web server on provided port",
	Long: `
	runs web server on provided port,
	the server exposes an endpoint named 'resolve' - 
	which can be accessed as http://localhost:<port>/resolve?ip=123.456.789.487 
 	resolves the provided ip address to JSON object containing Country, City, Latitude & Longitude etc.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// initialise a common logger
		logger := utils.GetLogger().WithFields(logrus.Fields{
			"command": "serve",
			"port":    cmd.Flag("port").Value,
		})

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

		// initialise location service
		locationSrv := service.NewLocationService(conn)

		// bind handler to location service
		http.HandleFunc("/resolve", locationSrv.Resolve)

		// init & start the server
		htttpSrvErrStream := make(chan error)
		s := &http.Server{
			Addr:    fmt.Sprintf(":%s", cmd.Flag("port").Value.String()),
			Handler: http.DefaultServeMux,
		}
		go func() {
			logger.Debug("listening")
			err := s.ListenAndServe()
			if err != nil {
				htttpSrvErrStream <- err
			}
		}()

		// cleanups & graceful shutdown
		quitSignalStream := make(chan os.Signal, 1)
		signal.Notify(quitSignalStream, os.Interrupt, syscall.SIGABRT, syscall.SIGTERM)

		select {
		case s := <-quitSignalStream:
			signal.Stop(quitSignalStream)
			logger.WithFields(logrus.Fields{"signal": s}).Debug(s.String() + "ed")
			return
		case err := <-htttpSrvErrStream:
			if err != nil {
				logger.WithFields(logrus.Fields{"err": err}).Debug("listening to port failed")
				return
			}
		case <-cmd.Context().Done():
			s.Shutdown(cmd.Context())
			logger.Debug("shutdown complete")
		}
		logger.Info("exiting")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int32P("port", "p", 8080, "web server port to listen on")
}
