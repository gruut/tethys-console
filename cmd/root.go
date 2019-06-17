// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	aurora "github.com/logrusorgru/aurora"
	homedir "github.com/mitchellh/go-homedir"
	pb "tethys-console/services/grpc_admin"

	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var cfgFile string
var address string

var errorLogger = log.New(os.Stdout, fmt.Sprint("[", aurora.Red("ERROR"), "] "), log.Lshortfile)
var infoLogger = log.New(os.Stdout, fmt.Sprint("[", aurora.Cyan("INFO"), "] "), 0)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tethys-console",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

type Response interface {
	Info() string
	Error() string
	Success() bool
}

func beforeExecute(ctx *context.Context) (*grpc.ClientConn, pb.TethysAdminServiceClient, error) {
	connOption := grpc.WithInsecure()
	conn, err := grpc.Dial(address, connOption)
	if err != nil {
		errorLogger.Fatalf("Failed to dial: %v", err)
		return nil, nil, err
	}

	client := pb.NewTethysAdminServiceClient(conn)

	return conn, client, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func afterExecute(resp Response, err error) {
	if err != nil {
		errorLogger.Println(err.Error())

		beepError := beeep.Notify("Error", err.Error(), "assets/warning.png")
		if beepError != nil {
			errorLogger.Fatal(beepError)
		}

		return
	}

	if resp.Success() == true {
		err := beeep.Notify("[SUCCESS]", resp.Info(), "assets/information.png")
		if err != nil {
			panic(err)
		}
	} else {
		errorLogger.Println(resp.Error())

		err := beeep.Notify("[FAILED]", resp.Error(), "assets/information.png")
		if err != nil {
			errorLogger.Fatal(err)
		}
	}	
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tethys-console.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tethys-console" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tethys-console")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		address = viper.Get("address").(string)
	}
}
