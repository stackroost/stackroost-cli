/*
Copyright © 2025 Stackroost CLI

*/
package remote

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage remote servers",
	Long:  `Add, list, and execute commands on remote servers via SSH.`,
}

var remoteAddCmd = &cobra.Command{
	Use:   "add [name] [user@host] --key [keyfile]",
	Short: "Add a remote server",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		userHost := args[1]
		key, _ := cmd.Flags().GetString("key")
		viper.Set("remotes."+name+".userhost", userHost)
		viper.Set("remotes."+name+".key", key)
		viper.WriteConfig()
		fmt.Printf("Added remote %s: %s\n", name, userHost)
	},
}

var remoteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List remote servers",
	Run: func(cmd *cobra.Command, args []string) {
		remotes := viper.GetStringMap("remotes")
		for name := range remotes {
			userhost := viper.GetString("remotes." + name + ".userhost")
			fmt.Printf("%s: %s\n", name, userhost)
		}
	},
}

var remoteExecCmd = &cobra.Command{
	Use:   "exec [name] [command]",
	Short: "Execute command on remote server",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		command := strings.Join(args[1:], " ")
		userhost := viper.GetString("remotes." + name + ".userhost")
		keyfile := viper.GetString("remotes." + name + ".key")
		executeRemoteCommand(userhost, keyfile, command)
	},
}

func AddRemoteCmd(root *cobra.Command) {
	root.AddCommand(remoteCmd)

	remoteCmd.AddCommand(remoteAddCmd)
	remoteCmd.AddCommand(remoteListCmd)
	remoteCmd.AddCommand(remoteExecCmd)

	remoteAddCmd.Flags().String("key", "", "SSH key file")
}

func executeRemoteCommand(userhost, keyfile, command string) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		fmt.Println("Unable to read private key:", err)
		return
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		fmt.Println("Unable to parse private key:", err)
		return
	}

	parts := strings.Split(userhost, "@")
	if len(parts) != 2 {
		fmt.Println("Invalid user@host format")
		return
	}
	user := parts[0]
	host := parts[1]

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		fmt.Println("Failed to dial:", err)
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Println("Failed to create session:", err)
		return
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	err = session.Run(command)
	if err != nil {
		fmt.Println("Failed to run command:", err)
	}
}