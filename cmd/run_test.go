package cmd

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCopyScriptForDebug(t *testing.T) {
	content := "temporary file's content"
	tmpfile, err := ioutil.TempFile("", "pogo-test")
	if err != nil {
		t.Error("Failed to create tempfile:", err)
	}
	_, err = tmpfile.Write([]byte(content))
	if err != nil {
		t.Error("Failed to write tempfile:", err)
	}
	defer os.Remove(tmpfile.Name())
	debugfile, err := CopyScriptForDebug(tmpfile.Name(), "test")
	if err != nil {
		t.Error("Failed to write debugfile:", err)
	}
	defer os.Remove(debugfile.Name())
	debugContent, err := ioutil.ReadFile(debugfile.Name())
	if err != nil {
		t.Error("Failed to read debugfile:", err)
	}
	if string(debugContent) != content {
		t.Errorf("Debug content %v does not match %v", debugContent, content)
	}
}
