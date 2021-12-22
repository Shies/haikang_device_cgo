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
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
	"unsafe"

	"go_device_haikang/app/config"
	"go_device_haikang/app/gdb"
	"go_device_haikang/utils/request"
	"go_device_haikang/utils/util"
)

const (
	UcharSize16   = int32(16)
	UcharSize32   = int32(32)
	UcharSize48   = int32(48)
	ScharSize16   = int32(16)
	ScharSize32   = int32(32)
	ScharSize200K = int64(204800)
)

type Service struct {
	context.Context
	Device *config.Device
	Scheme *Scheme
	close  chan bool
}

func NewService(ctx context.Context) *Service {
	s := &Service{
		Device: config.DeviceConfig,
		close:  make(chan bool, 1),
	}

	nPort, _ := strconv.Atoi(s.Device.Port)
	nChan, _ := strconv.Atoi(s.Device.Channel)
	s.Scheme = &Scheme{
		Address:  s.Device.Address,
		Username: s.Device.Username,
		Password: s.Device.Password,
		Port:     nPort,
		Channel:  nChan,
	}

	s.loadConfig()
	lDto := s.active(s.Scheme)
	s.loadDevice(int(lDto.lUserID))
	go s.activeDevice(int(lDto.lUserID))
	go s.listenAlarm(int(lDto.lUserID))
	go s.loadBadData()
	return s
}

// FindDevice查找设备
func (s *Service) activeDevice(lUserID int) bool {
	CDto := C.struct_NET_DVR_IPPARACFG_V40{}
	C.getDVRConfig(C.int(lUserID), C.int(0), (*C.struct_NET_DVR_IPPARACFG_V40)(unsafe.Pointer(&CDto)))

	// 监听设备
	health, _ := strconv.Atoi(s.Device.CheckHealth)
	for {
		device := s.resolveDeviceIp(CDto)
		if device != nil {
			for _, scheme := range device {
				scheme.IsSave = 0
				scheme.Active = int(checkDevice(scheme.Address, scheme.Port))
			}
			request.HttpPost(s.Device.DevPath+RemoteInfoUrl, util.MarshalJson(device), nil)
		}
		time.Sleep(time.Duration(health) * time.Second)
	}
}

func checkDevice(ip string, port int) int64 {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), time.Second)
	if err != nil {
		log.Printf("%v", err)
		return 0
	}
	defer conn.Close()
	return 1
}

// 断网重传，本地sqlite数据重新发送给接收方
func (s *Service) loadBadData() {
	health, _ := strconv.Atoi(s.Device.CheckHealth)
	for {
		var logs []gdb.DeviceLog
		gdb.DB.Table(gdb.TableNameDeviceLog).Limit(10).Scan(&logs)
		for _, log := range logs {
			if log.LogType > 0 {
				gdb.DB.Table(gdb.TableNameDeviceLog).Where("id=?", log.Id).Delete(nil)
			} else {
				var url string
				var record gdb.DeviceRecord
				gdb.DB.Table(gdb.TableNameDevice).Where("id=?", log.RecordId).Scan(&record)
				if record.PictureUrl != "" {
					url = s.Device.DevPath+RemoteRecordUrl
				} else {
					url = s.Device.DevPath+RemoteLogUrl
				}
				_, err := request.HttpPost(url, log.LogContent, nil)
				if err == nil {
					gdb.DB.Table(gdb.TableNameDeviceLog).Where("id=?", log.Id).Delete(nil)
					gdb.DB.Table(gdb.TableNameDevice).Where("id=?", log.RecordId).Delete(nil)
				}
			}
		}
		time.Sleep(time.Duration(health * 4) * time.Second)
	}
}

// 保存设备全局配置
func (s *Service) loadConfig() {
	C.SDKInit()
	param := LocalHealth(s.Device)
	//defer C.free(unsafe.Pointer(param.logToFile))

	C.checkHealth(&param, (C.ExceptionCallBack)(unsafe.Pointer(C.catchErrorCallback_cgo)))
	if s.Device != nil {
		request.HttpPost(s.Device.DevPath+RemoteConfigUrl, util.MarshalJson(s.Device), nil)
	}
	return
}

// 断线重连
func (s *Service) active(scheme *Scheme) C.struct_LoginDeviceDto {
	cscheme := GetCScheme(scheme)
	//defer C.free(unsafe.Pointer(cscheme.address))
	//defer C.free(unsafe.Pointer(cscheme.username))
	//defer C.free(unsafe.Pointer(cscheme.password))

	lDto := C.struct_LoginDeviceDto{}
	C.login(cscheme, C.int(0), (*C.struct_LoginDeviceDto)(unsafe.Pointer(&lDto)))
	return lDto
}

