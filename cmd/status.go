package cmd

import (
	"context"
	pb "tethys-console/services"
	"time"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of a node.",
	Long:  ``,
}

func status(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, client, err := beforeExecute(&ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	resp, err := client.CheckStatus(ctx, &pb.ReqStatus{})

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
