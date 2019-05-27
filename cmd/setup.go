package cmd

import (
	"fmt"
	"context"
	pb "gruut-console/services"
	. "github.com/logrusorgru/aurora"

	"github.com/gen2brain/beeep"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:     "setup",
	Short:   "Setup the node",
	Long:    ``,
	PostRun: postRun,
}

var password string

func setup(cmd *cobra.Command, args []string) {
	fmt.Println(Cyan("[INFO]"),  "It may take a while...")

	connOption := grpc.WithInsecure()
	conn, err := grpc.Dial(address, connOption)
	if err != nil {
		errorLogger.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewGruutAdminServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := client.Setup(ctx, &pb.ReqSetup{Password: password})

	if err != nil {
		beepError := beeep.Notify("Error", err.Error(), "assets/warning.png")
		if beepError != nil {
			panic(beepError)
		}
	}

	if resp.Success == true {
		infoLogger.Println("The Setup has successfully completed.")
	}
}

func postRun(cmd *cobra.Command, args []string) {
	err := beeep.Notify("Title", "Message body", "assets/information.png")
	if err != nil {
		panic(err)
	}
}

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.Run = setup

	setupCmd.PersistentFlags().StringVar(&password, "password", "", "")
}
