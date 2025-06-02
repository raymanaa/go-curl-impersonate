//go:build windows

package curl

// for GlobalInit(flag)
const (
	GLOBAL_SSL     = 1
	GLOBAL_WIN32   = 2
	GLOBAL_ALL     = 3
	GLOBAL_NOTHING = 0
	GLOBAL_DEFAULT = 3
)

// CURLMcode (from multi.h, typically CURLM_OK = 0, CURLM_CALL_MULTI_PERFORM = -1)
const (
	M_CALL_MULTI_PERFORM = -1
	M_OK                 = 0
	M_BAD_HANDLE         = 1
	M_BAD_EASY_HANDLE    = 2
	M_OUT_OF_MEMORY      = 3
	M_INTERNAL_ERROR     = 4
	M_BAD_SOCKET         = 5
	M_UNKNOWN_OPTION     = 6
)

// for multi.Setopt(flag, ...) (CURLMOPT_*)
// CURLOPTTYPE_LONG = 0, CURLOPTTYPE_OBJECTPOINT = 10000, CURLOPTTYPE_FUNCTIONPOINT = 20000
const (
	MOPT_SOCKETFUNCTION = 20000 + 1
	MOPT_SOCKETDATA     = 10000 + 2
	MOPT_PIPELINING     = 0 + 3
	MOPT_TIMERFUNCTION  = 20000 + 4
	MOPT_TIMERDATA      = 10000 + 5
	MOPT_MAXCONNECTS    = 0 + 6
)

// CURLSHcode
const (
	SHE_OK         = 0
	SHE_BAD_OPTION = 1
	SHE_IN_USE     = 2
	SHE_INVALID    = 3
	SHE_NOMEM      = 4
)

// for share.Setopt(flag, ...) (CURLSHOPT_*)
const (
	SHOPT_SHARE      = 1
	SHOPT_UNSHARE    = 2
	SHOPT_LOCKFUNC   = 3
	SHOPT_UNLOCKFUNC = 4
	SHOPT_USERDATA   = 5
)

// for share.Setopt(SHOPT_SHARE/SHOPT_UNSHARE, flag) (CURL_LOCK_DATA_*)
const (
	LOCK_DATA_SHARE       = 1
	LOCK_DATA_COOKIE      = 2
	LOCK_DATA_DNS         = 3
	LOCK_DATA_SSL_SESSION = 4
	LOCK_DATA_CONNECT     = 5
)

// for VersionInfo(flag) (CURLVERSION_*)
const (
	VERSION_FIRST  = 0
	VERSION_SECOND = 1
	VERSION_THIRD  = 2
	// VERSION_FOURTH = 3
	VERSION_LAST = 12
	VERSION_NOW  = 11
)

// for VersionInfo(...).Features mask flag (CURL_VERSION_*)
const (
	VERSION_IPV6         = 1 << 0
	VERSION_KERBEROS4    = 1 << 1
	VERSION_SSL          = 1 << 2
	VERSION_LIBZ         = 1 << 3
	VERSION_NTLM         = 1 << 4
	VERSION_GSSNEGOTIATE = 1 << 5
	VERSION_DEBUG        = 1 << 6
	VERSION_ASYNCHDNS    = 1 << 7
	VERSION_SPNEGO       = 1 << 8
	VERSION_LARGEFILE    = 1 << 9
	VERSION_IDN          = 1 << 10
	VERSION_SSPI         = 1 << 11
	VERSION_CONV         = 1 << 12
	VERSION_CURLDEBUG    = 1 << 13
	VERSION_TLSAUTH_SRP  = 1 << 14
	VERSION_NTLM_WB      = 1 << 15
)

// for OPT_READFUNCTION, return a int flag (CURL_READFUNC_*)
const (
	WRITEFUNC_PAUSE = 0x10000001
	READFUNC_ABORT  = 0x10000000
	READFUNC_PAUSE  = 0x10000001
)

// for easy.Setopt(OPT_HTTP_VERSION, flag) (CURL_HTTP_VERSION_*)
const (
	HTTP_VERSION_NONE = 0
	HTTP_VERSION_1_0  = 1
	HTTP_VERSION_1_1  = 2
)

// for easy.Setopt(OPT_PROXYTYPE, flag) (CURLPROXY_*)
const (
	PROXY_HTTP            = 0
	PROXY_HTTP_1_0        = 1
	PROXY_SOCKS4          = 4
	PROXY_SOCKS5          = 5
	PROXY_SOCKS4A         = 6
	PROXY_SOCKS5_HOSTNAME = 7
)

// for easy.Setopt(OPT_SSLVERSION, flag) (CURL_SSLVERSION_*)
const (
	SSLVERSION_DEFAULT = 0
	SSLVERSION_TLSv1   = 1
	SSLVERSION_SSLv2   = 2
	SSLVERSION_SSLv3   = 3
)

// for easy.Setopt(OPT_TIMECONDITION, flag) (CURL_TIMECOND_*)
const (
	TIMECOND_NONE         = 0
	TIMECOND_IFMODSINCE   = 1
	TIMECOND_IFUNMODSINCE = 2
	TIMECOND_LASTMOD      = 3
)

// for easy.Setopt(OPT_NETRC, flag) (CURL_NETRC_*)
const (
	NETRC_IGNORED  = 0
	NETRC_OPTIONAL = 1
	NETRC_REQUIRED = 2
)

// for easy.Setopt(OPT_FTP_CREATE_MISSING_DIRS, flag) (CURLFTP_CREATE_DIR_*)
const (
	FTP_CREATE_DIR_NONE  = 0
	FTP_CREATE_DIR       = 1
	FTP_CREATE_DIR_RETRY = 2
)

// for easy.Setopt(OPT_IPRESOLVE, flag) (CURL_IPRESOLVE_*)
const (
	IPRESOLVE_WHATEVER = 0
	IPRESOLVE_V4       = 1
	IPRESOLVE_V6       = 2
)

// for easy.Setopt(OPT_SSL_OPTIONS, flag) (CURLSSLOPT_*)
const (
	SSLOPT_ALLOW_BEAST = 1
)

// for easy.Pause(flag) (CURLPAUSE_*)
const (
	PAUSE_RECV      = 1
	PAUSE_RECV_CONT = 0
	PAUSE_SEND      = 4
	PAUSE_SEND_CONT = 0
	PAUSE_ALL       = 5
	PAUSE_CONT      = 0
)

// for multi.Info_read() (CURLMSG_*)
const (
	CURLMSG_NONE = 0
	CURLMSG_DONE = 1
	CURLMSG_LAST = 2
)

// CURLformoption values for curl_formadd
const (
	FORM_COPYNAME       = 1
	FORM_COPYCONTENTS   = 4
	FORM_CONTENTSLENGTH = 6
	FORM_FILE           = 10
	FORM_CONTENTTYPE    = 14
	FORM_END            = 17
)

// CURLFORMcode return values for curl_formadd
const (
	FORMADD_OK     = 0
	FORMADD_MEMORY = 1
)
