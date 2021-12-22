package haikang

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

/**
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
*/
