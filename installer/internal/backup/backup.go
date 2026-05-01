package backup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func SnapshotIfExists(dir string, files []string) error {
	any := false
	for _, f := range files {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			any = true
			break
		}
	}
	if !any {
		return nil
	}
	stamp := time.Now().Format("20060102-150405")
	dst := filepath.Join(dir, "backups", fmt.Sprintf("designapi-%s", stamp))
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	for _, f := range files {
		src := filepath.Join(dir, f)
		if _, err := os.Stat(src); err != nil {
			continue
		}
		if err := copyFile(src, filepath.Join(dst, f)); err != nil {
			return err
		}
	}
	fmt.Printf("  ↳ backup: %s\n", dst)
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
