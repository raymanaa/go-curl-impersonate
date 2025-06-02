//go:build windows

package curl

/*
int c_xferinfo_shim_entrypoint(void *clientp,
double dltotal,
double dlnow,
double ultotal,
double ulnow);
*/
import "C"
import "unsafe"

//export GoXferInfoCallback
func GoXferInfoCallback(clientp unsafe.Pointer, dltotal C.double, dlnow C.double, ultotal C.double, ulnow C.double) C.int {
	gdltotal := float64(dltotal)
	gdlnow := float64(dlnow)
	gultotal := float64(ultotal)
	gulnow := float64(ulnow)

	curlHandle := context_map.Get(uintptr(clientp))

	if curlHandle == nil {
		return 0
	}

	if curlHandle.progressFunction == nil {
		return 0
	}

	if (*curlHandle.progressFunction)(gdltotal, gdlnow, gultotal, gulnow, curlHandle.progressData) {
		return 0
	}
	return 1
}
