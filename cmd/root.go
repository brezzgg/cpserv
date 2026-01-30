package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/brezzgg/cpserv/clipboard"
	"github.com/brezzgg/cpserv/log"
	"github.com/brezzgg/cpserv/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	rootCmd = &cobra.Command{Use: "cpserv"}

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {

			clip, err := clipboard.GetClipboard()
			if err != nil {
				fmt.Printf("failed to get clipboard: %s\n", err.Error())
				os.Exit(0)
			}

			lis, err := net.Listen("tcp", host)
			if err != nil {
				fmt.Printf("failed to listen: %s\n", err.Error())
				os.Exit(0)
			}
			server := grpc.NewServer()
			svc := service.New(clip)
			service.RegisterClipboardServiceServer(server, svc)

			fmt.Printf("failed to listen: %s\n", server.Serve(lis))
		},
	}

	readCmd = &cobra.Command{
		Use:   "read",
		Short: "Read remote clipboard",
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.NewClient(remote, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Error(err.Error())
			}
			client := service.NewClipboardServiceClient(conn)
			r, err := client.Read(context.TODO(), nil)
			if err != nil {
				log.Error(fmt.Sprintf("failed to read clipboard from remote: %s\n", err.Error()))
			}
			log.Response(r.Text)
		},
	}

	writeCmd = &cobra.Command{
		Use:   "write",
		Short: "Write to remote clipboard",
		Run: func(cmd *cobra.Command, args []string) {
			var text string

			dashIndex := cmd.ArgsLenAtDash()
			if dashIndex >= 0 {
				text = strings.Join(args, " ")
			} else if len(args) > 0 {
				text = strings.Join(args, " ")
			} else {
				log.Error("no text provided")
				return
			}

			conn, err := grpc.NewClient(remote, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Error(err.Error())
				return
			}
			defer conn.Close()

			client := service.NewClipboardServiceClient(conn)
			_, err = client.Write(context.TODO(), &service.WriteReq{
				Auth: nil,
				Clipboard: &service.Clipboard{
					Text: text,
				},
			})
			if err != nil {
				log.Error(fmt.Sprintf("failed to write clipboard to remote: %s\n", err.Error()))
			}
		},
	}
)

const (
	defaultHost = "0.0.0.0:56384"
)

var (
	remote string
	host   string
)

func Execute() {
	rootCmd.PersistentFlags().StringVarP(
		&remote, "remote", "r", defaultHost, "set remote host",
	)
	serverCmd.PersistentFlags().StringVarP(
		&host, "host", "", defaultHost, "set host to listen",
	)
	rootCmd.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd: true,
	}

	rootCmd.AddCommand(
		serverCmd,
		writeCmd,
		readCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Error(err.Error())
	}
}
