package cmd

import (
	"context"
	pb "gruut-console/services"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of a node.",
	Long:  ``,
}

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
	} else {
		infoLogger.Println("A node is not running.")
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Run = status
}
