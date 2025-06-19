package ssh

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/scott-mescudi/taurine/pkg/input"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func writeHostKeyToKnownHosts(knownHostsPath, host string, key ssh.PublicKey) error {
	file, err := os.OpenFile(knownHostsPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("\n%s %s", host, ssh.MarshalAuthorizedKey(key)))
	return err
}

func getHostKeyCallback(capturedKey *ssh.PublicKey) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		*capturedKey = key
		return nil
	}
}

func establishInitialConnection(host, user string, authMethod ssh.AuthMethod, knownHostsPath string) (*ssh.Client, error) {
	var hostkey ssh.PublicKey

	conf := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: getHostKeyCallback(&hostkey),
	}

	client, err := ssh.Dial("tcp", host, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to establish initial SSH connection: %v", err)
	}
	defer client.Close()

	ok := input.GetConfirmationFromUser("The authenticity of host '%s' can't be established. add (%s) to know hosts file (yes/no)?\n", host, ssh.MarshalAuthorizedKey(hostkey))
	if !ok {
		return nil, fmt.Errorf("failed to connect to ssh server: couldn't esablish authenticity")
	}

	if err := writeHostKeyToKnownHosts(knownHostsPath, host, hostkey); err != nil {
		return nil, fmt.Errorf("failed to save host key: %v", err)
	}
	return nil, nil
}

func connectSSH(user, host, knownHostsPath string, authMethod ssh.AuthMethod) (*ssh.Client, error) {
	hostKeyCallback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load known_hosts: %v", err)
	}

	conf := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: hostKeyCallback,
	}

	client, err := ssh.Dial("tcp", host, conf)
	if err != nil && strings.Contains(err.Error(), "knownhosts: key is unknown") {
		_, err := establishInitialConnection(host, user, authMethod, knownHostsPath)
		if err != nil {
			return nil, err
		}

		hostKeyCallback, err = knownhosts.New(knownHostsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to reload known_hosts: %v", err)
		}

		conf.HostKeyCallback = hostKeyCallback
		return ssh.Dial("tcp", host, conf)
	}

	return client, err
}

func ConnectWithPassword(user, host, password, knownHostsPath string) (*ssh.Client, error) {
	return connectSSH(user, host, knownHostsPath, ssh.Password(password))
}

func ConnectWithPrivateKey(user, host string, key []byte, knownHostsPath string) (*ssh.Client, error) {
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	return connectSSH(user, host, knownHostsPath, ssh.PublicKeys(signer))
}
