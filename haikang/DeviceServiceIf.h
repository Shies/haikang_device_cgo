#ifndef DEVICESERVICEIF_H
#define DEVICESERVICEIF_H

#include "Types.h"
#include "DeviceService.h"

class DeviceServiceIf
{
protected:
    int port;
    int channel;
    char address[NET_DVR_DEV_ADDRESS_MAX_LEN];
    char username[NAME_LEN];
    char password[NAME_LEN];
    Scheme *scheme;

public:
    ~DeviceServiceIf();
    DeviceServiceIf();
    DeviceServiceIf(Scheme *scheme);
    Scheme init(std::string addr, std::string user, std::string pass, int port, int channel);

    // common
    int clean();
    int SDKInit();
    int ping(void *in, ExceptionCallBack g_ExceptionCallBack);

    // noAuth
    int login(BOOL async, void *dto);
    int active(LONG lUserID, LONG *status);
    int logout(LONG lUserID);
    int localIp(void *dto);
    int iPByResolveSvr(void *in, void *dto);

    // auth
    int getSadpInfoList(LONG lUserID, void *lpSadpInfoList);
    int stdXmlConfig(LONG lUserID, char *url, char *inbuf, char *outbuf, char *statbuf);
    int setupAlarmChan(LONG lUserID, LONG *returnHandle, MessageCallback callback);
    int closeAlarmChan(LONG lHandle);
    int getDVRConfig(LONG lUserID, LONG lChannel, void *dto);
    int realPlay(void *in, StdDataCallBack fStdDataCallBack);
    int realStopPlay(LONG lRealPlayHandle);
    int realCapPicture(LONG lRealPlayHandle, void *dto);
    int capPicture(LONG lUserID, LONG lChannel, void *dto);
    int getPicture(LONG lUserID, char *sDVRFileName, void *dto);
};

#endif
