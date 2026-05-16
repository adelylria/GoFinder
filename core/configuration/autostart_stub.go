//go:build !windows

package configuration

func ApplyAutoStart(enabled bool) error {
	return nil
}
