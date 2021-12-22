#ifndef TYPES_H
#define TYPES_H

#ifndef _HC_NET_SDK_H_
#define  CALLBACK
#define  BOOL  int

typedef  unsigned int       DWORD;
typedef  unsigned short     WORD;
typedef  unsigned short     USHORT;
typedef  short              SHORT;
typedef  int                LONG;
typedef  unsigned char      BYTE;
typedef  unsigned int       UINT;
typedef  void*              LPVOID;
typedef  void*              HANDLE;
typedef  unsigned int*      LPDWORD;
typedef  unsigned long long UINT64;
typedef  signed long long   INT64;

typedef struct NET_DVR_STD_CONFIG {
    LPVOID    lpCondBuffer; // in
    DWORD     dwCondSize;   // in
    LPVOID    lpInBuffer;   // in
    DWORD     dwInSize;     // in
    LPVOID    lpOutBuffer;  // out
    DWORD     dwOutSize;    // in
    LPVOID    lpStatusBuffer;   // out
    DWORD     dwStatusSize; // in
    LPVOID    lpXmlBuffer;  // in/out
    DWORD     dwXmlSize;    // in/out
    BYTE      byDataType;   // in
    BYTE      byRes[23];    // in
}NET_DVR_STD_CONFIG;

//临时回调函数需要
typedef struct NET_DVR_ALARMER
{
    BYTE byUserIDValid;
    BYTE bySerialValid;
    BYTE byVersionValid;
    BYTE byDeviceNameValid;
    BYTE byMacAddrValid;
    BYTE byLinkPortValid;
    BYTE byDeviceIPValid;
    BYTE bySocketIPValid;
    LONG lUserID;
    BYTE sSerialNumber[48];
    DWORD dwDeviceVersion;
    char sDeviceName[32];
    BYTE byMacAddr[6];
    WORD wLinkPort;
    char sDeviceIP[128];
    char sSocketIP[128];
    BYTE byIpProtocol;
    BYTE byRes1[2];
    BYTE bJSONBroken;
    WORD wSocketPort;
    BYTE byRes2[6];
}NET_DVR_ALARMER;

typedef struct NET_DVR_IPADDR
{
    char    sIpV4[16];
    BYTE    byIPv6[128];
}NET_DVR_IPADDR;

typedef struct NET_DVR_SADPINFO
{
    NET_DVR_IPADDR  struIP;
    WORD            wPort;
    WORD            wFactoryType;
    char            chSoftwareVersion[48];
    char            chSerialNo[16];
    WORD            wEncCnt;
    BYTE            byMACAddr[6];
    NET_DVR_IPADDR  struSubDVRIPMask;
    NET_DVR_IPADDR  struGatewayIpAddr;
    NET_DVR_IPADDR    struDnsServer1IpAddr;
    NET_DVR_IPADDR    struDnsServer2IpAddr;
    BYTE            byDns;
    BYTE            byDhcp;
    BYTE            szGB28181DevID[32];
    BYTE            byActivated;
    BYTE            byDeviceModel[24];
    BYTE            byRes[101];
}NET_DVR_SADPINFO;

typedef struct NET_DVR_SADPINFO_LIST
{
    DWORD            dwSize;
    WORD             wSadpNum;
    BYTE             byRes[6];
    NET_DVR_SADPINFO struSadpInfo[256];
}NET_DVR_SADPINFO_LIST;

//设备通道信息
typedef struct NET_DVR_DEV_CHAN_INFO
{
    NET_DVR_IPADDR     struIP;            //DVR IP地址
    WORD     wDVRPort;                 //端口号
    BYTE     byChannel;                //通道号
    BYTE    byTransProtocol;        //传输协议类型0-TCP，1-UDP
    BYTE    byTransMode;            //传输码流模式 0－主码流 1－子码流
    BYTE    byFactoryType;            /*前端设备厂家类型,通过接口获取*/
    BYTE    byDeviceType; //设备类型(视频综合平台智能板使用)，1-解码器（此时根据视频综合平台能力集中byVcaSupportChanMode字段来决定是使用解码通道还是显示通道），2-编码器
    BYTE    byDispChan;//显示通道号,智能配置使用
    BYTE    bySubDispChan;//显示通道子通道号，智能配置时使用
    BYTE    byResolution;    //; 1-CIF 2-4CIF 3-720P 4-1080P 5-500w大屏控制器使用，大屏控制器会根据该参数分配解码资源
    BYTE    byRes[2];
    BYTE    byDomain[64];    //设备域名
    BYTE    sUserName[32];    //监控主机登陆帐号
    BYTE    sPassword[16];    //监控主机密码
}NET_DVR_DEV_CHAN_INFO;

