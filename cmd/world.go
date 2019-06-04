package cmd

import (
	"os"
	"path/filepath"
	"context"
	pb "gruut-console/services"
	"time"

	"github.com/gen2brain/beeep"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// worldCmd represents the world command
var worldCmd = &cobra.Command{
	Use:   "world",
	Short: "Load the world config file",
	Long:  ``,
}

var path string

func world(cmd *cobra.Command, args []string) {
	infoLogger.Println("Load the config file from ", path)

	if checkFilePathExist() == false {
		return
	}

	connOption := grpc.WithInsecure()
	conn, err := grpc.Dial(address, connOption)
	if err != nil {
		errorLogger.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewGruutAdminServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.LoadWorld(ctx, &pb.ReqLoadWorld{Path: path})

	if err != nil {
		beepError := beeep.Notify("Error", err.Error(), "assets/warning.png")
		if beepError != nil {
			panic(beepError)
		}
	}

	if resp.Success == true {
		err := beeep.Notify("[LOAD_WORLD]", "The loading has successfully completed.", "assets/information.png")
		if err != nil {
			panic(err)
		}
	} else {
		infoLogger.Println(resp)

		err := beeep.Notify("[LOAD_WORLD]", "Failed to load.", "assets/information.png")
		if err != nil {
			panic(err)
		}
	}
}

func checkFilePathExist() bool {
	worldConfigFilePath := path

	if !filepath.IsAbs(path) {
		worldConfigFilePath, _ = filepath.Abs(worldConfigFilePath)
	}

	if _, err := os.Stat(worldConfigFilePath); os.IsNotExist(err) {
		errorLogger.Println("File does not exist.")
		return false
	}
 
	return true
}

func init() {
	rootCmd.AddCommand(worldCmd)

	worldCmd.Run = world

	worldCmd.Flags().StringVarP(&path, "path", "p", "", "'world_create.json' path")
}
