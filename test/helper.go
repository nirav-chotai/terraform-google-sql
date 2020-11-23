package test

// helper go file has all the helper functions, needed for testbase.go and the go tests.
import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/logger"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"testing"
)

// WhereAmI function returns the calling function's name that calls it and
// the line number from where it is called from the calling function
func WhereAmI(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, _, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("Function: %s Line: %d", runtime.FuncForPC(function).Name(), line)
}

// UpdateTerraformDirectory function will replace the value of 'source' from the 'filePath' with the 'modulePath' value
func UpdateTerraformDirectory(modulePath string, filePath string) error {
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	lines := strings.Split(string(input), "\n")
	moduleUpdated := false

	for i, line := range lines {
		if strings.Contains(line, "source") && (moduleUpdated == false) {
			moduleUpdated = true
			lines[i] = modulePath
			break
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filePath, []byte(output), 0644)
	return err
}

// UpdateModuleVariable function will add variable to the module under test
func UpdateModuleVariable(terraformVar string, filePath string) error {
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	lines := strings.Split(string(input), "\n")
	fileContent := ""
	variableUpdated := false

	for _, line := range lines {
		fileContent += line
		fileContent += "\n"
		if strings.Contains(line, "source") && (variableUpdated == false) {
			variableUpdated = true
			fileContent += terraformVar
			fileContent += "\n"
		}
	}

	err = ioutil.WriteFile(filePath, []byte(fileContent), 0644)
	return err
}

// RemoveModuleVariable function will remove variable from the module under test
func RemoveModuleVariable(terraformVar string, filePath string) error {
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err.Error())
	}

	lines := strings.Split(string(input), "\n")
	fileContent := ""

	for _, line := range lines {
		if strings.Contains(line, "source") {
			fileContent += line
			fileContent += "\n"
		}
	}

	err = ioutil.WriteFile(filePath, []byte(fileContent), 0644)
	return err
}

// CheckPanic will log the error and end the test
func CheckPanic(t *testing.T, e error, failTrace string, logMsg string) {
	if e != nil {
		logger.Logf(t, failTrace)
		logger.Logf(t, "PANIC: "+logMsg+": %s", e)
		TestPassed = false
		errorStatus = e
		panic(e)
	}
}

// CheckError will log the error at the end of the test, but will continue execution and mark it as failure in the end.
func CheckError(t *testing.T, e error, failTrace string, logMsg string) bool {
	if e != nil {
		logger.Logf(t, failTrace)
		logger.Logf(t, "ERROR: "+logMsg+": %s", e)
		TestPassed = false
		errorStatus = e
		t.Errorf("ERROR: "+logMsg+": %s", e)
		return false
	}
	return true
}

// CheckFatal will fail the test immediately
func CheckFatal(t *testing.T, e error, failTrace string, logMsg string) {
	if e != nil {
		logger.Logf(t, failTrace)
		logger.Logf(t, "FATAL: "+logMsg+": %s", e)
		TestPassed = false
		errorStatus = e
		afterTestMethod(t)
		t.Fatalf("FATAL: "+logMsg+": %s", e)
	}
}

// Copy a file from one path to another.copyFileContents copies the contents of the file named src to the file named by dst.
// The file will be created if it does not already exist. If the destination file exists,
// all it's contents will be replaced by the contents of the source file.
func CopyFile(src, dst, fileName string) (err error) {
	dst = fmt.Sprintf(("%s/%s"), dst, fileName)
	fileSrc, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("Error in describing the file info of the source file %s (%q)", fileSrc.Name(), fileSrc.Mode().String())
	}
	if !fileSrc.Mode().IsRegular() {
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", fileSrc.Name(), fileSrc.Mode().String())
	}
	fileDest, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(fileDest.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", fileDest.Name(), fileDest.Mode().String())
		}
		if os.SameFile(fileSrc, fileDest) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func testMethodStatus(errorStatus error) string {
	if errorStatus != nil {
		return "FAIL"
	} else {
		return "PASS"
	}
}

// Called to end the test gracefully if panic() is called in the test
func PanicCall(t *testing.T) {
	AfterSuite(t)
	TestStatus()
	os.Exit(1)
}