// 保存设备进行管理
func (s *Service) loadDevice(lUserID int) {
	CDto := C.struct_NET_DVR_IPPARACFG_V40{}
	C.getDVRConfig(C.int(lUserID), C.int(0), (*C.struct_NET_DVR_IPPARACFG_V40)(unsafe.Pointer(&CDto)))
	// 保存设备
	deivce := s.resolveDeviceIp(CDto)
	if deivce != nil {
		for _, scheme := range deivce {
			lDto := s.active(&Scheme{
				Address:  scheme.Address,
				Username: s.Device.DevUser,
				Password: s.Device.DevPass,
				Port:     scheme.Port,
				Channel:  s.Scheme.Channel,
			})
			scheme.Channel = s.Scheme.Channel
			scheme.DeviceName = strconv.Itoa(int(C.ushort(lDto.device.struDeviceV30.wDevType)))
			scheme.SerialNo = string(GetUchar48(lDto.device.struDeviceV30.sSerialNumber))
			scheme.IsSave = 1
			C.logout(C.int(lDto.lUserID))
		}
		log.Println(deivce)
		//C.logout(C.int(lUserID))
		//C.clean()
		request.HttpPost(s.Device.DevPath+RemoteInfoUrl, util.MarshalJson(deivce), nil)
	}
}

// 检测设备是否在线
/*
func (s *Service) activeDevice(lUserID int) {
	for {
		var status C.int
		C.active(C.int(lUserID), &status)
		if int(status) == 0 {
			device := ([]*Device)(nil)
			CDto := C.struct_NET_DVR_IPPARACFG_V40{}
			C.getDVRConfig(C.int(lUserID), C.int(0), (*C.struct_NET_DVR_IPPARACFG_V40)(unsafe.Pointer(&CDto)))
			for i := 0; i < int(CDto.dwDChanNum); i++ {
				stream := CDto.struStreamMode[i].uGetStream
				ipChan := *(*C.struct_NET_DVR_IPCHANINFO)(unsafe.Pointer(&stream))
				if uint(ipChan.byIPID) > 0 {
					// 下标0对应设备IP ID为1
					devInfo := CDto.struIPDevInfo[int(ipChan.byIPID)-1]
					device = append(device, &Device{
						IsSave:   0,
						Active:   int(uint8(C.uchar(ipChan.byEnable))),
						Address:  string(GetSchar16(devInfo.struIP.sIpV4)),
						Port:     int(C.ushort(devInfo.wDVRPort)),
						Username: string(GetUchar32(devInfo.sUserName)),
						Password: string(GetUchar16(devInfo.sPassword)),
						Channel:  int(uint8(C.uchar(ipChan.byChannel))),
					})
				}
			}
			C.logout(C.int(lUserID))
			request.HttpPost(s.Device.DevPath+RemoteInfoUrl, util.MarshalJson(device), nil)
			health, _ := strconv.Atoi(s.Device.CheckHealth)
			time.Sleep(time.Duration(health) * time.Second)
		}
		lDto := s.active(s.Scheme)
		lUserID = int(lDto.lUserID)
	}
}
*/

