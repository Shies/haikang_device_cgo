package haikang

/*
#cgo CXXFLAGS: -I./
#cgo LDFLAGS: -L. -Wl,-rpath=./:./lib:./lib/HCNetSDKCom -lhcnetsdk -lstdc++
#cgo LDFLAGS: -L. -ldevice_handler -lstdc++
#include <stdio.h>
#include <stdlib.h>
#include "DeviceHandler.h"
*/
import "C"
import (
	"context"
	"unsafe"

	"go_device_haikang/app/config"
	"go_device_haikang/utils/response"

	"github.com/gin-gonic/gin"
)

var UserApi = NewUserApi(context.Background())
type userApi struct {
	Device  *config.Device
	Service *Service
}

func NewUserApi(ctx context.Context) *userApi {
	a := &userApi{
		Service: NewService(ctx),
		Device: config.DeviceConfig,
	}

	return a
}

func RegisterRouterNoAuth(g *gin.RouterGroup) {
	g.POST("/user/get/extend/config", getExtendConfig)
	g.POST("/user/resolve/device/ip", resolveIp)
	g.POST("/user/setup/alarm/chan", setupAlarmChan)
	g.POST("/user/close/alarm/chan", closeAlarmChan)
	g.POST("/user/ping", ping)
	g.POST("/user/login", login)
	g.POST("/user/logout", logout)
	g.POST("/user/active", active)
	g.POST("/user/local/ip", localIp)
	g.POST("/user/resolve/server/ip", iPByResolveSvr)
	g.POST("/user/get/dvr/config", getDVRConfig)
	g.POST("/user/real/play", realPlay)
	g.POST("/user/real/stop/play", realStopPlay)
	g.POST("/user/real/cap/picture", realCapPicture)
	g.POST("/user/cap/picture", capPicture)
	return
}

// @summary 获取人脸比对库图片数据附加信息
// @tags    海康对接-获取人脸比对库图片数据附加信息
// @produce json
// @Param param body PingParam true "附件信息参数"
// @router  /haikang/user/get/extend/config [POST]
// @success 200 {object} response.Response "获取成功"
func getExtendConfig(c *gin.Context) {
	param := LocalHealth(UserApi.Service.Device)
	C.checkHealth(&param, (C.ExceptionCallBack)(unsafe.Pointer(C.catchErrorCallback_cgo)))

	lDto := C.struct_LoginDeviceDto{}
	scheme := GetCScheme(UserApi.Service.Scheme)
	C.login(scheme, C.int(0), (*C.struct_LoginDeviceDto)(unsafe.Pointer(&lDto)))
	lUserID := C.int(lDto.lUserID)

	url := C.CString("GET /ISAPI/Intelligent/FDLib/0/picture/32")
	defer C.free(unsafe.Pointer(url))
	inch := C.CString("")
	defer C.free(unsafe.Pointer(inch))
	outch := C.CString("")
	defer C.free(unsafe.Pointer(outch))
	statch := C.CString("")
	defer C.free(unsafe.Pointer(statch))

	C.stdXmlConfig(lUserID, url, inch, outch, statch)
	checkHealth := &PingParam{
		Welcome: C.GoString(outch),
	}
	response.Success(c, checkHealth)
	return
}

// @summary 发现设备
// @tags    海康对接-发现设备
// @produce json
// @Param param body PingParam true "发现设备参数"
// @router  /haikang/user/resolve/deivce/ip [POST]
// @success 200 {object} response.Response "获取成功"
func resolveIp(c *gin.Context) {
	param := LocalHealth(UserApi.Service.Device)
	C.checkHealth(&param, (C.ExceptionCallBack)(unsafe.Pointer(C.catchErrorCallback_cgo)))
	checkHealth := &PingParam{
		Welcome: "ok",
	}
	response.Success(c, checkHealth)
	return
}

// @summary 心跳检测
// @tags    海康对接-心跳检测
// @produce json
// @Param param body PingParam true "心跳检测参数"
// @router  /haikang/user/ping [POST]
// @success 200 {object} response.Response "获取成功"
func ping(c *gin.Context) {
	param := LocalHealth(UserApi.Service.Device)
	C.checkHealth(&param, (C.ExceptionCallBack)(unsafe.Pointer(C.catchErrorCallback_cgo)))
	checkHealth := &PingParam{
		Welcome: "ok",
	}
	response.Success(c, checkHealth)
	return
}

