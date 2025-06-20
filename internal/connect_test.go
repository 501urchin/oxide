package internal

import (
	"fmt"
	"os"
	"strings"
	"testing"

	taurinetesting "github.com/scott-mescudi/taurine/pkg/testing"
)

func TestMain(m *testing.M) {
	lis, err := taurinetesting.StartMockSSHServer("127.0.0.1:3098")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lis.Close()

	m.Run()
}

func TestPasswordAuth(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		user        string
		host        string
		expectError bool
	}{
		{name: "valid password login", password: taurinetesting.TestPassword, user: taurinetesting.TestUser, host: "127.0.0.1:3098", expectError: false},
	}


	c := &TaurineContext{}

	tmpDir := t.TempDir()

	tmpFilePath := tmpDir + "/known_hosts"
	err := os.WriteFile(tmpFilePath, []byte(`127.0.0.1:3098 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBWiHlrQ6HS7vytwfKb32R70waRKqJ9cZOWx8RDfm4HX`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := c.ConnectWithPassword(taurinetesting.TestUser, "127.0.0.1:3098", taurinetesting.TestPassword, tmpFilePath)
			if tt.expectError && err == nil {
				client.Close()
				t.Fatal("failed to throw err")
			}

			if !tt.expectError && err != nil {
				client.Close()
				t.Fatal(err)
			}

			client.Close()
		})
	}

}
func TestKeyAuth(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		user        string
		host        string
		expectError bool
	}{
		{name: "valid key login", key: taurinetesting.PrivateKeyPEM, user: taurinetesting.TestUser, host: "127.0.0.1:3098", expectError: false},
	}

	tmpDir := t.TempDir()

	c := &TaurineContext{}


	tmpFilePath := tmpDir + "/known_hosts"
	err := os.WriteFile(tmpFilePath, []byte(`127.0.0.1:3098 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBWiHlrQ6HS7vytwfKb32R70waRKqJ9cZOWx8RDfm4HX`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := c.ConnectWithPrivateKey(taurinetesting.TestUser, "127.0.0.1:3098", []byte(taurinetesting.PrivateKeyPEM), tmpFilePath)
			if tt.expectError && err == nil {
				client.Close()
				t.Fatal("failed to throw err")
			}

			if !tt.expectError && err != nil {
				t.Fatal(err)
			}

			if client != nil {
				client.Close()
			}
		})
	}
}

func TestPasswordAutUnknownHost(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		user        string
		host        string
		expectError bool
	}{
		{name: "valid password login", password: taurinetesting.TestPassword, user: taurinetesting.TestUser, host: "127.0.0.1:3098", expectError: false},
	}
	c := &TaurineContext{}


	tmpDir := t.TempDir()

	tmpFilePath := tmpDir + "/known_hosts"
	err := os.WriteFile(tmpFilePath, []byte(`127.0.0.5:3098 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBWiHlrQ6HS7vytwfKb32R70waRKqJ9cZOWx8RDfm4HX`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	clean, err := taurinetesting.MockStdin(tmpDir+"rozz", "yes")
	if err != nil {
		t.Fatal(err)
	}
	defer clean()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := c.ConnectWithPassword(taurinetesting.TestUser, "127.0.0.1:3098", taurinetesting.TestPassword, tmpFilePath)
			if tt.expectError && err == nil {
				client.Close()
				t.Fatal("failed to throw err")
			}

			if !tt.expectError && err != nil {
				client.Close()
				t.Fatal(err)
			}

			f, err := os.ReadFile(tmpFilePath)
			if err != nil {
				t.Fatal(err)
			}

			if !strings.Contains(string(f), "127.0.0.1") {
				t.Fatalf("failed to write to known host")
			}

			client.Close()
		})
	}

}
func TestKeyAuthUnknownHost(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		user        string
		host        string
		expectError bool
	}{
		{name: "valid key login", key: taurinetesting.PrivateKeyPEM, user: taurinetesting.TestUser, host: "127.0.0.1:3098", expectError: false},
	}

	tmpDir := t.TempDir()
	tmpFilePath := tmpDir + "/known_hosts"
	err := os.WriteFile(tmpFilePath, []byte(`127.0.0.6:3098 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBWiHlrQ6HS7vytwfKb32R70waRKqJ9cZOWx8RDfm4HX`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	c := &TaurineContext{}

	clean, err := taurinetesting.MockStdin(tmpDir+"rozz", "yes")
	if err != nil {
		t.Fatal(err)
	}
	defer clean()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := c.ConnectWithPrivateKey(taurinetesting.TestUser, "127.0.0.1:3098", []byte(taurinetesting.PrivateKeyPEM), tmpFilePath)
			if tt.expectError && err == nil {
				client.Close()
				t.Fatal("failed to throw err")
			}

			if !tt.expectError && err != nil {
				t.Fatal(err)
			}

			f, err := os.ReadFile(tmpFilePath)
			if err != nil {
				t.Fatal(err)
			}

			if !strings.Contains(string(f), "127.0.0.1") {
				t.Fatalf("failed to write to known host")
			}

			if client != nil {
				client.Close()
			}
		})
	}
}
