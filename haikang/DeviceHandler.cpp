#include <stdio.h>
#include <stdlib.h>
#include <iostream>
#include <string.h>
#include "DeviceServiceIf.h"
#include "DeviceHandler.h"
using namespace std;

// The gateway function
void catchErrorCallback_cgo(DWORD dwType, LONG lUserID, LONG lHandle, void *pUser) {
    GoCatchErrorCallback(dwType, lUserID, lHandle, pUser);
}
void realDataCallBack_cgo(LONG lRealHandle, DWORD dwDataType, BYTE *pBuffer, DWORD dwBufSize, void *pUser) {
    GoRealDataCallBack(lRealHandle, dwDataType, pBuffer, dwBufSize, pUser);
}
BOOL alarmMsgCallback_cgo(LONG lCommand, NET_DVR_ALARMER *pAlarmer, char *pAlarmInfo, DWORD dwBufLen, void *pUser) {
    NET_VCA_FACESNAP_MATCH_ALARM struFaceMatchAlarm = {0};
    memcpy(&struFaceMatchAlarm, pAlarmInfo, sizeof(NET_VCA_FACESNAP_MATCH_ALARM));
    GoAlarmMsgCallback(lCommand, pAlarmer, &struFaceMatchAlarm, dwBufLen, pUser);
    return true;
}

void setupAlarmChan(LONG lUserID, LONG *returnHandle, MessageCallback fn) {
    /*
    NET_VCA_FACESNAP_MATCH_ALARM struFaceMatchAlarm = {0};
    char ip[16] = "127.0.0.1";
    memcpy(struFaceMatchAlarm.sStorageIP, ip, strlen(ip));
    GoAlarmMsgCallback(1, NULL, &struFaceMatchAlarm, 2, NULL);
    */
    DeviceServiceIf *s = new DeviceServiceIf();
    s->setupAlarmChan(lUserID, returnHandle, fn);
    delete s;
    std::cout << "\n" << std::endl;
}

void closeAlarmChan(LONG lHandle) {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->closeAlarmChan(lHandle);
    delete s;
    std::cout << "\n" << std::endl;
}

void SDKInit() {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->SDKInit();
    delete s;
    // std::cout << "\n" << std::endl;
}

void clean() {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->clean();
    delete s;
    // std::cout << "\n" << std::endl;
}

void checkHealth(HealthParam *param, ExceptionCallBack fn) {
    fn(1, 2, 3, NULL);
    DeviceServiceIf *s = new DeviceServiceIf();
    s->ping(param, fn);
    delete s;
    std::cout << "\n" << std::endl;
}

void login(Scheme *scheme, BOOL async, LoginDeviceDto *dto) {
    DeviceServiceIf *s = new DeviceServiceIf(scheme);
    s->login(async, dto);
    delete s;
    //std::cout << "\n" << std::endl;
}

void active(LONG lUserID, LONG *status) {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->active(lUserID, status);
    delete s;
    //std::cout << "\n" << std::endl;
}

void logout(LONG lUserID) {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->logout(lUserID);
    delete s;
    //std::cout << "\n" << std::endl;
}

void localIp(LocalIpDto *dto) {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->localIp(dto);
    delete s;
    //std::cout << "\n" << std::endl;
}

void iPByResolveSvr(ResolveSvrParam *in, ResolveSvrDto *dto) {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->iPByResolveSvr(in, dto);
    delete s;
    //std::cout << "\n" << std::endl;
}

void getDVRConfig(LONG lUserID, LONG lChannel, NET_DVR_IPPARACFG_V40 *dto) {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->getDVRConfig(lUserID, lChannel, dto);
    delete s;
    //std::cout << "\n" << std::endl;
}

void realPlay(RealPlayParam *in, StdDataCallBack fn) {
    fn(4, 5, NULL, 6, NULL);
    DeviceServiceIf *s = new DeviceServiceIf();
    s->realPlay(in, fn);
    delete s;
    //std::cout << "\n" << std::endl;
}

void realStopPlay(LONG lRealPlayHandle) {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->realStopPlay(lRealPlayHandle);
    delete s;
    //std::cout << "\n" << std::endl;
}

void realCapPicture(LONG lRealPlayHandle, RealCapPictureDto *dto) {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->realCapPicture(lRealPlayHandle, dto);
    delete s;
    //std::cout << "\n" << std::endl;
}

void capPicture(LONG lUserID, LONG lChannel, CapPictureDto *dto) {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->capPicture(lUserID, lChannel, dto);
    delete s;
    //std::cout << "\n" << std::endl;
}

void stdXmlConfig(LONG lUserID, char *url, char *inbuf, char *outbuf, char *statbuf) {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->stdXmlConfig(lUserID, url, inbuf, outbuf, statbuf);
    delete s;
    //std::cout << "\n" << std::endl;
}

void getSadpInfoList(LONG lUserID, NET_DVR_SADPINFO_LIST *lpSadpInfoList) {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->getSadpInfoList(lUserID, lpSadpInfoList);
    delete s;
    //std::cout << "\n" << std::endl;
}

void getPicture(LONG lUserID, char *sDVRFileName, CapPictureDto *dto) {
    DeviceServiceIf *s = new DeviceServiceIf();
    s->getPicture(lUserID, sDVRFileName, dto);
    delete s;
    //std::cout << "\n" << std::endl;
}

/*
int main()
{
    Scheme scheme;
    DeviceServiceIf *s = new DeviceServiceIf();
    scheme = s->init("192.168.5.165", "admin", "zxw123456", 8000, 1);

//    checkHealth(catchErrorCallback_cgo);

    LoginDeviceDto dto = {0};
    login(&scheme, false, &dto);
    printf("%d\n", dto.lUserID);
    printf("%d\n", dto.device.struDeviceV30.byStartChan);

//    CapPictureDto capDto = {0};
//    capPicture(dto.lUserID, dto.device.struDeviceV30.byStartChan, &capDto);
//    printf("%d\n", capDto.lpSizeReturned);
//    for (int i = 0; i < 204800; i++) {
//        printf("%d", capDto.sJpegPicBuffer[i]);
//    }

    char url[512];
    char inbuf[10240];
    char outbuf[10240];
    char statbuf[1024];
    sprintf(url, "%s", "GET /ISAPI/Intelligent/FDLib/capabilities\r\n");
    stdXmlConfig(dto.lUserID, url, inbuf, outbuf, statbuf);
    printf("%s\n", statbuf);
    printf("%s\n", outbuf);

    // logout(dto.lUserID);
    delete s;
    return 0;
}
*/