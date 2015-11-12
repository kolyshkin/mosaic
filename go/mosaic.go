package mosaic

// #include <mosaic/mosaic.h>
import "C"
import (
	"fmt"
	"syscall"
)

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

/*
 * Volume functions
 */

// OpenVol opens a volume
func (m Mosaic) OpenVol(name string, flags int) (v Volume, e error) {
	cname := C.CString(name)
	defer cfree(cname)

	e = nil
	v.v = C.mosaic_open_vol(m.m, cname, C.int(flags))
	if v.v == nil {
		e = lastError()
	}

	return
}

// Close closes a volume
func (v Volume) Close() {
	C.mosaic_close_vol(v.v)
}

// HaveVol checks if a volume with given name exists in a mosaic.
func (m Mosaic) HaveVol(name string, flags int) (bool, error) {
	cname := C.CString(name)
	defer cfree(cname)

	ret := C.mosaic_have_vol(m.m, cname, C.int(flags))
	if ret < 0 {
		return false, lastError()
	}

	return ret == 1, nil
}

// CreateVol creates a new volume. If fs argument is true, a filesystem
// is created, otherwise just a block device.
func (m Mosaic) CreateVol(name string, size uint64, flags int, fs bool) error {
	cname := C.CString(name)
	defer cfree(cname)
	csize := C.ulong(size)
	cflags := C.int(flags)
	var ret C.int

	if fs {
		ret = C.mosaic_make_vol_fs(m.m, cname, csize, cflags)
	} else {
		ret = C.mosaic_make_vol(m.m, cname, csize, cflags)
	}

	return condErr(ret)
}

// Clone creates a clone of existing volume v.
func (v Volume) Clone(name string, flags int) error {
	cname := C.CString(name)
	defer cfree(cname)

	ret := C.mosaic_clone_vol(v.v, cname, C.int(flags))

	return condErr(ret)
}

// Drop removes a volume.
func (v Volume) Drop(flags int) error {
	return condErr(C.mosaic_drop_vol(v.v, C.int(flags)))
}

// Resize changes the size of the volume
func (v Volume) Resize(size uint64, flags int) error {
	return condErr(C.mosaic_resize_vol(v.v, C.ulong(size), C.int(flags)))
}

// Mount mounts a volume, i.e. makes it filesystem available at mnt
func (v Volume) Mount(path string, flags int) error {
	cpath := C.CString(path)
	defer cfree(cpath)

	return condErr(C.mosaic_mount_vol(v.v, cpath, C.int(flags)))
}

// Umount unmounts a volume
func (v Volume) Umount(path string, flags int) error {
	cpath := C.CString(path)
	defer cfree(cpath)

	return condErr(C.mosaic_umount_vol(v.v, cpath, C.int(flags)))
}

// GetBlockDev prepares volume's block device and returns path to it
func (v Volume) GetBlockDev(flags int) (string, error) {
	const len = 4096 // PATH_MAX
	var out [len]C.char

	ret := int(C.mosaic_get_vol_bdev(v.v, &out[0], len, C.int(flags)))
	if ret < 0 {
		return "", lastError()
	}
	if ret >= len {
		return "", fmt.Errorf("mosaic.GetBlockDev error: %d >= %d", ret, len)
	}

	path := C.GoString(&out[0])
	return path, nil
}

// PutBlockDev releases volume's block device
func (v Volume) PutBlockDev() error {
	return condErr(C.mosaic_put_vol_bdev(v.v))
}

// Size returns a size of block device, in 512-byte blocks
func (v Volume) Size() (uint64, error) {
	var csize C.ulong

	ret := C.mosaic_get_vol_size(v.v, &csize)
	if ret < 0 {
		return 0, lastError()
	}

	return uint64(csize), nil
}
