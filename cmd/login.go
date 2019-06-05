package cmd

import (
	"context"

	pb "gruut-console/services"

	"github.com/gen2brain/beeep"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short:   "Setup the node",
	Long:    ``,
}

var password string

func login(cmd *cobra.Command, args []string) {
	connOption := grpc.WithInsecure()
	conn, err := grpc.Dial(address, connOption)
	if err != nil {
		errorLogger.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewTethysAdminServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := client.Login(ctx, &pb.ReqLogin{Password: password})

	if err != nil {
		beepError := beeep.Notify("Error", err.Error(), "assets/warning.png")
		if beepError != nil {
			panic(beepError)
		}
	}

	if resp.Success == true {
		err := beeep.Notify("[LOGIN]", "Login has successfully completed.", "assets/information.png")
		if err != nil {
			panic(err)
		}
	} else {
		infoLogger.Println(resp)
		
		err := beeep.Notify("[LOGIN]", "Failed to login.", "assets/information.png")
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Run = login

	loginCmd.Flags().StringVarP(&password, "password", "p", "", "password to decode a pem format private key")
}
