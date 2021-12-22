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
	"crypto/md5"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unsafe"
)

const (
	RemoteInfoUrl   = "/v21/haikang/open/device/scheme"
	RemoteConfigUrl = "/v21/haikang/open/device/health"
	RemoteLogUrl    = "/v21/haikang/open/device/log"
	RemoteRecordUrl = "/v21/haikang/open/device/record"
)

//export GoRealDataCallBack
func GoRealDataCallBack(lRealHandle C.int, dwDataType C.uint, pBuffer *C.uchar, dwBufSize C.uint, pUser unsafe.Pointer) {
	fmt.Println("GoRealDataCallBack")
	fmt.Println(C.uint(dwBufSize))
	fmt.Println((*C.char)(unsafe.Pointer(pBuffer)))
	fmt.Println(C.int(lRealHandle))
	fmt.Println(C.uint(dwDataType))
	fmt.Println((*C.struct_RealPlayInfo)(pUser))
	return
}

//export GoCatchErrorCallback
func GoCatchErrorCallback(dwType C.uint, lUserID C.int, lHandle C.int, pUser unsafe.Pointer) {
	switch dwType {
	// EXCEPTION_RECONNECT
	case 0x8005:
		fmt.Println("GoCatchErrorCallback")
		fmt.Println(C.uint(dwType))
		fmt.Println(C.int(lHandle))
		fmt.Println(C.int(lUserID))
		fmt.Println((*C.uint)(pUser), time.Now().Unix())
	}
	return
}

