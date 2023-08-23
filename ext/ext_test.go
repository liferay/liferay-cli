package ext

import (
	"regexp"
	"testing"
)

func TestMakeWorkspaceDirKeyLinux(t *testing.T) {
	sep := "/"
	baseDirName := "foo"
	dirPath := "/home/user"
	doTest(baseDirName, dirPath, sep, t)
}

func TestMakeWorkspaceDirKeyWindows(t *testing.T) {
	sep := "\\"
	baseDirName := "foo"
	dirPath := "C:\\home\\user"
	doTest(baseDirName, dirPath, sep, t)
}

func doTest(baseDirName string, dirPath string, sep string, t *testing.T) {
	want := regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_.-]*$")
	result := MakeWorkspaceDirKey(baseDirName, dirPath, sep)
	if !want.MatchString(result) {
		t.Fatalf(`MakeWorkspaceDirKey("%s") = %q, want match for %#q`, dirPath+sep+baseDirName, result, want)
	}
	t.Logf(`MakeWorkspaceDirKey("%s") = %q`, dirPath+sep+baseDirName, result)
}