//流媒体服务器基本配置
typedef struct NET_DVR_STREAM_MEDIA_SERVER_CFG
{
    BYTE    byValid;            /*是否可用*/
    BYTE    byRes1[3];
    NET_DVR_IPADDR  struDevIP;
    WORD    wDevPort;            /*流媒体服务器端口*/
    BYTE    byTransmitType;        /*传输协议类型 0-TCP，1-UDP*/
    BYTE    byRes2[69];
}NET_DVR_STREAM_MEDIA_SERVER_CFG;

//直接通过流媒体取流
typedef struct NET_DVR_PU_STREAM_CFG
{
    DWORD                                dwSize;
    NET_DVR_STREAM_MEDIA_SERVER_CFG    struStreamMediaSvrCfg;
    NET_DVR_DEV_CHAN_INFO                struDevChanInfo;
}NET_DVR_PU_STREAM_CFG;

typedef  struct NET_DVR_IPSERVER_STREAM
{
    BYTE    byEnable;   // 是否在线
    BYTE    byRes[3];               // 保留字节
    NET_DVR_IPADDR struIPServer;    //IPServer 地址
    WORD    wPort;                  //IPServer 端口
    WORD    wDvrNameLen;            // DVR 名称长度
    BYTE    byDVRName[32];    // DVR名称
    WORD    wDVRSerialLen;          // 序列号长度
    WORD    byRes1[2];              // 保留字节
    BYTE    byDVRSerialNumber[48];    // DVR序列号长度
    BYTE    byUserName[32];   // DVR 登陆用户名
    BYTE    byPassWord[16]; // DVR登陆密码
    BYTE    byChannel;              // DVR 通道
    BYTE    byRes2[11];             //  保留字节
}NET_DVR_IPSERVER_STREAM;

/* IP通道匹配参数 */
typedef struct NET_DVR_IPCHANINFO
{
    BYTE byEnable;                    /* 该通道是否在线 */
    BYTE byIPID;                    //IP设备ID低8位，当设备ID为0时表示通道不可用
    BYTE byChannel;                    /* 通道号 */
    BYTE byIPIDHigh;                // IP设备ID的高8位
    BYTE byTransProtocol;            //传输协议类型0-TCP/auto(具体有设备决定)，1-UDP 2-多播 3-仅TCP 4-auto
    BYTE byGetStream;         /* 是否对该通道取流，0-是，1-否*/
    BYTE byres[30];                    /* 保留 */
} NET_DVR_IPCHANINFO;

typedef struct NET_DVR_DDNS_STREAM_CFG
{
    BYTE   byEnable;   // 是否启用
    BYTE   byRes1[3];
    NET_DVR_IPADDR  struStreamServer;            //流媒体服务器地址
    WORD   wStreamServerPort;           //流媒体服务器端口
    BYTE   byStreamServerTransmitType;  //流媒体传输协议类型 0-TCP，1-UDP
    BYTE   byRes2;
    NET_DVR_IPADDR   struIPServer;          //IPSERVER地址
    WORD   wIPServerPort;        //IPserver端口号
    BYTE   byRes3[2];
    BYTE   sDVRName[32];   //DVR名称
    WORD   wDVRNameLen;            // DVR名称长度
    WORD   wDVRSerialLen;          // 序列号长度
    BYTE   sDVRSerialNumber[48];    // DVR序列号
    BYTE   sUserName[32];   // DVR 登陆用户名
    BYTE   sPassWord[16]; // DVR登陆密码
    WORD   wDVRPort;   //DVR端口号
    BYTE   byRes4[2];
    BYTE   byChannel;              // DVR 通道
    BYTE   byTransProtocol; //传输协议类型0-TCP，1-UDP
    BYTE   byTransMode; //传输码流模式 0－主码流 1－子码流
    BYTE   byFactoryType; //前端设备厂家类型,通过接口获取
}NET_DVR_DDNS_STREAM_CFG;