//export GoAlarmMsgCallback
func GoAlarmMsgCallback(lCommand C.int, pAlarmer unsafe.Pointer, pAlarmInfo unsafe.Pointer, dwBufLen C.uint, pUser unsafe.Pointer) bool {
	alarm := (*C.struct_NET_DVR_ALARMER)(pAlarmer)
	// fmt.Println("报警设备用户ID", alarm.lUserID, "\n")
	httpPath := UserApi.Device.DevPath
	switch lCommand {
	// COMM_SNAP_MATCH_ALARM
	case 0x2902:
		match := (*C.struct_NET_VCA_FACESNAP_MATCH_ALARM)(pAlarmInfo)
		if match == nil {
			return false
		}
		uname := string(GetUchar32(match.struBlockListInfo.struBlockListInfo.struAttribute.byName))
		err := StrToUtf8(&uname)
		if err != nil {
			log.Printf("%v", err)
		}
		var v XmlSchema
		pExtendLen := int(match.struBlockListInfo.struBlockListInfo.struAttribute.dwPersonInfoExtendLen)
		pExtendBuf := (*C.char)(unsafe.Pointer(match.struBlockListInfo.struBlockListInfo.struAttribute.pPersonInfoExtend))
		if pExtendLen > 0 && pExtendBuf != nil {
			schema := "<XmlSchema>" + C.GoString(pExtendBuf) + "</XmlSchema>"
			err := xml.Unmarshal([]byte(schema), &v)
			if err != nil {
				log.Printf("%v", err)
			}
		}
		record := DeviceRecord{
			Address:     string(GetSchar16(match.struSnapInfo.struDevInfo.struDevIP.sIpV4)),
			Port:        int(match.struSnapInfo.struDevInfo.wPort),
			Channel:     int(match.struSnapInfo.struDevInfo.byChannel),
			DVRName:     string(GetSchar32(alarm.sDeviceName)),
			DVRSerialNo: string(GetUchar48(alarm.sSerialNumber)),
			Mobile:      v.FDDesc.PhoneNumber,
			NowTime:     time.Now().Unix(),
			FDID:        int(match.struBlockListInfo.struBlockListInfo.dwGroupNo),
			PID:         int(match.struBlockListInfo.struBlockListInfo.dwRegisterID),
			Realname:    uname,
			Username:    string(GetUchar32(match.struBlockListInfo.struBlockListInfo.struAttribute.byCertificateNumber)),
		}
		snapPicLen := int(match.struSnapInfo.dwSnapFacePicLen)
		snapPicBuf := (*C.char)(unsafe.Pointer(match.struSnapInfo.pBuffer1))
		if snapPicLen > 0 && snapPicBuf != nil {
			if int(match.byPicTransType) > 0 {
				lBuf := capStream(snapPicBuf, C.int(snapPicLen))
				if string(lBuf) != "" {
					result, err := UploadImg(lBuf, UserApi.Device.OssPath)
					if err != nil {
						record.LogContent = "保存抓拍图片失败" + err.Error()
						_, err := request.HttpPost(httpPath+RemoteLogUrl, util.MarshalJson(record), nil)
						if err != nil {
							saveDeviceLog(record)
						}
					} else {
						record.Images = []string{result.Url}
						record.LogContent = util.MarshalJson(result)
					}
				} else {
					record.LogContent = "保存抓拍图片失败"
					_, err := request.HttpPost(httpPath+RemoteLogUrl, util.MarshalJson(record), nil)
					if err != nil {
						saveDeviceLog(record)
					}
				}
			} else {
				lBuf := C.GoBytes(unsafe.Pointer(match.struSnapInfo.pBuffer1), C.int(snapPicLen))
				result, err := UploadImg(lBuf, UserApi.Device.OssPath)
				if err != nil {
					record.LogContent = "保存抓拍图片失败" + err.Error()
					_, err := request.HttpPost(httpPath+RemoteLogUrl, util.MarshalJson(record), nil)
					if err != nil {
						saveDeviceLog(record)
					}
				} else {
					record.Images = []string{result.Url}
					record.LogContent = util.MarshalJson(result)
				}
			}
			_, err := request.HttpPost(httpPath+RemoteRecordUrl, util.MarshalJson(record), nil)
			if err != nil {
				saveDeviceRecord(record)
			}
		}
		blockPicLen := int(match.struBlockListInfo.dwBlockListPicLen)
		blockPicBuf := (*C.char)(unsafe.Pointer(match.struBlockListInfo.pBuffer1))
		if blockPicLen > 0 && blockPicBuf != nil {
			if int(match.byPicTransType) > 0 {
				lBuf := capStream(blockPicBuf, C.int(blockPicLen))
				if string(lBuf) != "" {
					result, err := UploadImg(lBuf, UserApi.Device.OssPath)
					if err != nil {
						record.LogContent = "保存黑名单人脸图片失败" + err.Error()
						_, err := request.HttpPost(httpPath+RemoteLogUrl, util.MarshalJson(record), nil)
						if err != nil {
							saveDeviceLog(record)
						}
					} else {
						record.Images = []string{result.Url}
						record.LogContent = util.MarshalJson(result)
					}
				} else {
					record.LogContent = "保存黑名单人脸图片失败"
					_, err := request.HttpPost(httpPath+RemoteLogUrl, util.MarshalJson(record), nil)
					if err != nil {
						saveDeviceLog(record)
					}
				}
			} else {
				lBuf := C.GoBytes(unsafe.Pointer(match.struBlockListInfo.pBuffer1), C.int(blockPicLen))
				result, err := UploadImg(lBuf, UserApi.Device.OssPath)
				if err != nil {
					record.LogContent = "保存黑名单人脸图片失败" + err.Error()
					_, err := request.HttpPost(httpPath+RemoteLogUrl, util.MarshalJson(record), nil)
					if err != nil {
						saveDeviceLog(record)
					}
				} else {
					record.Images = []string{result.Url}
					record.LogContent = util.MarshalJson(result)
				}
			}
			_, err := request.HttpPost(httpPath+RemoteRecordUrl, util.MarshalJson(record), nil)
			if err != nil {
				saveDeviceRecord(record)
			}
		}
	}
	return false
}

