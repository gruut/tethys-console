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
	pb "gruut-console/services"
	"time"

	"google.golang.org/grpc"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the node",
	Long: ``,
}

var password string

func setup(cmd *cobra.Command, args []string) {
	connOption := grpc.WithInsecure()
	conn, err := grpc.Dial(address, connOption)
	if err != nil {
		errorLogger.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewGruutAdminServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Setup(ctx, &pb.ReqSetup{Password: password})
	if err != nil {
		errorLogger.Fatalf(err.Error())
	}

	if resp.Success == true {
		infoLogger.Println("The Setup has successfully completed.")
	}
}

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.Run = setup
	
	statusCmd.PersistentFlags().StringVar(&password, "password", "", "")
}
