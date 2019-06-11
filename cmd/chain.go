package cmd

import (
	"context"
	pb "tethys-console/services"
	"time"

	"github.com/spf13/cobra"
)

type ResLoadChain struct {
	res *pb.ResLoadChain
}

func (resp ResLoadChain) Info() string {
	return resp.res.Info
}

func (resp ResLoadChain) Error() string {
	return resp.res.Info
}

func (resp ResLoadChain) Success() bool {
	return resp.res.Success
}

// chainCmd represents the chain command
var chainCmd = &cobra.Command{
	Use:   "chain",
	Short: "Load the chain config file",
	Long: `
	Load the chain config file
	This command have to execute preceding a 'world' command.
	`,
}

var chainPath string

func chain(cmd *cobra.Command, args []string) {
	infoLogger.Println("Load the config file from ", chainPath)

	if checkFilePathExist(chainPath) == false {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, client, err := beforeExecute(&ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	resp, err := client.LoadChain(ctx, &pb.ReqLoadChain{Path: chainPath})

	responseLoadChain := ResLoadChain{resp}
	afterExecute(responseLoadChain, err)
}

func init() {
	rootCmd.AddCommand(chainCmd)

	chainCmd.Run = chain

	chainCmd.Flags().StringVarP(&chainPath, "path", "p", "", "'load_chain.json' path")
}
