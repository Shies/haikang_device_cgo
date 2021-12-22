package haikang

type Scheme struct {
	Address  string
	Username string
	Password string
	Port     int
	Channel  int
}

type UploadImageResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
	Url     string `json:"url"`
}

type ResolveSvrParam struct {
	ServerIP        string
	ServerPort      int
	DVRName         string
	DVRNameLen      int
	DVRSerialNumber string
	DVRSerialLen    int
}

type ResolveSvrDto struct {
	SGetIP string
	DwPort int
}

type Device struct {
	IsSave     int    `json:"is_save"`                                 // 保存设备
	Active     int    `json:"active"`                                  // 设备状态 0=离线 1=在线
	Address    string `gorm:"address"              json:"address"`     // 设备IP地址
	Port       int    `gorm:"port"                 json:"port"`        // 设备IP端口
	Username   string `gorm:"username"             json:"username"`    // 登录用户名
	Password   string `gorm:"password"             json:"password"`    // 登录密码
	Channel    int    `gorm:"channel"              json:"channel"`     // 设备通道号
	DeviceName string `gorm:"device_name"          json:"device_name"` // 设备名称
	SerialNo   string `gorm:"serial_no"            json:"serial_no"`   // 设备序列号
}

type DeviceRecord struct {
	Address     string   `json:"address"`
	Port        int      `json:"port"`
	Channel     int      `json:"channel"`
	DVRName     string   `json:"dvrname"`
	DVRSerialNo string   `json:"dvrserialno"`
	Images      []string `json:"images"`
	Mobile      string   `json:"mobile"`
	NowTime     int64    `json:"nowTime"`
	FDID        int      `json:"fdid"`
	PID         int      `json:"pid"`
	Realname    string   `json:"realname"`
	Username    string   `json:"username"`
	LogType     int      `json:"logType"`
	LogContent  string   `json:"logContent"`
}