#define URL_LEN     240   //URL长度
typedef struct NET_DVR_PU_STREAM_URL
{
    BYTE    byEnable;
    BYTE    strURL[URL_LEN];
    BYTE    byTransPortocol ; // 传输协议类型 0-tcp  1-UDP
    WORD    wIPID;  //设备ID号，wIPID = iDevInfoIndex + iGroupNO*64 +1
    BYTE    byChannel;  //通道号
    BYTE    byRes[7];
}NET_DVR_PU_STREAM_URL;

typedef struct NET_DVR_HKDDNS_STREAM
{
    BYTE    byEnable;                 // 是否在线
    BYTE    byRes[3];               // 保留字节
    BYTE    byDDNSDomain[64];        // hiDDNS服务器
    WORD    wPort;                  // hiDDNS 端口
    WORD    wAliasLen;              // 别名长度
    BYTE    byAlias[32];         // 别名
    WORD    wDVRSerialLen;          // 序列号长度
    BYTE    byRes1[2];              // 保留字节
    BYTE    byDVRSerialNumber[48];    // DVR序列号
    BYTE    byUserName[32];   // DVR 登陆用户名
    BYTE    byPassWord[16]; // DVR登陆密码
    BYTE    byChannel;              // DVR通道
    BYTE    byRes2[11];             // 保留字
}NET_DVR_HKDDNS_STREAM;

typedef struct NET_DVR_IPCHANINFO_V40
{
    BYTE    byEnable;                /* 该通道是否在线 */
    BYTE    byRes1;
    WORD    wIPID;                  //IP设备ID
    DWORD     dwChannel;                //通道号
    BYTE    byTransProtocol;        //传输协议类型0-TCP，1-UDP，2- 多播，3-RTSP，0xff- auto(自动)
    BYTE    byTransMode;            //传输码流模式 0－主码流 1－子码流
    BYTE    byFactoryType;            /*前端设备厂家类型,通过接口获取*/
    BYTE    byRes;
    BYTE    strURL[URL_LEN/*240*/];   /*RTSP协议取流URL （仅RTSP协议时有效）*/
}NET_DVR_IPCHANINFO_V40;

typedef union NET_DVR_GET_STREAM_UNION
{
    NET_DVR_IPCHANINFO      struChanInfo;    /*IP通道信息*/
    NET_DVR_IPSERVER_STREAM struIPServerStream;  // IPServer去流
    NET_DVR_PU_STREAM_CFG   struPUStream;     //  通过前端设备获取流媒体去流
    NET_DVR_DDNS_STREAM_CFG struDDNSStream;     //通过IPServer和流媒体取流
    NET_DVR_PU_STREAM_URL   struStreamUrl;        //通过流媒体到url取流
    NET_DVR_HKDDNS_STREAM    struHkDDNSStream;   //通过hiDDNS去取流
    NET_DVR_IPCHANINFO_V40 struIPChan; //直接从设备取流（扩展）
}NET_DVR_GET_STREAM_UNION;

typedef struct NET_DVR_STREAM_MODE
{
    BYTE    byGetStreamType; //取流方式GET_STREAM_TYPE，0-直接从设备取流，1-从流媒体取流、2-通过IPServer获得ip地址后取流,3.通过IPServer找到设备，再通过流媒体去设备的流
    //4-通过流媒体由URL去取流,5-通过hkDDNS取流，6-直接从设备取流(扩展)，使用NET_DVR_IPCHANINFO_V40结构, 7-通过RTSP协议方式进行取流
    BYTE    byRes[3];        //保留字节
    NET_DVR_GET_STREAM_UNION uGetStream;    // 不同取流方式结构体
}NET_DVR_STREAM_MODE;

typedef struct NET_DVR_IPDEVINFO_V31
{
    BYTE byEnable;
    BYTE byProType;
    BYTE byEnableQuickAdd;

    BYTE byCameraType;
    BYTE sUserName[32];
    BYTE sPassword[16];
    BYTE byDomain[64];
    NET_DVR_IPADDR struIP;
    WORD wDVRPort;
    BYTE szDeviceID[32];
    BYTE byEnableTiming;
    BYTE byCertificateValidation;
}NET_DVR_IPDEVINFO_V31;

