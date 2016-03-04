package main

import "os"
import "github.com/codegangsta/cli"
import "github.com/tmaher/go-ssh-add/ssh-agent-client"

func main() {
  app := cli.NewApp()
  app.Name = "go-ssh-add"
  app.Usage = "Load SSH private keys to ssh-agent, requring confirmation on use"

  app.Commands = []cli.Command{
    {
      Name:     "list",
      Aliases:  []string{"l"},
      Usage:    "List loaded keys",
      Action:   sshAgentClient.CliList,
    },
    {
      Name:      "add",
      Aliases:   []string{"a"},
      Usage:     "add a key",
      Action:    sshAgentClient.CliAdd,
    },
    {
      Name:      "delete",
      Aliases:   []string{"d"},
      Usage:     "Delete specified key",
      Action:    sshAgentClient.CliDelete,
    },
    {
      Name:      "delete-all",
      Usage:     "Delete *ALL* loaded keys from agent",
      Action:    sshAgentClient.CliDeleteAll,
    },
    {
      Name:       "add-all",
      Aliases:    []string{"add-all"},
      Usage:      "auto-load all keys listed in keychain",
      Action:     sshAgentClient.CliAddAll,
    },
  }

  app.Run(os.Args)


}
