package cmd

import (
	"github.com/gen2brain/beeep"
	"google.golang.org/grpc"
	"github.com/spf13/viper"
	pb "gruut-console/services"
	"context"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a LDS",
	Long: `Start to running LDS. Login should be preceded to produce, query the block.`,
	Run: start,
}

func start(cmd *cobra.Command, args []string)  {
	monitorMode := viper.GetBool("monitor")

	var mode pb.ReqStart_Mode

	if monitorMode == true {
		infoLogger.Println("monitor mode")

		mode = pb.ReqStart_MONITOR
	} else {
		infoLogger.Println("default mode")
		mode = pb.ReqStart_DEFAULT
	}

	connOption := grpc.WithInsecure()
	conn, err := grpc.Dial(address, connOption)
	if err != nil {
		errorLogger.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewGruutAdminServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := client.Start(ctx, &pb.ReqStart{Mode: mode})

	if err != nil {
		beepError := beeep.Notify("Error", err.Error(), "assets/warning.png")
		if beepError != nil {
			panic(beepError)
		}
	}

	if resp.Success == true {
		err := beeep.Notify("[START KEY]", resp.Info, "assets/information.png")
		if err != nil {
			panic(err)
		}
	} else {
		infoLogger.Println(resp)
		
		err := beeep.Notify("[START KEY]", "Failed to start.", "assets/information.png")
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolP("monitor", "m", false, "")
	viper.BindPFlag("monitor", startCmd.Flags().Lookup("monitor"))
}
