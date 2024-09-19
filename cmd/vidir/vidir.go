package vidir

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/gabe565/moreutils/cmd/cmdutil"
	"github.com/gabe565/moreutils/internal/editor"
	"github.com/spf13/cobra"
)

const (
	Name          = "vidir"
	FlagVerbose   = "verbose"
	FlagRecursive = "recursive"
)

func New(opts ...cmdutil.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     Name + " [file | dir]...",
		Short:   "Edit a directory in your text editor",
		RunE:    run,
		GroupID: cmdutil.Applet,
	}

	cmd.Flags().BoolP(FlagVerbose, "v", false, "Verbosely display the actions taken by the program.")
	cmd.Flags().BoolP(FlagRecursive, "r", false, "Recurses into subdirectories")

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

var (
	ErrInvalidIndex = errors.New("invalid index")
	ErrFieldCount   = errors.New("invalid field count")
)

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	tmp, err := os.CreateTemp("", "vidir-*.txt")
	if err != nil {
		return err
	}
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
	}()

	recursive, err := cmd.Flags().GetBool(FlagRecursive)
	if err != nil {
		panic(err)
	}

	paths, err := createListing(tmp, args, recursive)
	if err != nil {
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	if err := editor.Edit(tmp.Name(), true); err != nil {
		return err
	}

	tmp, err = os.Open(tmp.Name())
	if err != nil {
		return err
	}

	verbose, err := cmd.Flags().GetBool(FlagVerbose)
	if err != nil {
		panic(err)
	}

	seen := make([]int, 0, len(paths))

	scanner := bufio.NewScanner(tmp)
	for scanner.Scan() {
		parts := bytes.SplitN(bytes.TrimSpace(scanner.Bytes()), []byte("\t"), 2)
		switch len(parts) {
		case 0:
			continue
		case 1:
			return fmt.Errorf("%w: %d", ErrFieldCount, len(parts))
		}

		i, err := strconv.Atoi(string(parts[0]))
		if err != nil {
			return err
		}
		i--
		if i < 0 || i > len(paths)-1 {
			return fmt.Errorf("%w: %d", ErrInvalidIndex, i+1)
		}

		oldName := paths[i]
		newName := string(parts[1])

		seen = append(seen, i)
		if oldName == newName {
			continue
		}

		tmpName := newName
		for i := 0; ; i++ {
			if _, err := os.Stat(tmpName); err == nil {
				// New file already exists
				if i == 0 {
					tmpName += "~"
				} else {
					tmpName = newName + "~" + strconv.Itoa(i)
				}
			} else if !os.IsNotExist(err) {
				return err
			} else {
				break
			}
		}
		if tmpName != newName {
			if err := rename(cmd.OutOrStdout(), newName, tmpName, verbose); err != nil {
				return err
			}
			for k, v := range paths {
				if newName == v {
					paths[k] = tmpName
				}
			}
		}

		if err := rename(cmd.OutOrStdout(), oldName, newName, verbose); err != nil {
			return err
		}
	}

	for i, name := range paths {
		if !slices.Contains(seen, i) {
			if err := os.Remove(name); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

func createListing(w io.Writer, args []string, recursive bool) ([]string, error) {
	if len(args) == 0 {
		args = append(args, ".")
	}

	paths := make([]string, 0, len(args))

	buf := bufio.NewWriter(w)
	var i int
	for _, arg := range args {
		glob, err := filepath.Glob(arg)
		if err != nil {
			return nil, err
		}

		for _, globPath := range glob {
			if err := filepath.WalkDir(globPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() && path == globPath {
					// Only show files when a directory is passed as a param
					return nil
				}

				i++
				if _, err := buf.WriteString(fmt.Sprintf("%04d\t%s\n", i, path)); err != nil {
					return err
				}

				paths = append(paths, path)

				if !recursive && d.IsDir() && !slices.Contains(args, path) {
					// Do not recurse more than 1 level
					return filepath.SkipDir
				}
				return nil
			}); err != nil {
				return nil, err
			}
		}
	}

	return paths, buf.Flush()
}

func rename(w io.Writer, oldname, newname string, verbose bool) error {
	if verbose {
		_, _ = fmt.Fprintf(w, "%q => %q\n", oldname, newname)
	}
	return os.Rename(oldname, newname)
}
