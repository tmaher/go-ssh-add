package sshAgentClient

import "os"
import "net"
import "fmt"
import "io/ioutil"
import "encoding/base64"
import "crypto/md5"
import "crypto/sha1"
import "crypto/x509"
import "encoding/pem"
//import "golang.org/x/crypto/ssh"
import "golang.org/x/crypto/ssh"
import "golang.org/x/crypto/ssh/agent"
import "github.com/keybase/go-keychain"
import "github.com/codegangsta/cli"

var MyAgent (agent.Agent)

func Abort(msg string){
  fmt.Fprintf(os.Stderr, "ERROR: %s\n", msg)
  os.Exit(1)
}

func AbortOnError(err error, msg string) {
  if err == nil { return }
  Abort(msg + "\n" + err.Error() + "\n")
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
    Abort("Key is plaintext!!!")
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

func CliAddAll(c *cli.Context){
  accts, err := keychain.GetGenericPasswordAccounts("SSH")
  AbortOnError(err, "Can't get list of keychain'd SSH keys")
  for _, key_file := range accts { Add(key_file) }
  return
}

func CliAdd(c *cli.Context){
  Add(c.Args().First())
  return
}

func CliDelete(c *cli.Context){
  Delete(c.Args().First())
}

func CliList(c *cli.Context) {
  List()
}

func CliDeleteAll(c *cli.Context) {
  DeleteAll()
}

func Delete(key_path string) int {
  keys, err := GetAgent().List()
  AbortOnError(err, "Can't get list of loaded keys from agent")
  for _, key := range keys {
    if key_path != key.Comment { continue }
    pubkey, err := ssh.ParsePublicKey(key.Blob)
    AbortOnError(err, "Can't parse public key " + key.Comment + " for deletion")
    err = GetAgent().Remove(pubkey)
    AbortOnError(err, "Can't remove public key " + key_path)
    return 0
  }
  Abort("agent does not contain key " + key_path)
  return 1
}

func DeleteAll() int {
  keys, err := GetAgent().List()
  AbortOnError(err, "Can't get list of loaded keys from agent")
  for _, key := range keys {
    pubkey, err := ssh.ParsePublicKey(key.Blob)
    AbortOnError(err, "Can't paarse public key " + key.Comment + " for deletion")
    err = GetAgent().Remove(pubkey)
    AbortOnError(err, "Can't reemove public key " + key.Comment)
  }
  return 0
}

func List() int {
  keys, err := GetAgent().List()
  AbortOnError(err, "Can't get list of loaded keys from agent")
  for _, key := range keys {
    fmt.Printf("%s %s\n", key.Format, key.Comment)
    fmt.Printf("MD5:  %x\n", md5.Sum(key.Blob))
    sha1_fp := sha1.Sum(key.Blob)
    fmt.Printf("SHA1: %s\n", base64.StdEncoding.EncodeToString(sha1_fp[:]))
    //pubkey, err := ssh.ParsePublicKey(key.Blob)
    //AbortOnError(err, "Can't parse public key %s", key.Comment)
  }
  return 0
}