func capStream(pBuf *C.char, pLen C.int) []byte {
	newHost := "https://" + UserApi.Device.Address + ":443"
	rawurl := strings.Replace(C.GoStringN(pBuf, pLen), "\n", "", -1)
	paths, _ := url.Parse(rawurl)

	newUrl := newHost + paths.Path + "?" + paths.Query().Encode()
	method := "POST"
	req, err := http.NewRequest(method, newUrl, nil)
	req.Header.Set("Content-Type", "application/json")
	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("%v", err)
		return []byte("")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		log.Printf("Recieved status code '%v' auth skipped", resp.StatusCode)
		return []byte("")
	}
	digestParts := digestParts(resp)
	digestParts["uri"] = paths.Path + "?" + paths.Query().Encode()
	digestParts["method"] = method
	digestParts["username"] = UserApi.Device.Username
	digestParts["password"] = UserApi.Device.Password
	req, err = http.NewRequest(method, newUrl, bytes.NewBuffer([]byte("")))
	req.Header.Set("Authorization", getDigestAuthrization(digestParts))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		log.Printf("%v", err)
		return []byte("")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []byte("")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%v", err)
		return []byte("")
	}
	return body
}

func digestParts(resp *http.Response) map[string]string {
	result := map[string]string{}
	if len(resp.Header["Www-Authenticate"]) > 0 {
		wantedHeaders := []string{"nonce", "realm", "qop"}
		responseHeaders := strings.Split(resp.Header["Www-Authenticate"][0], ",")
		for _, r := range responseHeaders {
			for _, w := range wantedHeaders {
				if strings.Contains(r, w) {
					result[w] = strings.Split(r, `"`)[1]
				}
			}
		}
	}
	return result
}

func getMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getCnonce() string {
	b := make([]byte, 8)
	io.ReadFull(rand.Reader, b)
	return fmt.Sprintf("%x", b)[:16]
}

func getDigestAuthrization(digestParts map[string]string) string {
	d := digestParts
	ha1 := getMD5(d["username"] + ":" + d["realm"] + ":" + d["password"])
	ha2 := getMD5(d["method"] + ":" + d["uri"])
	nonceCount := 00000001
	cnonce := getCnonce()
	response := getMD5(fmt.Sprintf("%s:%s:%v:%s:%s:%s", ha1, d["nonce"], nonceCount, cnonce, d["qop"], ha2))
	authorization := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc="%v", qop="%s", response="%s"`,
		d["username"], d["realm"], d["nonce"], d["uri"], cnonce, nonceCount, d["qop"], response)
	return authorization
}

func saveDeviceRecord(device DeviceRecord) bool {
	var hkRecord []gdb.DeviceRecord
	hkRecord = append(hkRecord, gdb.DeviceRecord{
		Id:         util.RandomUUID(),
		DeviceId:   device.DVRSerialNo,
		PictureUrl: device.Images[0],
		UserId:     "",
		SchoolId:   "",
		AttendTime: device.NowTime * 1000,
		AttendType: 300, // 刷脸
		Status:     0,   // 打卡成功
	})
	if len(hkRecord) > 0 {
		data := hkRecord[0]
		if gdb.DB.Table(gdb.TableNameDevice).CreateInBatches(hkRecord, len(hkRecord)).RowsAffected <= 0 {
			log.Printf("%s", "批量保存数据失败")
		}
		dlog := gdb.DeviceLog{}
		dlog.Id = util.RandomUUID()
		dlog.Username = device.Username
		dlog.LogContent = util.MarshalJson(device)
		dlog.DeviceId = data.DeviceId
		dlog.RecordId = data.Id
		gdb.DB.Table(gdb.TableNameDeviceLog).Create(&dlog)
	}
	return true
}

func saveDeviceLog(device DeviceRecord) bool {
	var hkLog []gdb.DeviceLog
	hkLog = append(hkLog, gdb.DeviceLog{
		Id:         util.RandomUUID(),
		DeviceId:   device.DVRSerialNo,
		SchoolId:   "",
		UserId:     "",
		Username:   device.Username,
		LogType:    1, // 错误日志类型
		LogContent: util.MarshalJson(device),
	})
	if len(hkLog) > 0 {
		if gdb.DB.Table(gdb.TableNameDeviceLog).CreateInBatches(hkLog, len(hkLog)).RowsAffected <= 0 {
			log.Printf("%s", "批量保存数据失败")
		}
	}
	return true
}
