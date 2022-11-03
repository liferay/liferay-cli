package ext

import (
	"regexp"
	"testing"
)

func TestMakeExtensionDirKeyLinux(t *testing.T) {
	sep := "/"
	baseDirName := "foo"
	dirPath := "/home/user"
	doTest(baseDirName, dirPath, sep, t)
}

func TestMakeExtensionDirKeyWindows(t *testing.T) {
	sep := "\\"
	baseDirName := "foo"
	dirPath := "C:\\home\\user"
	doTest(baseDirName, dirPath, sep, t)
}

func doTest(baseDirName string, dirPath string, sep string, t *testing.T) {
	want := regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_.-]*$")
	result := MakeExtensionDirKey(baseDirName, dirPath, sep)
	if !want.MatchString(result) {
		t.Fatalf(`MakeExtensionDirKey("%s") = %q, want match for %#q`, dirPath+sep+baseDirName, result, want)
	}
	t.Logf(`MakeExtensionDirKey("%s") = %q`, dirPath+sep+baseDirName, result)
}
