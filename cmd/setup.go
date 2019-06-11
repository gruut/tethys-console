package cmd

import (
	"context"
	pb "tethys-console/services"

	"github.com/spf13/cobra"
)

type ResSetup struct {
	res *pb.ResSetupKey
}

func (resp ResSetup) Info() string {
	return resp.res.Info
}

func (resp ResSetup) Error() string {
	return resp.res.Info
}

func (resp ResSetup) Success() bool {
	return resp.res.Success
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:     "setup",
	Short:   "Setup the node",
	Long:    ``,
}

var port string

func setup(cmd *cobra.Command, args []string) {
	infoLogger.Println("It may take a while...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, client, err := beforeExecute(&ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	resp, err := client.SetupKey(ctx, &pb.ReqSetupKey{SetupPort: port})
	if err != nil {
		errorLogger.Fatalln("A connection was not established because of following error: ", err.Error())
	}

	responseSetup := ResSetup{resp}
	afterExecute(responseSetup, err)
}

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.Run = setup

	setupCmd.Flags().StringVarP(&port, "port", "p", "49090", "port number to get a certificate and a secret key from a user")
}
