/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"log"

	v "github.com/aadityadev/mstodo/pkg/api/v1"
	"google.golang.org/grpc"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("client called")

		fmt.Println("Initialising gRPC client ...")

		opts := grpc.WithInsecure()
		cc, err := grpc.Dial("localhost:8000", opts)
		if err != nil {
			log.Fatal(err)
		}
		defer cc.Close()

		client := v.NewUserServiceClient(cc)
		name := args[0];
		fmt.Println("User Name/(0) is: %d", name)

		request := &v.UserRequest{Id: 1}

		resp, _ := client.GetUser(context.Background(), request)
		fmt.Printf("User's name received: %v\n", resp.Name)

	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
