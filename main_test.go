package main

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainHelp(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"eyez", "--help"}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "USAGE:") {
		t.Errorf("expected terminal output to contain USAGE block, got\n%s", output)
	}
}

func TestMainVersion(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"eyez", "--version"}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Version:") {
		t.Errorf("expected terminal output to contain Version: block, got\n%s", output)
	}
}

func TestMainMissingArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"eyez"}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Error: missing required arguments") {
		t.Errorf("expected terminal output to instruct missing arguments, got\n%s", output)
	}
}

func createTempPNG(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test_image.png")

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	err = png.Encode(f, img)
	if err != nil {
		t.Fatalf("failed to encode png: %v", err)
	}
	return path
}

func createTempSVG(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test_image.svg")
	svgContent := `<svg viewBox="0 0 50 50" xmlns="http://www.w3.org/2000/svg"><rect width="50" height="50" fill="red"/></svg>`
	if err := os.WriteFile(path, []byte(svgContent), 0644); err != nil {
		t.Fatalf("could not create temp svg file: %v", err)
	}
	return path
}

func TestMainValidFile(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	path := createTempPNG(t)
	os.Args = []string{"eyez", path}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	go io.Copy(io.Discard, r)

	main()

	w.Close()
	os.Stdout = oldStdout
}

func TestMainValidSVGFile(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	path := createTempSVG(t)
	os.Args = []string{"eyez", path}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	go io.Copy(io.Discard, r)

	main()

	w.Close()
	os.Stdout = oldStdout
}

func TestMainValidSVGPipe(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"eyez"}

	svgData := []byte(`<svg viewBox="0 0 50 50" xmlns="http://www.w3.org/2000/svg"><circle cx="25" cy="25" r="20" fill="green"/></svg>`)

	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	w.Write(svgData)
	w.Close()
	defer func() { os.Stdin = oldStdin }()

	oldStdout := os.Stdout
	rStd, wStd, _ := os.Pipe()
	os.Stdout = wStd

	go io.Copy(io.Discard, rStd)

	main()

	wStd.Close()
	os.Stdout = oldStdout
}

func TestMainValidPipe(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"eyez"}

	var buf bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	png.Encode(&buf, img)

	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	w.Write(buf.Bytes())
	w.Close()
	defer func() { os.Stdin = oldStdin }()

	oldStdout := os.Stdout
	rStd, wStd, _ := os.Pipe()
	os.Stdout = wStd

	go io.Copy(io.Discard, rStd)

	main()

	wStd.Close()
	os.Stdout = oldStdout
}

func TestMainInvalidFileExit(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"eyez", "nonexistent_file.png"}
		main()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMainInvalidFileExit")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestMainInvalidStdinExit(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		os.Args = []string{"eyez"}
		main()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMainInvalidStdinExit")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	cmd.Stdin = strings.NewReader("not a real image")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
