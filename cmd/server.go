/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	v "github.com/aadityadev/mstodo/pkg/api/v1"
	"golang.org/x/xerrors"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type server struct {
	v.UnimplementedUserServiceServer
}

func (*server) GetUser(ctx context.Context, request *v.UserRequest) (*v.UserResponse, error) {
	response := &v.UserResponse{};

	if request == nil {
		fmt.Println("request must not be nil")
		return response, xerrors.Errorf("request must not be nil")
	}
	id := request.Id
	fmt.Printf("User ID is: %d", id)

	name := request.Name;
	fmt.Println("User Name is: %d", name)

	response = &v.UserResponse{
		Name: "John Doe",
	}
	return response, nil
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")

		if len(os.Args) != 2 {
			log.Fatal("Please specify to start the server or client.")
		}
	
		osArgs := os.Args[1:]
		fmt.Println(osArgs)
	
		// if osArgs[0] == "server" {
			address := "0.0.0.0:8000"
			
			lis, err := net.Listen("tcp", address)
			// lis, err := 
			http.HandleFunc("/createTodo", CreateTodo)
			http.HandleFunc("/hello", HelloServer)
			http.HandleFunc("/", HelloServer)
			http.ListenAndServe(":8080", nil)
			
			db, err := sql.Open("mysql", "user7:s$cret@tcp(127.0.0.1:3306)/testdb")
			defer db.Close()

			if err != nil {
				log.Fatal(err)
			}
			
			if err != nil {
				log.Fatalf("Error %v", err)
			}
	
			fmt.Printf("Server is listening on %v ...\n", address)
	
			s := grpc.NewServer()
			v.RegisterUserServiceServer(s, &server{})
	
			s.Serve(lis)
		// } else if osArgs[0] == "client" {
		// 	fmt.Printf("Serv");
		// 	RunClient()
		// }
	},
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	dat, err := json.Marshal(r.Body)
	fmt.Println("dat %s", dat);
	fmt.Println("err %s", err);
	w.Header().Add("Content-Type", "application/json")
    fmt.Fprintf(w, "Hello, %s!", dat)
	// w.Write(dat);
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
