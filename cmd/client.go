/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"encoding/json"
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
		// x := map[string]string{}
		var x map[string]interface{}


		// var p v.UserRequest;
		// var sb strings.Builder
		// sb.WriteString("")
		// sb.WriteString(args[0])
		// sb.WriteString("")

		// map1 := map[string]string{args[0]}
		// name := 
		json.Unmarshal([]byte(args[0]), &x);
		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	return
		// }
		// name := make(map[string]interface{args[0]});
		// name, err := json.Marshal([]byte(args[0]));

		// map1 := map[string]string{sb};
		// fmt.Println("User Name/(0) is: %d", sb.String())
		// fmt.Println("User Name/(name) is: %d", name)
		// fmt.Println("User Name/(name) is: %d", json.Valid(x))
		// fmt.Println("User Name/(mar) is: %d", json.Unmarshal([]byte(name)))
		fmt.Println("User Name/(map) is: %d", x)
		fmt.Println("User Name/(p) is: %s", x["Name"])

		request := &v.UserRequest{Id:1};
		// request := &v.UserRequest{Id: args[0].id, Name: args[0].name}

		resp, _ := client.GetUser(context.Background(), request)
		fmt.Printf("User's name received: %v\n", resp)

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
