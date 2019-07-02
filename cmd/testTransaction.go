package cmd

import (
	"io/ioutil"
	"os"
	"fmt"
	"crypto/sha256"
	"context"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	pb "tethys-console/services/grpc_merger"
)

// testTransactionCmd represents the testTransaction command
var testTransactionCmd = &cobra.Command{
	Use:   "testTransaction",
	Short: "Command in order to send transactions to a node",
}

var p2pAddr string

func testTransaction(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connOption := grpc.WithInsecure()
	conn, err := grpc.Dial(p2pAddr, connOption)
	if err != nil {
		errorLogger.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewTethysMergerServiceClient(conn)
	if err != nil {
		errorLogger.Fatalf("Failed to initialize a client: %v", err)
	}

	message_id := generateMessageID()
	message := generateMessageBody()

	msgStatus, err := client.MergerService(ctx, &pb.RequestMsg{Broadcast: false, MessageId: message_id, Message: message})
	if err != nil {
		errorLogger.Fatalln("A connection was not established because of following error: ", err.Error())
	}
	
	fmt.Println(msgStatus)
}

func generateMessageID() string {
	h := sha256.New()
	h.Write([]byte("hello world\n"))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func generateMessageBody() []byte {
	jsonFile, err := os.Open("test_tx.json")
	if err != nil {
		errorLogger.Fatalln("Cannot read the json file")
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	
	return byteValue
}

func init() {
	rootCmd.AddCommand(testTransactionCmd)

	testTransactionCmd.Run = testTransaction

	loginCmd.Flags().StringVarP(&p2pAddr, "address", "addr", "10.10.10.200:59001", "The address of a merger (Default address: 10.10.10.200:59001)")
}
