// Copyright © 2019 Gruut network <contact@gruut.net>
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
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var errorLogger = log.New(os.Stdout, "[Status Request Error] ", log.Lshortfile)
var infoLogger = log.New(os.Stdout, "[Status Request Info] ", 0)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of a node.",
	Long:  ``,
}

var address string

func status(cmd *cobra.Command, args []string) {
	connOption := grpc.WithInsecure()
	conn, err := grpc.Dial(address, connOption)
	if err != nil {
		errorLogger.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewGruutAdminServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.CheckStatus(ctx, &pb.ReqStatus{})
	if err != nil {
		errorLogger.Fatalf(err.Error())
	}

	if resp.Alive == true {
		infoLogger.Println("A node is running.")
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Run = status

	statusCmd.PersistentFlags().StringVar(&address, "address", "10.10.10.200:59002", "A node address (default is localhost:59001)")
}
