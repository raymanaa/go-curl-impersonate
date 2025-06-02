package curl

import (
	"testing"
)

func TestVersionInfo(t *testing.T) {
	info := VersionInfo(VERSION_FIRST)
	expectedProtocols := []string{
		"dict", "file", "ftp", "ftps", "gopher", "gophers",
		"http", "https", "imap", "imaps", "ldap", "ldaps",
		"mqtt",
		"pop3", "pop3s", "rtmp", "rtsp", "scp", "sftp", "smb", "smbs",
		"smtp", "smtps", "telnet", "tftp",
		"ws", "wss",
	}
	protocols := info.Protocols
	for _, protocol := range protocols {
		found := false
		for _, expectedProtocol := range expectedProtocols {
			if expectedProtocol == protocol {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("protocol should be in %v and is %v.", expectedProtocols, protocol)
		}
	}
}
