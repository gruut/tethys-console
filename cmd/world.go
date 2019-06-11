package cmd

import (
	"context"
	pb "tethys-console/services"
	"time"

	"github.com/spf13/cobra"
)

type ResLoadWorld struct {
	res *pb.ResLoadWorld
}

func (resp ResLoadWorld) Info() string {
	return resp.res.Info
}

func (resp ResLoadWorld) Error() string {
	return resp.res.Info
}

func (resp ResLoadWorld) Success() bool {
	return resp.res.Success
}

// worldCmd represents the world command
var worldCmd = &cobra.Command{
	Use:   "world",
	Short: "Load the world config file",
	Long:  ``,
}

var worldPath string

func world(cmd *cobra.Command, args []string) {
	infoLogger.Println("Load the config file from ", worldPath)

	if checkFilePathExist(worldPath) == false {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, client, err := beforeExecute(&ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	resp, err := client.LoadWorld(ctx, &pb.ReqLoadWorld{Path: worldPath})

	responseLoadWorld := ResLoadWorld{resp}
	afterExecute(responseLoadWorld, err)	
}

func init() {
	rootCmd.AddCommand(worldCmd)

	worldCmd.Run = world

	worldCmd.Flags().StringVarP(&worldPath, "path", "p", "", "'world_create.json' path")
}
