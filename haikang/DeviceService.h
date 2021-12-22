#ifndef DEVICESERVICE_H
#define DEVICESERVICE_H

#define HPR_OK 0
#define HPR_ERROR -1

#if (defined(_WIN32)) //windows
#include <windows.h>
#endif

#ifndef _HC_NET_SDK_H_
#define NET_DVR_DEV_ADDRESS_MAX_LEN 129
#define NAME_LEN 32
#endif

/*
stub class
class DeviceService
{
public:
    // common
    virtual int clean();
    virtual int SDKInit();
    virtual int ping(void *in, ExceptionCallBack g_ExceptionCallBack);

private:
    // noAuth
    virtual int login(BOOL async, void *dto);
    virtual int active(LONG lUserID, LONG *status);
    virtual int logout(LONG lUserID);
    virtual int localIp(void *dto);
    virtual int iPByResolveSvr(void *in, void *dto);

protected:
    void catchError(ExceptionCallBack g_ExceptionCallBack) {
        ping(NULL, g_ExceptionCallBack);
    };

    void checkHealth(void *in) {
        ping(in, NULL);
    };
};
*/

#endif