// 设备告警监听
func (s *Service) listenAlarm(lUserID int) {
	for {
		lDto := s.active(s.Scheme)
		lUserID = int(lDto.lUserID)

		var status C.int
		C.active(C.int(lUserID), &status)
		if int(status) == 0 {
			var lHandle int
			C.setupAlarmChan(C.int(lUserID), (*C.int)(unsafe.Pointer(&lHandle)), (C.MessageCallback)(unsafe.Pointer(C.alarmMsgCallback_cgo)))
			log.Printf("setupAlarmChan start...")
			<-s.close
			s.closeAlarm(lUserID, lHandle)
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
}

// 关闭告警
func (s *Service) closeAlarm(lUserID int, lHandle int) bool {
	C.closeAlarmChan(C.int(lHandle))
	C.logout(C.int(lUserID))
	C.clean()
	return false
}

// 发现设备
func (s *Service) resolveDeviceIp(CDto C.struct_NET_DVR_IPPARACFG_V40) []*Device {
	device := ([]*Device)(nil)
	for i := 0; i < int(CDto.dwDChanNum); i++ {
		ipDevInfo := CDto.struIPDevInfo[i]
		if uint(ipDevInfo.byEnable) != 0 {
			device = append(device, &Device{
				IsSave:   0,
				Channel:  1,
				Active:   int(uint8(C.uchar(ipDevInfo.byEnable))),
				Address:  string(GetSchar16(ipDevInfo.struIP.sIpV4)),
				Port:     int(C.ushort(ipDevInfo.wDVRPort)),
				Username: string(GetUchar32(ipDevInfo.sUserName)),
				Password: string(GetUchar16(ipDevInfo.sPassword)),
			})
		}
	}

	return device
}

// 转换CScheme唤醒海康摄像头
func GetCScheme(scheme *Scheme) *C.struct_Scheme {
	address := C.CString(scheme.Address)
	username := C.CString(scheme.Username)
	password := C.CString(scheme.Password)

	var CScheme C.struct_Scheme
	CScheme = C.struct_Scheme{
		channel:  C.uint(scheme.Channel),
		port:     C.uint(scheme.Port),
		address:  address,
		username: username,
		password: password,
	}
	return (*C.struct_Scheme)(unsafe.Pointer(&CScheme))
}

// 发现设备ip需要的参数
func GetIpResolveParam(param *ResolveSvrParam) C.struct_ResolveSvrParam {
	serverIp := C.CString(param.ServerIP)
	dvrName := C.CString(param.DVRName)
	dvrSerial := C.CString(param.DVRSerialNumber)
	//defer C.free(unsafe.Pointer(dvrName))
	//defer C.free(unsafe.Pointer(serverIp))
	//defer C.free(unsafe.Pointer(dvrSerial))

	var CParam C.struct_ResolveSvrParam
	CParam = C.struct_ResolveSvrParam{
		sServerIP:        serverIp,
		wServerPort:      C.ushort(param.ServerPort),
		sDVRName:         (*C.uchar)(unsafe.Pointer(dvrName)),
		sDVRSerialNumber: (*C.uchar)(unsafe.Pointer(dvrSerial)),
	}

	return CParam
}

// 获取本地服务器多网卡IP地址
func GetLocalIp(strIp [ScharSize16][ScharSize16]C.char) []string {
	var welcome []string
	for i := 0; i < len(strIp); i++ {
		var ipch []byte
		for j := 0; j < len(strIp[i]); j++ {
			if strIp[i][j] != 0 {
				ipch = append(ipch, byte(strIp[i][j]))
			}
		}
		welcome = append(welcome, string(ipch))
	}

	return welcome
}

// 获取设备的名称和软件版本
func GetUchar32(str [UcharSize32]C.uchar) []byte {
	var soft []byte
	for j := 0; j < len(str); j++ {
		if byte(str[j]) != 0 {
			soft = append(soft, byte(str[j]))
		}
	}

	return soft
}

func GetUchar48(str [UcharSize48]C.uchar) []byte {
	var soft []byte
	for j := 0; j < len(str); j++ {
		if byte(str[j]) != 0 {
			soft = append(soft, byte(str[j]))
		}
	}

	return soft
}

func GetUchar16(str [UcharSize16]C.uchar) []byte {
	var soft []byte
	for j := 0; j < len(str); j++ {
		if byte(str[j]) != 0 {
			soft = append(soft, byte(str[j]))
		}
	}

	return soft
}

func GetSchar16(str [ScharSize16]C.char) []byte {
	var out []byte
	for i := 0; i < len(str); i++ {
		if str[i] != 0 {
			out = append(out, byte(str[i]))
		}
	}

	return out
}

func GetSchar32(str [ScharSize32]C.char) []byte {
	var out []byte
	for i := 0; i < len(str); i++ {
		if str[i] != 0 {
			out = append(out, byte(str[i]))
		}
	}

	return out
}

// 抓拍图片流解析
func GetStream(buf [ScharSize200K]C.char, size int) []byte {
	var out []byte
	for i := 0; i < size; i++ {
		if buf[i] != 0 {
			out = append(out, byte(buf[i]))
		}
	}

	return out
}

// SDK参数设置
func LocalHealth(health *config.Device) C.struct_HealthParam {
	ct, _ := strconv.Atoi(health.ConnectTime)
	ro, _ := strconv.Atoi(health.RecvTimeOut)
	rc, _ := strconv.Atoi(health.Reconnect)
	ll, _ := strconv.Atoi(health.LogLevel)

	var param C.struct_HealthParam
	param = C.struct_HealthParam{
		connectTime: C.int(ct),
		recvTimeOut: C.int(ro),
		reconnect:   C.int(rc),
		logToFile:   C.CString(health.LogToFile),
		logLevel:    C.uint(ll),
	}

	return param
}

// UploadImg 上传图片
func UploadImg(buf []byte, fileUploadPath string) (UploadImageResult, error) {
	result := UploadImageResult{}
	bodyBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuf)
	// 创建一个文件句柄
	ioWriter, err := bodyWriter.CreateFormFile("file", "cap-"+util.RandomUUID()+".jpeg")
	if err != nil {
		return result, err
	}
	// 将图片流写入文件句柄
	ioWriter.Write(buf)
	boundary := bodyWriter.Boundary()
	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	requestReader := io.MultiReader(bodyBuf, closeBuf)
	req, err := http.NewRequest("POST", fileUploadPath, requestReader)
	if err != nil {
		return result, err
	}
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = int64(bodyBuf.Len()) + int64(closeBuf.Len())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	util.UnmarshalJson(body, &result)
	return result, nil
}

func getLocal(url string) ([]byte, error) {
	fp, err := os.OpenFile(url, os.O_CREATE|os.O_APPEND, 6) // 读写方式打开
	if err != nil {
		// 如果有错误返回错误内容
		return nil, err
	}

	defer fp.Close()
	bytes, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	return bytes, err
}
