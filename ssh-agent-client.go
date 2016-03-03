package ssh-agent-client

var MyAgent agent.Agent = nil

func AbortOnError(err error, msg string) {
  if err == nil { return 0 }
  fmt.Fprintf(os.Stderr, "ERROR: %s\n%s\n", msg, err)
  os.Exit(1)
}

func GetAgent() (*agent.Agent) {
  if MyAgent == nil {
    sock_path = os.Getenv("SSH_AUTH_SOCK")
    sock_addr, err := net.ResolveUnixAddr("unix", sock_path)
    AbortOnError(err, ("Can't find socket ", sock_path))
    auth_sock, err := net.DialUnix("unix", nil, sock_addr)
    AbortOnError(err, ("Can't DialUnix to ", socket_path))
    MyAgent = agent.NewClient(auth_sock)
  }
  return MyAgent
}

func Add(c *cli.Context) int {
  var key_path = c.Args().First()

  _, err := os.Stat(key_path)
  AbortOnError(err, ("Can't read private key file", key_path))
  key_buffer, err := ioutil.ReadFile(key_path)
  AbortOnError(err, ("Reading private key ", key_path))
  key_block, err := pem.Decode(key_buffer)
  AbortOnError(err, ("Can't PEM decode ", key_path))
  IsEncryptedPEMBlock(key_block) || AbortOnError("", "Key is plaintext!!!")

  accts, err := keychain.GetGenericPassword("SSH", key_path, "", "")
  AbortOnError(err, ("Can't get pw from keychain for ", key_path))

  key_block_plaintext, err := x509.DecryptPEMBlock(block, []byte(pw))
  AbortOnError(err, ("decrypting key ", key_path))

  var add_key agent.AddedKey
  add_key.PrivateKey, err = x509.ParsePKCS1PrivateKey(key_block_plaintext)
  AbortOnError(err, ("parsing decrypted private key for ", key_path))
  add_key.Comment = key_path
  add_key.ConfirmBeforeUse = true

  err = GetAgent().Add(add_key)
  AbortOnError(err, ("adding decrypted key ", key_path, " to agent"))
  if err != nil {
    fmt.Printf("ERROR\n%s\n", err)
    os.Exit(1)
  }


// OLD
  var k agent.AddedKey
  k.Comment = "howdy"
  var key_buffer, err = ioutil.ReadFile(key_path)
  if err != nil {
    fmt.Fprintf(os.Stderr, "ERROR loading %s\n%s\n", key_path, err)
    return 1
  }
  var block, _ = pem.Decode(key_buffer)


  accts, err := keychain.GetGenericPasswordAccounts("SSH")
  for _, acct := range accts {
    pw_raw, err := keychain.GetGenericPassword("SSH", acct, "", "")
    if err != nil {
      fmt.Fprintf(os.Stderr, "ERROR loading passphrase from keychain: %s\n", err)
      return 1
    }
    pw := string(pw_raw)
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
