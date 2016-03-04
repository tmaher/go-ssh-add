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
      Action:   func(c *cli.Context) {
        println("unimplemented: Listing keys")
      },
    },
    {
      Name:      "add",
      Aliases:   []string{"a"},
      Usage:     "add a key",
      Action:    sshAgentClient.CliAdd,
    },
    {
      Name:      "delete",
      Aliases:    []string{"d"},
      Usage:     "complete a task on the list",
      Action: func(c *cli.Context) {
        println("delete key: ", c.Args().First())
      },
    },
    {
      Name:       "auto",
      Usage:      "auto-load all keys listed in keychain",
      Action:     func(c *cli.Context) {
        println("auto-load")
      },
    },
  }

  app.Run(os.Args)


}
