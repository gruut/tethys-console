package cmd

import (
	"context"
	pb "gruut-console/services"

	"github.com/gen2brain/beeep"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:     "setup",
	Short:   "Setup the node",
	Long:    ``,
}

var port string

func setup(cmd *cobra.Command, args []string) {
	infoLogger.Println("It may take a while...")

	connOption := grpc.WithInsecure()
	conn, err := grpc.Dial(address, connOption)
	if err != nil {
		errorLogger.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewGruutAdminServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := client.SetupKey(ctx, &pb.ReqSetupKey{SetupPort: port})

	if err != nil {
		beepError := beeep.Notify("Error", err.Error(), "assets/warning.png")
		if beepError != nil {
			panic(beepError)
		}
	}

	if resp.Success == true {
		err := beeep.Notify("[SETUP KEY]", "The Setup has successfully completed.", "assets/information.png")
		if err != nil {
			panic(err)
		}
	} else {
		infoLogger.Println(resp)
		
		err := beeep.Notify("[SETUP KEY]", "Failed to setup.", "assets/information.png")
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.Run = setup

	setupCmd.Flags().StringVarP(&port, "port", "p", "49090", "port number to get a certificate and a secret key from a user")
}
