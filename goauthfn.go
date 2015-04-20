package gosmbclient

// #include <stdlib.h>
// #include <string.h>
import "C"
import (
	"unsafe"
)

type Auth struct {
	Wg, Un, Pw string
}
type AuthFn func(server, share, workgroup, username string) (*Auth, error)
type CAuthFn func(srv *C.char, shr *C.char, wg *C.char, wglen C.int, un *C.char, unlen C.int, pw *C.char, pwlen C.int)

var serverAuth = make(map[string]Auth)
var authenticationFunction AuthFn

//export goAuthFn
func goAuthFn(srv *C.char, shr *C.char, wg *C.char, wglen C.int,
	un *C.char, unlen C.int, pw *C.char, pwlen C.int) {
	if authenticationFunction != nil {
		wrapAuthFn(authenticationFunction)(srv, shr, wg, wglen, un, unlen, pw, pwlen)
	}
}

func ccopy(goVal string, cVal *C.char, cValLen C.int) {
	cgoVal := C.CString(goVal)
	defer C.free(unsafe.Pointer(cgoVal))
	C.memset(unsafe.Pointer(cVal), 0, C.size_t(cValLen))
	C.strncpy(cVal, cgoVal, C.size_t(cValLen))
}

func wrapAuthFn(authfn AuthFn) CAuthFn {
	return func(srv *C.char, shr *C.char, wg *C.char, wglen C.int, un *C.char, unlen C.int, pw *C.char, pwlen C.int) {
		server := C.GoString(srv)
		share := C.GoString(shr)
		workgroup := C.GoString(wg)
		username := C.GoString(un)
		auth, err := authfn(server, share, workgroup, username)
		if err != nil {
			return
		}
		ccopy(auth.Wg, wg, wglen)
		ccopy(auth.Un, un, unlen)
		ccopy(auth.Pw, pw, pwlen)
	}
}