typedef struct NET_DVR_IPPARACFG_V40
{
    DWORD      dwSize;
    DWORD      dwGroupNum;
    DWORD      dwAChanNum;
    DWORD      dwDChanNum;
    DWORD      dwStartDChan;
    BYTE       byAnalogChanEnable[64];
    NET_DVR_IPDEVINFO_V31   struIPDevInfo[64];
    NET_DVR_STREAM_MODE     struStreamMode[64];
    BYTE            byRes2[20];
}NET_DVR_IPPARACFG_V40;

typedef struct NET_DVR_AREAINFOCFG
{
    WORD wNationalityID;
    WORD wProvinceID;
    WORD wCityID;
    WORD wCountyID;
    DWORD dwCode;
}NET_DVR_AREAINFOCFG;

typedef struct NET_VCA_DEV_INFO
{
    NET_DVR_IPADDR  struDevIP;
    WORD wPort;
    BYTE byChannel;
    BYTE byIvmsChannel;
}NET_VCA_DEV_INFO;

typedef struct NET_VCA_RECT
{
    float fX;
    float fY;
    float fWidth;
    float fHeight;
}NET_VCA_RECT;

typedef struct NET_VCA_HUMAN_ATTRIBUTE
{
    BYTE   bySex;
    BYTE   byCertificateType;
    BYTE   byBirthDate[10];
    BYTE   byName[32];
    NET_DVR_AREAINFOCFG struNativePlace;
    BYTE   byCertificateNumber[32];
    DWORD  dwPersonInfoExtendLen;
    BYTE   *pPersonInfoExtend;
    BYTE   byAgeGroup;
    BYTE   byRes2[3];
    BYTE*  pThermalData;
    BYTE   byRes3[4];
}NET_VCA_HUMAN_ATTRIBUTE;

typedef struct NET_VCA_BLOCKLIST_INFO
{
    DWORD  dwSize;
    DWORD  dwRegisterID;
    DWORD  dwGroupNo;
    BYTE   byType;
    BYTE   byLevel;
    BYTE   byRes1[2];
    NET_VCA_HUMAN_ATTRIBUTE struAttribute;
    BYTE   byRemark[32];
    DWORD  dwFDDescriptionLen;
    BYTE   *pFDDescriptionBuffer;
    DWORD  dwFCAdditionInfoLen;
    BYTE   *pFCAdditionInfoBuffer;
    DWORD  dwThermalDataLen;
}NET_VCA_BLOCKLIST_INFO;

typedef struct NET_VCA_FACESNAP_INFO_ALARM
{
    DWORD dwRelativeTime;
    DWORD dwAbsTime;
    DWORD dwSnapFacePicID;
    DWORD dwSnapFacePicLen;
    NET_VCA_DEV_INFO struDevInfo;
    BYTE  byFaceScore;
    BYTE bySex;
    BYTE byGlasses;
    BYTE byAge;
    BYTE byAgeDeviation;
    BYTE byAgeGroup;
    BYTE byFacePicQuality;
    BYTE  byRes;
    DWORD dwUIDLen;
    BYTE  *pUIDBuffer;
    float fStayDuration;
    BYTE  *pBuffer1;
}NET_VCA_FACESNAP_INFO_ALARM;

typedef struct NET_VCA_BLOCKLIST_INFO_ALARM
{
    NET_VCA_BLOCKLIST_INFO struBlockListInfo;
    DWORD dwBlockListPicLen;
    DWORD  dwFDIDLen;
    BYTE  *pFDID;
    DWORD  dwPIDLen;
    BYTE  *pPID;
    WORD  wThresholdValue;
    BYTE  byIsNoSaveFDPicture;
    BYTE  byRealTimeContrast;
    BYTE  *pBuffer1;
}NET_VCA_BLOCKLIST_INFO_ALARM;

