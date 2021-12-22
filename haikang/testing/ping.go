package testing

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
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go_device_haikang/utils/util"
)

type XmlSchema struct {
//	ExtendList ExtendList `xml:"PersonInfoExtendList"`
	FDDesc     FDDesc	  `xml:"FDDescription"`
}

type ExtendList struct {
	ExtendInfo []ExtendInfo `xml:"PersonInfoExtend"`
}

type ExtendInfo struct {
	Id     string `xml:"id"`
	Enable string `xml:"enable"`
}

type FDDesc struct {
	//	Name        string `xml:"name"`
	PhoneNumber string `xml:"phoneNumber"`
	//	Prompt      string `xml:"prompt"`
}

const xmlschema = ` 
<XmlSchema>
<PersonInfoExtendList>
	<PersonInfoExtend>
		<id>1</id>
		<enable>false</enable>
	</PersonInfoExtend>
	<PersonInfoExtend>
		<id>2</id>
		<enable>false</enable>
	</PersonInfoExtend>
	<PersonInfoExtend>
		<id>3</id>
		<enable>false</enable>
	</PersonInfoExtend>
	<PersonInfoExtend>
		<id>4</id>
		<enable>false</enable>
	</PersonInfoExtend>
</PersonInfoExtendList>
16377675381830746D833CE9F69245F7AC12EF31682A91616AE426306102402EA662B1337DE9C238
<FDDescription>
	<name>教育云考勤</name>
	<phoneNumber>188xxxxxxxx</phoneNumber>
	<prompt>欢迎光临</prompt>
</FDDescription>
</XmlSchema>
`


func digestPost(host string, uri string, postBody []byte) bool {
	url := host + uri
	method := "POST"
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Content-Type", "application/json")
	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Timeout: 30 * time.Second, Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		log.Printf("Recieved status code '%v' auth skipped", resp.StatusCode)
		return true
	}
	digestParts := digestParts(resp)
	digestParts["uri"] = uri
	digestParts["method"] = method
	digestParts["username"] = "admin"
	digestParts["password"] = "zxw123456"
	req, err = http.NewRequest(method, url, bytes.NewBuffer(postBody))
	req.Header.Set("Authorization", getDigestAuthrization(digestParts))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		log.Println("response body: ", string(body))
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result, _ := UploadImg(body, "https://oss.school.xxx.com/upload/image")
	fmt.Println(result)
	return true
}

type UploadImageResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
	Url     string `json:"url"`
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


// FindDevice查找设备
func FindDevice(ip string) bool {
	return checkDevice(ip)
}

func checkDevice(ip string) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:8090", ip), time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func stageDev() {
	var v XmlSchema
	xml.Unmarshal([]byte(xmlschema), &v)
	fmt.Println(v.FDDesc.PhoneNumber)
	/*
	rawurl := strings.Replace("http://192.168.88.101:80/picture/Streaming/tracks/303/?name=ch00003_00000001405018879334400036459&size=36459", "\n", "", -1)
	t := digest.NewTransport("admin", "zxw123456")
	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		log.Printf("%v", err)
	}
	resp, err := t.RoundTrip(req)
	if err != nil {
		log.Printf("%v", err)
	}

	fmt.Println(resp)
	//defer resp.Body.Close()
	//buf, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("%v", err)
	//}
	return []byte("")
	 */
	// curl -u "admin:zxw123456" --digest http://192.168.88.101:8000/picture/Streaming/tracks/303/?name=ch00003_00000001405018879334400036459&size=36459
	// digestPost("https://192.168.88.101:443", "/picture/Streaming/tracks/303/?name=ch00003_00000001405018879334400036459&size=36459", []byte(""))
	paths, _ := url.Parse("http://192.168.88.101:80/picture/Streaming/tracks/303/?name=ch00003_00000001405018879334400036459&size=36459")
	paths.Host = "https://192.168.88.101:443"
	newUrl := paths.Host+paths.Path+"?"+paths.Query().Encode()
	fmt.Println(newUrl, paths.Path+"?"+paths.Query().Encode())
	return
}
