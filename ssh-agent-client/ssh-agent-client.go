package sshAgentClient

import "os"
import "net"
import "fmt"
import "errors"
import "io/ioutil"
import "crypto/x509"
import "encoding/pem"
// import "golang.org/x/crypto/ssh"
import "golang.org/x/crypto/ssh/agent"
// import "github.com/lunixbochs/go-keychain"
import "github.com/keybase/go-keychain"
import "github.com/codegangsta/cli"

var MyAgent (agent.Agent)

func AbortOnError(err error, msg string) {
  if err == nil { return }
  fmt.Fprintf(os.Stderr, "ERROR: %s\n%s\n", msg, err)
  os.Exit(1)
}

func GetAgent() (agent.Agent) {
  if MyAgent == nil {
    sock_path := os.Getenv("SSH_AUTH_SOCK")
    sock_addr, err := net.ResolveUnixAddr("unix", sock_path)
    AbortOnError(err, "Can't find socket " + sock_path)
    auth_sock, err := net.DialUnix("unix", nil, sock_addr)
    AbortOnError(err, "Can't DialUnix to " + sock_path)
    MyAgent = agent.NewClient(auth_sock)
  }
  return MyAgent
}

func Add(key_path string) int {
  _, err := os.Stat(key_path)
  AbortOnError(err, "Can't read private key file" + key_path)
  key_buffer, err := ioutil.ReadFile(key_path)
  AbortOnError(err, "Reading private key " + key_path)
  key_block, _ := pem.Decode(key_buffer)
  if ! x509.IsEncryptedPEMBlock(key_block){
    AbortOnError(errors.New(""), "Key is plaintext!!!")
  }

  pw, err := keychain.GetGenericPassword("SSH", key_path, "", "")
  AbortOnError(err, "Can't get pw from keychain for " + key_path)
  key_block_plaintext, err := x509.DecryptPEMBlock(key_block, []byte(pw))
  AbortOnError(err, "decrypting key " + key_path)

  var add_key agent.AddedKey
  add_key.PrivateKey, err = x509.ParsePKCS1PrivateKey(key_block_plaintext)
  AbortOnError(err, "parsing decrypted private key for " + key_path)
  add_key.Comment = key_path
  add_key.ConfirmBeforeUse = true

  err = GetAgent().Add(add_key)
  AbortOnError(err, "adding decrypted key " + key_path + " to agent")

  return 0
}

func CliAdd(c *cli.Context){
  Add(c.Args().First())
  return
}
