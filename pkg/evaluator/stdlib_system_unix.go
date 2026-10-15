//go:build !windows

package evaluator

import "os"

func isRootOrAdmin() bool {
	return os.Geteuid() == 0
}
