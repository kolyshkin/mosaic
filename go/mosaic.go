package mosaic

// #include <mosaic/mosaic.h>
import "C"
import "syscall"

// Mosaic is a type holding a reference to a mosaic
type Mosaic struct {
	m *C.struct_mosaic
}

// Volume is a type holding a reference to a volume
type Volume struct {
	v *C.struct_volume
}

// Open opens a mosaic
func Open(name string, flags int) (m Mosaic, e error) {
	cname := C.CString(name)
	defer cfree(cname)

	e = nil
	m.m = C.mosaic_open(cname, C.int(flags))
	if m.m == nil {
		e = lastError()
	}

	return
}

// Close closes a mosaic
func (m Mosaic) Close() {
	C.mosaic_close(m.m)
}

// Mount mounts a mosaic
func (m Mosaic) Mount(path string, flags int) error {
	cpath := C.CString(path)
	defer cfree(cpath)

	ret := C.mosaic_mount(m.m, cpath, C.int(flags))
	if ret == 0 {
		return nil
	}

	return lastError()
}

// Umount umounts a mosaic
func (m Mosaic) Umount(path string) error {
	return syscall.Unmount(path, 0)
}

/*
 * Mosaic features
 */

// Possible bits set by Features()
const (
	FeatureClone      = C.MOSAIC_FEATURE_CLONE
	FeatureManageSize = C.MOSAIC_FEATURE_DISK_SIZE_MGMT
	FeatureBlockDev   = C.MOSAIC_FEATURE_BDEV
	FeatureMigrate    = C.MOSAIC_FEATURE_MIGRATE
)

// Features returns mosaic features
func (m Mosaic) Features() uint64 {
	var f C.ulonglong
	// Assuming mosaic_get_features() never fails
	C.mosaic_get_features(m.m, &f)

	return uint64(f)
}

// CanClone checks if a mosaic support volume cloning
func (m Mosaic) CanClone() bool {
	return (m.Features() & FeatureClone) > 0
}

// CanManageSize checks if a mosaic honors size argument of CreateVol(),
// and supports Resize().
func (m Mosaic) CanManageSize() bool {
	return (m.Features() & FeatureManageSize) > 0
}

// CanBlockDev checks if a mosaic supports raw block devices
// (i.e. GetBlockDev() works).
func (m Mosaic) CanBlockDev() bool {
	return (m.Features() & FeatureBlockDev) > 0
}

// CanMigrate checks if a mosaic supports volume (live) migration.
func (m Mosaic) CanMigrate() bool {
	return (m.Features() & FeatureMigrate) > 0
}
