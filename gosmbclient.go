package gosmbclient

/*
#include <samba-4.0/libsmbclient.h>
#include <stdlib.h>

void goAuthFn(const char *srv, const char *shr, char *wg,
int wglen, char *un, int unlen, char *pw, int pwlen);

#cgo LDFLAGS: -lsmbclient
*/
import "C"
import (
	"io"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	_                = iota
	WORKGROUP C.uint = iota
	SERVER
	FILE_SHARE
	PRINTER_SHARE
	COMMS_SHARE
	IPC_SHARE
	DIR
	FILE
	LINK
)

func isDirOrFile(t C.uint) bool {
	return t == DIR || t == FILE
}

func addSlash(url string) string {
	if url[len(url)-1] == '/' {
		return url
	}
	return url + "/"
}

type Dir struct {
	fd  C.uint
	url string
}

func (d Dir) Close() {
	C.smbc_close(C.int(d.fd))
}

func (d Dir) list() ([]string, error) {
	size := 4096
	buf := C.malloc(C.size_t(size))
	if buf == nil {
		return nil, syscall.ENOMEM
	}
	defer C.free(unsafe.Pointer(buf))
	max, err := C.smbc_getdents(d.fd, (*C.struct_smbc_dirent)(buf), C.int(size))
	if err != nil {
		return nil, err
	}
	ret := []string{}
	current := unsafe.Pointer(buf)
	for (uintptr(current) - uintptr(buf)) < uintptr(max) {
		var d *C.struct_smbc_dirent
		d = (*C.struct_smbc_dirent)(current)
		if isDirOrFile(d.smbc_type) {
			ret = append(ret, C.GoString((*C.char)(&d.name[0])))
		}
		current = unsafe.Pointer(uintptr(current) + uintptr(d.dirlen))
	}
	return ret, nil
}

func (d Dir) List() ([]string, error) {
	names := []string{}
	for {
		list, err := d.list()
		if err != nil {
			return nil, err
		}
		if len(list) == 0 {
			break
		}
		names = append(names, list...)
	}
	for i, name := range names {
		names[i] = d.url + name
	}
	return names, nil
}

func Init(authfn AuthFn, debugLevel int) error {
	authenticationFunction = authfn
	auth := C.smbc_get_auth_data_fn(C.goAuthFn)
	ec, err := C.smbc_init(auth, C.int(debugLevel))
	if ec != 0 {
		return err
	}
	return nil
}

func OpenDir(durl string) (*Dir, error) {
	cdurl := C.CString(durl)
	defer C.free(unsafe.Pointer(cdurl))
	fd, err := C.smbc_opendir(cdurl)
	if fd < 0 {
		return nil, err
	}
	return &Dir{C.uint(fd), addSlash(durl)}, nil
}

func MkDir(durl string, mode os.FileMode) error {
	cdurl := C.CString(durl)
	defer C.free(unsafe.Pointer(cdurl))
	ec, err := C.smbc_mkdir(cdurl, C.mode_t(mode))
	if ec != 0 {
		return err
	}
	return nil
}

func RmDir(durl string) error {
	cdurl := C.CString(durl)
	defer C.free(unsafe.Pointer(cdurl))
	ec, err := C.smbc_rmdir(cdurl)
	if ec != 0 {
		return err
	}
	return nil
}

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (f fileInfo) Name() string       { return f.name }
func (f fileInfo) Size() int64        { return f.size }
func (f fileInfo) Mode() os.FileMode  { return f.mode }
func (f fileInfo) ModTime() time.Time { return f.modTime }
func (f fileInfo) IsDir() bool        { return f.mode.IsDir() }
func (f fileInfo) Sys() interface{}   { return nil }

func Stat(durl string) (os.FileInfo, error) {
	cdurl := C.CString(durl)
	defer C.free(unsafe.Pointer(cdurl))
	var st C.struct_stat
	ec, err := C.smbc_stat(cdurl, &st)
	if ec != 0 {
		return nil, err
	}
	// TODO(klak): isdir
	t := time.Unix(int64(st.st_mtim.tv_sec), int64(st.st_mtim.tv_nsec))
	mode := os.FileMode(st.st_mode)
	info := fileInfo{durl, int64(st.st_size), mode, t, false}

	return info, nil
}

func Unlink(durl string) error {
	cdurl := C.CString(durl)
	defer C.free(unsafe.Pointer(cdurl))
	ec, err := C.smbc_unlink(cdurl)
	if ec < 0 {
		return err
	}
	return nil
}

type File struct {
	fd  C.int
	url string
}

func Open(durl string, flag int, perm os.FileMode) (*File, error) {
	cdurl := C.CString(durl)
	defer C.free(unsafe.Pointer(cdurl))
	fd, err := C.smbc_open(cdurl, C.int(flag), C.mode_t(perm))
	if fd < 0 {
		return nil, err
	}
	return &File{fd, durl}, nil
}

func Create(durl string, perm os.FileMode) (*File, error) {
	cdurl := C.CString(durl)
	defer C.free(unsafe.Pointer(cdurl))
	fd, err := C.smbc_creat(cdurl, C.mode_t(perm))
	if fd < 0 {
		return nil, err
	}
	return &File{fd, durl}, nil
}

func (f File) Close() {
	// TODO(klak): errors
	C.smbc_close(f.fd)
}

func (f File) Read(buf []byte) (int, error) {
	l := C.size_t(len(buf))
	arr := &buf[0]
	size, err := C.smbc_read(f.fd, unsafe.Pointer(arr), l)
	if size < 0 {
		return int(size), err
	} else if size == 0 {
		return 0, io.EOF
	}
	return int(size), nil
}

func (f File) Write(buf []byte) (int, error) {
	l := C.size_t(len(buf))
	arr := unsafe.Pointer(&buf[0])
	size, err := C.smbc_write(f.fd, arr, l)
	return int(size), err
}

func Rename(src, dst string) error {
	csrc := C.CString(src)
	defer C.free(unsafe.Pointer(csrc))
	cdst := C.CString(dst)
	defer C.free(unsafe.Pointer(cdst))
	if ec, err := C.smbc_rename(csrc, cdst); ec != 0 {
		return err
	}
	return nil
}
