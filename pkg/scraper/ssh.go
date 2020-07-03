package scraper

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"strings"
)

// copy from ssh exporter by Nordstrom (https://github.com/Nordstrom/ssh_exporter)

//
// LogMsg logs a string to stdout with timestamp.
//
func LogMsg(s string) {

	log.Printf("spu_exporter :: %s", fmt.Sprintf("%s", s))
}

//
// SoftCheck logs non-nil errors to stderr. Used for runtime errors that should
// not kill the server.
//
func SoftCheck(e error) bool {

	if e != nil {
		LogMsg(fmt.Sprintf("%v", e))
		return true
	} else {
		return false
	}
}

//
// executeScriptOnHost executes a given script on a given host.
//
func executeScriptOnHost(host, port, user, keyfile, script string) (string, int, error) {

	client, session, err := sshConnectToHost(host, port, user, keyfile)
	if SoftCheck(err) {
		return "", -1, err
	}

	out, err := session.CombinedOutput(script)
	if SoftCheck(err) {
		var errorStatusCode int
		fmt.Sscanf(fmt.Sprintf("%v", err), "Process exited with status %d", &errorStatusCode)
		if errorStatusCode != 0 {
			return "", errorStatusCode, err
		} else {
			return "", -1, err
		}
	}
	defer client.Close()

	return literalFormat(string(out)), 0, nil

}

//
// sshConnectToHost connects to a given host with the given keyfile.
//
func sshConnectToHost(host, port, user, keyfile string) (*ssh.Client, *ssh.Session, error) {

	key, err := getKeyFile(keyfile)
	SoftCheck(err)

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConfig.SetDefaults()

	fullHost := fmt.Sprintf("%s:%s", host, port)
	client, err := ssh.Dial("tcp", fullHost, sshConfig)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, nil, err
	}

	return client, session, nil
}

//
// getKeyFile provides an ssh.Signer for the given keyfile (path to a private key).
//
func getKeyFile(keyfile string) (ssh.Signer, error) {

	buf, err := ioutil.ReadFile(keyfile)
	SoftCheck(err)

	key, err := ssh.ParsePrivateKey(buf)
	SoftCheck(err)

	return key, nil
}

//
// literalFormat formats a string to be included in an endpoint to be scraped by Prometheus.
//
// Turns newline characters into '\n' characters.
//
func literalFormat(input string) string {

	return strings.Replace(input, "\n", "\\n", -1)
}
