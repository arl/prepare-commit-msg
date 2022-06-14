package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	commitFilename := ""
	if len(os.Args) < 2 {
		log.Fatalf("not enough arguments")
	}
	if len(os.Args) > 1 {
		commitFilename = os.Args[1]
	}
	_ = commitFilename

	out, err := exec.Command("git", "status", "--porcelain=v1", "-z", "--untracked-files=no").CombinedOutput()
	if err != nil {
		log.Fatalf("git status error: %s", out)
	}

	r := bytes.NewReader(out)
	fnames, err := parse(r)
	if err != nil {
		log.Fatalf("error parsing git porcelain output: %v", err)
	}

	msgPrefix := ""
	switch len(fnames) {
	case 0:
		return // nothing to do
	case 1:
		msgPrefix = fnames[0] + ": "
	default:
		// multiple filenames, find the common path prefix, if any.
		common := commonPrefix(fnames)
		if len(common) == 0 {
			msgPrefix = "all: "
		} else {
			msgPrefix = common + ": "
		}
	}

	if err := appendMsg(commitFilename, msgPrefix); err != nil {
		log.Fatalf("error writing into commit message file: %v", err)
	}
}

func appendMsg(fname, msg string) error {
	f, err := os.OpenFile(fname, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	bb := bytes.Buffer{}
	bb.WriteString(msg)
	bb.Write(buf)

	if err := f.Truncate(0); err != nil {
		return err
	}

	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return err
	}
	_, err = f.Write(bb.Bytes())
	return err
}

func commonPrefix(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	if len(paths) == 1 {
		return paths[0]
	}

	common := strings.Split(paths[0], string(filepath.Separator))
	for i := 1; i < len(paths); i++ {
		dirs := strings.Split(paths[i], string(filepath.Separator))
		j := 0
		for ; j < len(common); j++ {
			if common[j] == dirs[j] {
				continue
			}
			break
		}
		common = common[:j]
		if len(common) == 0 {
			break
		}
	}
	return filepath.Join(common...)
}

func parse(r io.Reader) ([]string, error) {
	var fnames []string

	scan := bufio.NewScanner(r)
	scan.Split(scanNilBytes)
	for scan.Scan() {
		if len(scan.Text()) < 4 {
			continue
		}
		fnames = append(fnames, scan.Text()[3:])
	}

	return fnames, scan.Err()
}

// scanNilBytes is a bufio.SplitFunc function used to tokenize the input with
// nil bytes. The last byte should always be a nil byte or scanNilBytes returns
// an error.
func scanNilBytes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, 0); i >= 0 {
		// We have a full nil-terminated line.
		return i + 1, data[0:i], nil
	}

	// If we're at EOF, we would have a final not ending with a nil byte, we
	// won't allow that.
	if atEOF {
		return 0, nil, errors.New("last line doesn't end with a nil byte")
	}
	// Request more data.
	return 0, nil, nil
}
