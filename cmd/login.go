package cmd

import (
	"context"

	pb "tethys-console/services/grpc_admin"

	"github.com/spf13/cobra"
)

type ResLogin struct {
	res *pb.ResLogin
}

func (resp ResLogin) Info() string {
	return resp.res.Info
}

func (resp ResLogin) Error() string {
	return resp.res.Info
}

func (resp ResLogin) Success() bool {
	return resp.res.Success
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Setup the node",
	Long:  ``,
}

var password string

func login(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, client, err := beforeExecute(&ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	resp, err := client.Login(ctx, &pb.ReqLogin{Password: password})
	if err != nil {
		errorLogger.Fatalln("A connection was not established because of following error: ", err.Error())
	}
	
	responseLogin := ResLogin{resp}
	afterExecute(responseLogin, err)
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Run = login

	loginCmd.Flags().StringVarP(&password, "password", "p", "", "password to decode a pem format private key")
}