typedef struct NET_VCA_FACESNAP_MATCH_ALARM
{
    DWORD dwSize;
    float fSimilarity;
    NET_VCA_FACESNAP_INFO_ALARM  struSnapInfo;
    NET_VCA_BLOCKLIST_INFO_ALARM struBlockListInfo;
    char         sStorageIP[16];
    WORD            wStoragePort;
    BYTE  byMatchPicNum;
    BYTE  byPicTransType;
    DWORD dwSnapPicLen;
    BYTE  *pSnapPicBuffer;
    NET_VCA_RECT  struRegion;
    DWORD dwModelDataLen;
    BYTE  *pModelDataBuffer;
    BYTE  byModelingStatus;
    BYTE  byLivenessDetectionStatus;
    char  cTimeDifferenceH;
    char  cTimeDifferenceM;
    BYTE  byMask;
    BYTE  bySmile;
    BYTE  byContrastStatus;
    BYTE  byBrokenNetHttp;
}NET_VCA_FACESNAP_MATCH_ALARM;
#endif

typedef BOOL(CALLBACK *MessageCallback)(LONG lCommand, NET_DVR_ALARMER *pAlarmer, char *pAlarmInfo, DWORD dwBufLen, void *pUser);
typedef void(CALLBACK *ExceptionCallBack)(DWORD dwType, LONG lUserID, LONG lHandle, void *pUser);
typedef void(CALLBACK *StdDataCallBack)(LONG lRealHandle, DWORD dwDataType, BYTE *pBuffer, DWORD dwBufSize, void *pUser);

typedef struct Scheme
{
    DWORD port;
    DWORD channel;
    char *address;
    char *username;
    char *password;
}Scheme;

typedef struct LocalIpDto
{
    // char (*strIp)[16];
    char strIp[16][16];
    DWORD pValidNum;
    BOOL pEnableBind;
}LocalIpDto;

typedef struct ResolveSvrParam
{
    char *sServerIP;
    WORD wServerPort;
    BYTE *sDVRName;
    WORD wDVRNameLen;
    BYTE *sDVRSerialNumber;
    WORD wDVRSerialLen;
}ResolveSvrParam;

typedef struct ResolveSvrDto
{
    char *sGetIP;
    DWORD dwPort;
}ResolveSvrDto;

typedef struct RealPlayParam
{
    LONG lUserID;
    LONG lChannel;
    DWORD dwStreamType;
    DWORD dwLinkMode;
    DWORD bBlocked;
    DWORD bPassbackRecord;
}RealPlayParam;

typedef struct RealPlayInfo
{
    char szIP[16];
    LONG lUserID;
    DWORD lChannel;
}RealPlayInfo;

typedef struct DeviceInfoV30
{
    BYTE sSerialNumber[48];
    BYTE byDVRType;
    WORD wDevType;//设备型号扩展
    BYTE byChanNum;
    BYTE byStartChan;
}DeviceInfoV30;

typedef struct DeviceInfoV40
{
    DeviceInfoV30 struDeviceV30;
    BYTE byRetryLoginTime;
    BYTE byPasswordLevel;
}DeviceInfoV40;

typedef struct LoginDeviceDto
{
    LONG lUserID;
    DeviceInfoV40 device;
}LoginDeviceDto;

typedef struct RealCapPictureDto
{
    char  pPicBuf[204800];
    DWORD lpSizeReturned;
}RealCapPictureDto;

typedef struct CapPictureDto
{
    WORD  wPicSize;
    WORD  wPicQuality;
    char  sJpegPicBuffer[204800];
    DWORD lpSizeReturned;
}CapPictureDto;

typedef struct HealthParam
{
    LONG connectTime;
    LONG recvTimeOut;
    LONG reconnect;
    char *logToFile;
    DWORD logLevel;
}HealthParam;

typedef struct FaceDetection
{
    DWORD     dwSize;
    DWORD     dwRelativeTime;
    DWORD     dwAbsTime;
    DWORD     dwBackgroundPicLen;
    NET_VCA_RECT struFacePic[30];
    BYTE*     pBackgroundPicpBuffer;
}FaceDetection;

typedef struct GoAlarmer
{
    LONG lUserID;
    BYTE sSerialNumber[48];
    DWORD dwDeviceVersion;
    char sDeviceName[128];
}GoAlarmer;

#endif
