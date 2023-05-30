// Package netrc provides methods to handle basic authentification used by the auto-login process in .netrc f.
// See https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html.
package netrc

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jdxcode/netrc"
	"github.com/rvflash/goup/internal/vcs"
)

// File represents a netrc f.
type File struct {
	f *netrc.Netrc
}

// Parse the netrc f behind the NETRC environment variable.
// If undefined, we try to find the .netrc f in the user home directory.
func Parse() (rc File, err error) {
	path := os.Getenv("NETRC")
	if path == "" {
		path, err = os.UserHomeDir()
		if err != nil {
			return rc, fmt.Errorf("user home directory: %w", err)
		}
		path = filepath.Join(path, ".netrc")
	}
	f, err := netrc.Parse(path)
	if err != nil {
		if os.IsNotExist(err) {
			return rc, nil
		}
		return rc, fmt.Errorf("netrc parsing %q: %w", path, err)
	}
	return File{f: f}, nil
}

// BasicAuth returns if available the basic auth information for this hostname.
func (c File) BasicAuth(host string) *vcs.BasicAuth {
	if c.f == nil {
		return nil
	}
	h := c.f.Machine(host)
	if h == nil {
		return nil
	}
	return &vcs.BasicAuth{
		Username: h.Get("login"),
		Password: h.Get("password"),
	}
}