// @summary 布防报警
// @tags    海康对接-布防报警
// @produce json
// @Param param body PingParam true "布防报警参数"
// @router  /haikang/user/setup/alarm/chan [POST]
// @success 200 {object} response.Response "获取成功"
func setupAlarmChan(c *gin.Context) {
	param := LocalHealth(UserApi.Service.Device)
	C.checkHealth(&param, (C.ExceptionCallBack)(unsafe.Pointer(C.catchErrorCallback_cgo)))

	lDto := C.struct_LoginDeviceDto{}
	scheme := GetCScheme(UserApi.Service.Scheme)
	C.login(scheme, C.int(0), (*C.struct_LoginDeviceDto)(unsafe.Pointer(&lDto)))
	lUserID := C.int(lDto.lUserID)

	var lHandle int
	C.setupAlarmChan(lUserID, (*C.int)(unsafe.Pointer(&lHandle)), (C.MessageCallback)(unsafe.Pointer(C.alarmMsgCallback_cgo)))
	C.logout(lUserID)
	checkHealth := &PingParam{
		Welcome: lHandle,
	}
	response.Success(c, checkHealth)
	return
}

// @summary 关闭报警
// @tags    海康对接-关闭报警
// @produce json
// @Param param body PingParam true "关闭报警参数"
// @router  /haikang/user/close/alarm/chan [POST]
// @success 200 {object} response.Response "获取成功"
func closeAlarmChan(c *gin.Context) {
	param := LocalHealth(UserApi.Service.Device)
	C.checkHealth(&param, (C.ExceptionCallBack)(unsafe.Pointer(C.catchErrorCallback_cgo)))

	lHandle := 0
	C.closeAlarmChan(C.int(lHandle))
	checkHealth := &PingParam{
		Welcome: lHandle,
	}
	response.Success(c, checkHealth)
	return
}

// @summary 设备激活
// @tags    海康对接-设备激活
// @produce json
// @Param param body PingParam true "设备激活参数"
// @router  /haikang/user/active [POST]
// @success 200 {object} response.Response "获取成功"
func active(c *gin.Context) {
	param := LocalHealth(UserApi.Service.Device)
	C.checkHealth(&param, (C.ExceptionCallBack)(unsafe.Pointer(C.catchErrorCallback_cgo)))

	lDto := C.struct_LoginDeviceDto{}
	scheme := GetCScheme(UserApi.Service.Scheme)
	C.login(scheme, C.int(0), (*C.struct_LoginDeviceDto)(unsafe.Pointer(&lDto)))
	lUserID := C.int(lDto.lUserID)

	dDto := C.struct_NET_DVR_SADPINFO_LIST{}
	C.getSadpInfoList(lUserID, (*C.struct_NET_DVR_SADPINFO_LIST)(unsafe.Pointer(&dDto)))
	checkHealth := &PingParam{
		Welcome: C.ushort(dDto.wSadpNum),
	}
	response.Success(c, checkHealth)
	return
}

// @summary 设备登录
// @tags    海康对接-设备登录
// @produce json
// @Param param body PingParam true "设备登录参数"
// @router  /haikang/user/login [POST]
// @success 200 {object} response.Response "获取成功"
func login(c *gin.Context) {
	CDto := C.struct_LoginDeviceDto{}

	useAsync := 0
	scheme := GetCScheme(UserApi.Service.Scheme)
	C.login(scheme, C.int(useAsync), (*C.struct_LoginDeviceDto)(unsafe.Pointer(&CDto)))
	checkHealth := &PingParam{
		Welcome: C.int(CDto.lUserID),
	}
	response.Success(c, checkHealth)
	return
}

// @summary 设备登出
// @tags    海康对接-设备登出
// @produce json
// @Param param body PingParam true "设备登出参数"
// @router  /haikang/user/logout [POST]
// @success 200 {object} response.Response "获取成功"
func logout(c *gin.Context) {
	lUserID := 0
	C.logout(C.int(lUserID))
	checkHealth := &PingParam{
		Welcome: C.int(lUserID),
	}
	response.Success(c, checkHealth)
	return
}

// @summary 本地多网卡IP获取
// @tags    海康对接-本地多网卡IP获取
// @produce json
// @Param param body PingParam true "本地多网卡IP获取参数"
// @router  /haikang/user/local/ip [POST]
// @success 200 {object} response.Response "获取成功"
func localIp(c *gin.Context) {
	CDto := C.struct_LocalIpDto{}

	C.localIp((*C.struct_LocalIpDto)(unsafe.Pointer(&CDto)))
	checkHealth := &PingParam{
		Welcome: GetLocalIp(CDto.strIp),
	}
	response.Success(c, checkHealth)
	return
}

