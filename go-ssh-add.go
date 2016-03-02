package main

import "os"
import "net"
import "fmt"
import "io/ioutil"
import "crypto/x509"
import "encoding/pem"
import "golang.org/x/crypto/ssh"
import "golang.org/x/crypto/ssh/agent"
// import "github.com/lunixbochs/go-keychain"
import "github.com/keybase/go-keychain"


func main() {
  fmt.Printf("hello, %s\n", ssh.CertAlgoRSAv01)

  var key_path = "/Users/tmaher/.ssh/id_rsa"
  var k agent.AddedKey
  k.Comment = "howdy"
  var key_buffer, err = ioutil.ReadFile(key_path)
  if err != nil {
    fmt.Printf("err is non-nill: %s\n", err)
    os.Exit(1)
  }
  var block, _ = pem.Decode(key_buffer)

  // lunixbochs code
  // var pw, _ = keychain.Find("SSH", key_path)
  // fmt.Printf("pw: %s\n", pw)

  // keybase code
  var pw = ""
  accts, err := keychain.GetGenericPasswordAccounts("SSH")
  for _, acct := range accts {
    fmt.Printf("found account %s\n", acct)
    pw_raw, err := keychain.GetGenericPassword("SSH", acct, "", "")
    if err != nil {
      fmt.Printf("Noooo can't find pw for %s\n", acct)
    } else {
      pw = string(pw_raw)
      fmt.Printf("%s: %s\n", acct, pw)
    }
  }
  os.Exit(1)

  var block_plaintext, _ = x509.DecryptPEMBlock(block, []byte(pw))
  var add_key agent.AddedKey
  add_key.PrivateKey, err = x509.ParsePKCS1PrivateKey(block_plaintext)
  add_key.Comment = key_path
  add_key.ConfirmBeforeUse = true

  unix_addr, err := net.ResolveUnixAddr("unix", os.Getenv("SSH_AUTH_SOCK"))
  auth_sock, err := net.DialUnix("unix", nil, unix_addr)
  var a agent.Agent = agent.NewClient(auth_sock)
  var keys []*agent.Key
  keys, err = a.List()
  for index, element := range keys {
    fmt.Printf("index: %d\n", index)
    fmt.Printf("key: %s\n", element.Comment)
  }
  err = a.Add(add_key)
  if err != nil {
    fmt.Printf("ERROR\n%s\n", err)
    os.Exit(1)
  }

}
