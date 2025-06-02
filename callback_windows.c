#include <stddef.h>

extern int GoXferInfoCallback(void *clientp,
                              double dltotal,
                              double dlnow,
                              double ultotal,
                              double ulnow);

int c_xferinfo_shim_entrypoint(void *clientp,
                               double dltotal,
                               double dlnow,
                               double ultotal,
                               double ulnow) {
    return GoXferInfoCallback(clientp, dltotal, dlnow, ultotal, ulnow);
}