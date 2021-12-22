#ifndef DEVICEHANDLER_H
#define DEVICEHANDLER_H

#include "Types.h"
#include "_obj/_cgo_export.h"

#ifdef __cplusplus
extern "C" {
#endif

// The gateway function
void catchErrorCallback_cgo(DWORD dwType, LONG lUserID, LONG lHandle, void *pUser);
void realDataCallBack_cgo(LONG lRealHandle, DWORD dwDataType, BYTE *pBuffer, DWORD dwBufSize, void *pUser);
BOOL alarmMsgCallback_cgo(LONG lCommand, NET_DVR_ALARMER *pAlarmer, char *pAlarmInfo, DWORD dwBufLen, void *pUser);

// noAuth
void clean();
void SDKInit();
void checkHealth(HealthParam *param, ExceptionCallBack fn);
void login(Scheme *scheme, BOOL async, LoginDeviceDto *dto);
void active(LONG lUserID, LONG *status);
void logout(LONG lUserID);
void localIp(LocalIpDto *dto);
void iPByResolveSvr(ResolveSvrParam *in, ResolveSvrDto *dto);

// auth
void getSadpInfoList(LONG lUserID, NET_DVR_SADPINFO_LIST *lpSadpInfoList);
void stdXmlConfig(LONG lUserID, char *url, char *inbuf, char *outbuf, char *statbuf);
void setupAlarmChan(LONG lUserID, LONG *returnHandle, MessageCallback callback);
void closeAlarmChan(LONG lHandle);
void getDVRConfig(LONG lUserID, LONG lChannel, NET_DVR_IPPARACFG_V40 *dto);
void realPlay(RealPlayParam *in, StdDataCallBack fn);
void realStopPlay(LONG lRealPlayHandle);
void realCapPicture(LONG lRealPlayHandle, RealCapPictureDto *dto);
void capPicture(LONG lUserID, LONG lChannel, CapPictureDto *dto);
void getPicture(LONG lUserID, char *sDVRFileName, CapPictureDto *dto);

#ifdef __cplusplus
}
#endif

#endif