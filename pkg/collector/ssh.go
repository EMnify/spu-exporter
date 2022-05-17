package collector

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-kit/kit/log/level"
	"golang.org/x/crypto/ssh"
)

// copy from ssh exporter by Nordstrom (https://github.com/Nordstrom/ssh_exporter)
var allowedHostKeyTypes = []string{
	"ssh-rsa-cert-v01@openssh.com",
	"ssh-dss-cert-v01@openssh.com",
	"ecdsa-sha2-nistp256-cert-v01@openssh.com",
	"ecdsa-sha2-nistp384-cert-v01@openssh.com",
	"ecdsa-sha2-nistp521-cert-v01@openssh.com",
	"ssh-ed25519-cert-v01@openssh.com",
	"ecdsa-sha2-nistp256 ecdsa-sha2-nistp384",
	"ecdsa-sha2-nistp521",
	"ssh-rsa",
	"ssh-dss",
	"ssh-ed25519",
	"rsa-sha2-256",
	"rsa-sha2-512",
}

//
// SoftCheck logs non-nil errors to stderr. Used for runtime errors that should
// not kill the server.
//
func (d *SpuMetricsDaemon) SoftCheck(e error) bool {

	if e != nil {
		_ = level.Warn(d.logger).Log("message", "Error in ssh connection to spu application", "error", e)
		return true
	}
	return false
}

//
// executeScriptOnHost executes a given script on a given host.
//
func (d *SpuMetricsDaemon) executeScriptOnHost(host, port, user, keyfile, script string) (string, int, error) {

	client, session, err := d.sshConnectToHost(host, port, user, keyfile)
	if d.SoftCheck(err) {
		return "", -1, err
	}
	defer client.Close()
	defer session.Close()

	out, err := session.CombinedOutput(script)
	if d.SoftCheck(err) {
		var errorStatusCode int
		fmt.Sscanf(fmt.Sprintf("%v", err), "Process exited with status %d", &errorStatusCode)
		if errorStatusCode != 0 {
			return "", errorStatusCode, err
		}
		return "", -1, err
	}

	return literalFormat(string(out)), 0, nil

}

//
// sshConnectToHost connects to a given host with the given keyfile.
//
func (d *SpuMetricsDaemon) sshConnectToHost(host, port, user, keyfile string) (*ssh.Client, *ssh.Session, error) {

	key, err := d.getKeyFile(keyfile)
	d.SoftCheck(err)

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
		HostKeyAlgorithms: allowedHostKeyTypes,
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
func (d *SpuMetricsDaemon) getKeyFile(keyfile string) (ssh.Signer, error) {

	buf, err := ioutil.ReadFile(keyfile)
	if d.SoftCheck(err) {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buf)
	if d.SoftCheck(err) {
		return nil, err
	}

	return key, nil
}

//
// literalFormat formats a string to be included in an endpoint to be scraped by Prometheus.
//
// Turns newline characters into '\n' characters.
//
func literalFormat(input string) string {
	s1 := strings.Replace(input, "\r\n", "\\n", -1)
	return strings.Replace(s1, "\n", "\\n", -1)
}
