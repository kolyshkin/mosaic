package mosaic

// A test suite

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"testing"
)

var (
	oldPwd  string
	testDir string
	drivers []string
)

func init() {
	drivers = []string{"plain", "fsimg", "btrfs", "ploop"}
}

// abort is used when the tests after the current one
// can't be run as one of their prerequisite(s) failed
func abort(format string, args ...interface{}) {
	s := fmt.Sprintf("ABORT: "+format+"\n", args...)
	f := bufio.NewWriter(os.Stderr)
	f.Write([]byte(s))
	f.Flush()
	dirCleanup()
	os.Exit(1)
}

// Check for a fatal error, call abort() if it is
func chk(err error) {
	if err != nil {
		abort("%s", err)
	}
}

func dirPrepare(dir string) {
	var err error

	oldPwd, err = os.Getwd()
	chk(err)

	testDir, err = ioutil.TempDir(oldPwd, dir)
	chk(err)

	err = os.Chdir(testDir)
	chk(err)
}

func dirCleanup() {
	if oldPwd != "" {
		os.Chdir(oldPwd)
	}
	if testDir != "" {
		os.RemoveAll(testDir)
	}
}

type skipErr struct {
}

func (e *skipErr) Error() string {
	return "driver not supported"
}

func isSkipErr(e error) bool {
	_, ok := e.(*skipErr)
	return ok
}

func runDriverShellFn(drv, fn string) (string, error) {
	mosfile := "./" + drv + ".mos"
	script := "../../test/" + drv + ".sh"

	cmd := exec.Command("bash", "-c", "source "+script+"; "+fn)
	//	cmd.Stdout = os.Stdout // be verbose
	cmd.Stderr = os.Stderr // always show errors
	env := os.Environ()
	env = append(env, fmt.Sprintf("MOS_FILE=%s", mosfile))
	cmd.Env = env

	err := cmd.Run()

	if err == nil {
		return mosfile, nil
	}

	// Check for a specific exit code of 22 (meaning "driver not supported")

	// Get the exit code (Unix-specific)
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			errCode := status.ExitStatus()
			if errCode == 22 {
				return mosfile, &skipErr{}
			}
		}
	}

	return "", err
}

func mosPrepare(drv string) (s string, e error) {
	dirPrepare("tmp-test-" + drv)

	s, e = runDriverShellFn(drv, "prepare")
	if e != nil {
		dirCleanup()
	}

	return
}

func mosCleanup(drv string) {
	runDriverShellFn(drv, "cleanup")
	dirCleanup()
}

func testMosaicMountUmount(t *testing.T, drv string) {
	t.Logf("Mosaic mount/umount test for %s", drv)

	mosfile, err := mosPrepare(drv)
	if isSkipErr(err) {
		// skip the test for this driver
		t.Logf("SKIP %s: %s", drv, err)
		return
	}
	chk(err)
	defer mosCleanup(drv)

	mntdir := "mmnt"
	err = os.Mkdir(mntdir, 0755)
	chk(err)

	m, err := Open(mosfile, 0)
	chk(err)
	defer m.Close()

	err = m.Mount(mntdir, 0)
	chk(err)
	defer m.Umount(mntdir)

	tmpfile := mntdir + "/tfile"
	err = ioutil.WriteFile(tmpfile, []byte("test"), 0644)
	chk(err)

	err = m.Umount(mntdir)
	chk(err)

	_, err = ioutil.ReadFile(tmpfile)
	// expecting ENOENT
	if err == nil || !os.IsNotExist(err) {
		t.Fatalf("Unexpectedly found %s: %s", tmpfile, err)
	}

	err = m.Mount(mntdir, 0)
	chk(err)

	err = os.Remove(tmpfile)
	chk(err)
}

func TestMosaicMountUmount(t *testing.T) {
	for _, drv := range drivers {
		testMosaicMountUmount(t, drv)
	}
}

func bstr(b bool) string {
	if b {
		return "YES"
	}
	return "no "
}

func testMosaicFeatures(t *testing.T, drv string) {
	mosfile, err := mosPrepare(drv)
	if isSkipErr(err) {
		// skip the test for this driver
		t.Logf("SKIP %s: %s", drv, err)
		return
	}
	chk(err)
	defer mosCleanup(drv)

	m, err := Open(mosfile, 0)
	chk(err)
	defer m.Close()

	res := fmt.Sprintf("%s features: 0x%x (", drv, m.Features())
	res += fmt.Sprintf(" clone: %s ", bstr(m.CanClone()))
	res += fmt.Sprintf(" size: %s ", bstr(m.CanManageSize()))
	res += fmt.Sprintf(" bdev: %s ", bstr(m.CanBlockDev()))
	res += fmt.Sprintf(" migrate: %s ", bstr(m.CanMigrate()))
	res += ")"
	t.Logf(res)
}

func TestMosaicFeatures(t *testing.T) {
	for _, drv := range drivers {
		testMosaicFeatures(t, drv)
	}
}

func testVolumeMountUmount(t *testing.T, drv string) {
	vol := "test_vol"
	var size uint64 = 1024 * 1024 // 512MB in 512-byte blocks
	mosfile, err := mosPrepare(drv)
	if isSkipErr(err) {
		// skip the test for this driver
		t.Logf("SKIP %s: %s", drv, err)
		return
	}
	chk(err)
	defer mosCleanup(drv)

	m, err := Open(mosfile, 0)
	chk(err)
	defer m.Close()

	have, err := m.HaveVol(vol, 0)
	chk(err)
	if have {
		t.Fatalf("HaveVol unexpectedly returned true")
	}

	err = m.CreateVol(vol, size, 0, true)
	chk(err)

	have, err = m.HaveVol(vol, 0)
	chk(err)
	if !have {
		t.Fatalf("HaveVol unexpectedly returned false")
	}

	v, err := m.OpenVol(vol, 0)
	chk(err)
	defer v.Close()

	mntdir := "vol_mnt"
	err = os.Mkdir(mntdir, 0755)
	chk(err)

	err = v.Mount(mntdir, 0)
	chk(err)
	defer v.Umount(mntdir, 0)

	tmpfile := mntdir + "/tfile"
	err = ioutil.WriteFile(tmpfile, []byte("test"), 0644)
	chk(err)

	err = v.Umount(mntdir, 0)
	chk(err)

	_, err = ioutil.ReadFile(tmpfile)
	// expecting ENOENT
	if err == nil || !os.IsNotExist(err) {
		t.Fatalf("Unexpectedly found %s: %s", tmpfile, err)
	}

	err = v.Mount(mntdir, 0)
	chk(err)

	os.Remove(tmpfile)

	err = v.Umount(mntdir, 0)
	chk(err)

	err = v.Drop(0)
	chk(err)
}

func TestVolumeMountUmount(t *testing.T) {
	for _, drv := range drivers {
		testVolumeMountUmount(t, drv)
	}
}
