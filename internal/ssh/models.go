package ssh

import "golang.org/x/crypto/ssh"

type sshClient struct {
	host string
	user string
	client *ssh.Client
}