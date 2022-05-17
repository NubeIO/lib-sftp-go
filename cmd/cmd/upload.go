package cmd

import (
	"fmt"
	"github.com/NubeIO/lib-sftp-go/pkg/ftp"
	"github.com/spf13/cobra"
	"time"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload a file",
	Long:  ``,
	RunE:  runUpload,
}

func runUpload(cmd *cobra.Command, args []string) error {
	fmt.Println("Sftp Upload started ", time.Now().String())

	session := &ftp.Session{
		FromPath: flags.fromPath, //..
		FromFile: flags.toFile,   //test.txt
		ToPath:   flags.toPath,   //home/user/new-dir
		ToFile:   flags.toFile,   //new.txt

		Host: &ftp.Host{
			Host:     flags.host,
			Port:     flags.port,
			User:     flags.username,
			Password: flags.password,
		},
	}
	ftpClient, err := ftp.NewConn(session)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = ftpClient.Upload()
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Println("Sftp Upload finished ", time.Now().String())
	}
	return nil
}

var flags struct {
	host     string
	username string
	password string
	port     int
	verbose  bool
	fromPath string
	fromFile string
	toPath   string
	toFile   string
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.PersistentFlags().StringVarP(&flags.host, "host", "", "localhost", "host hostname")
	uploadCmd.PersistentFlags().IntVarP(&flags.port, "port", "", 22, "host port")
	uploadCmd.PersistentFlags().StringVarP(&flags.username, "user", "", "username", "host username")
	uploadCmd.PersistentFlags().StringVarP(&flags.password, "pass", "", "password", "host password")

	uploadCmd.PersistentFlags().StringVarP(&flags.fromPath, "from-path", "", "", "path of the from host")
	uploadCmd.PersistentFlags().StringVarP(&flags.fromFile, "from-file", "", "", "file name of the from host")
	uploadCmd.PersistentFlags().StringVarP(&flags.toPath, "to-path", "", "", "path of the to host")
	uploadCmd.PersistentFlags().StringVarP(&flags.toFile, "to-file", "", "", "file name of the to host")
}