// @summary 自动发现设备IP和端口
// @tags    海康对接-自动发现设备IP和端口
// @produce json
// @Param param body PingParam true "自动发现设备IP和端口"
// @router  /haikang/user/resolve/server/ip [POST]
// @success 200 {object} response.Response "获取成功"
func iPByResolveSvr(c *gin.Context) {
	CDto := C.struct_ResolveSvrDto{}
	CParam := GetIpResolveParam(nil)

	C.iPByResolveSvr((*C.struct_ResolveSvrParam)(unsafe.Pointer(&CParam)), (*C.struct_ResolveSvrDto)(unsafe.Pointer(&CDto)))
	checkHealth := &PingParam{
		Welcome: C.GoString(CDto.sGetIP),
	}
	response.Success(c, checkHealth)
	return
}

// @summary errno-23-设备不支持 系统参数配置
// @tags    海康对接-系统参数配置
// @produce json
// @Param param body PingParam true "系统参数配置"
// @router  /haikang/user/ip/config/param [POST]
// @success 200 {object} response.Response "获取成功"
func getDVRConfig(c *gin.Context) {
	lUserID := 0
	lChannel := 1
	CDto := C.struct_NET_DVR_IPPARACFG_V40{}

	C.getDVRConfig(C.int(lUserID), C.int(lChannel), (*C.struct_NET_DVR_IPPARACFG_V40)(unsafe.Pointer(&CDto)))
	checkHealth := &PingParam{
		Welcome: C.int(CDto.dwDChanNum),
	}
	response.Success(c, checkHealth)
	return
}

// @summary 实时预览
// @tags    海康对接-实时预览
// @produce json
// @Param param body PingParam true "实时预览参数"
// @router  /haikang/user/real/play [POST]
// @success 200 {object} response.Response "获取成功"
func realPlay(c *gin.Context) {
	CParam := C.struct_RealPlayParam{
		lUserID:  C.int(0),
		lChannel: C.int(1),
	}
	C.realPlay((*C.struct_RealPlayParam)(unsafe.Pointer(&CParam)), (C.StdDataCallBack)(unsafe.Pointer(C.realDataCallBack_cgo)))
	checkHealth := &PingParam{
		Welcome: "ok",
	}
	response.Success(c, checkHealth)
	return
}

// @summary 停止预览
// @tags    海康对接-停止预览
// @produce json
// @Param param body PingParam true "停止预览参数"
// @router  /haikang/user/real/stop/play [POST]
// @success 200 {object} response.Response "获取成功"
func realStopPlay(c *gin.Context) {
	lRealPlayHandle := 0
	C.realStopPlay(C.int(lRealPlayHandle))
	checkHealth := &PingParam{
		Welcome: C.int(lRealPlayHandle),
	}
	response.Success(c, checkHealth)
	return
}

// @summary 预览抓拍
// @tags    海康对接-预览抓拍
// @produce json
// @Param param body PingParam true "预览抓拍参数"
// @router  /haikang/user/real/cap/picture [POST]
// @success 200 {object} response.Response "获取成功"
func realCapPicture(c *gin.Context) {
	lRealPlayHandle := 1
	CDto := C.struct_RealCapPictureDto{}

	C.realCapPicture(C.int(lRealPlayHandle), (*C.struct_RealCapPictureDto)(unsafe.Pointer(&CDto)))
	length := int(C.uint(CDto.lpSizeReturned))
	buffer := GetStream(CDto.pPicBuf, length)
	//ioutil.WriteFile("./record/preview-ssss.jpeg", buffer, 0644)

	checkHealth := &PingParam{
		Welcome: string(buffer),
	}
	response.Success(c, checkHealth)
	return
}

// @summary 设备抓拍
// @tags    海康对接-设备抓拍
// @produce json
// @Param param body PingParam true "设备抓拍参数"
// @router  /haikang/user/cap/picture [POST]
// @success 200 {object} response.Response "获取成功"
func capPicture(c *gin.Context) {
	param := LocalHealth(UserApi.Service.Device)
	C.checkHealth(&param, (C.ExceptionCallBack)(unsafe.Pointer(C.catchErrorCallback_cgo)))

	lDto := C.struct_LoginDeviceDto{}
	// 这里是登录设备scheme配置，不是nvr配置
	scheme := GetCScheme(UserApi.Service.Scheme)
	C.login(scheme, C.int(0), (*C.struct_LoginDeviceDto)(unsafe.Pointer(&lDto)))
	lUserID := C.int(lDto.lUserID)
	lChannel := C.int(int(byte(lDto.device.struDeviceV30.byStartChan)))

	CDto := C.struct_CapPictureDto{}
	C.capPicture(lUserID, lChannel, (*C.struct_CapPictureDto)(unsafe.Pointer(&CDto)))
	length := int(C.uint(CDto.lpSizeReturned))
	buffer := GetStream(CDto.sJpegPicBuffer, length)

	result, _ := UploadImg(buffer, UserApi.Device.OssPath)
	C.logout(lUserID)
	checkHealth := &PingParam{
		Welcome: result,
	}
	response.Success(c, checkHealth)
	return
}
