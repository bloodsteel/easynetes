package utils

// 表示API调用的状态
/*
1. success 表示API调用成功
2. fail 表示API调用失败; 通常是由于数据无效或者调用条件而被拒绝; 通俗讲此类调用整个程序并没有产生error类型, 而需要用户自己手动创建一个error对象; 比如 errors.New("密码复杂度不够") 、 errors.New("email字段必须提供")
3. error 表示API调用发生错误; 通常是在执行代码逻辑中产生了一个err != nil的错误
*/
const (
	StatusSuccess string = "success"
	StatusFail    string = "fail"
	StatusError   string = "error"
)

// 所有返回的数据 携带 自定义的 code, 方便 查询错误
// 格式为 httpCode-CustomCode
/*
1. 所有请求正常的API响应
	http_code都为200; code都为10000;
	不使用201(创建成功)、204(无响应内容,通常用来表示删除成功)此类的2xx的响应码,一律使用200;
	message字段在正常情况没有
2. 所有请求失败的API响应
	http_code都为4xx; code都为2000x;
	message字段告知用户为何失败, 比如 "用户名长度, 邮箱格式不正确",  "token过期了",  "权限不够"等等
3. 所有请求错误的API响应
	http_code都为500; code都为3000x;
	message字段为程序内部err的信息
*/
const (
	SCodeOK                                 string = "200-10000"
	SCodeNotFoundWithDao                    string = "404-20001"
	SCodeForbidden                          string = "403-20002"
	SCodeUnauthenticateWithLogin            string = "401-20003"
	SCodeUnauthenticateWithNotExpired       string = "401-20004"
	SCodeUnauthenticateNotFoundClaims       string = "401-20005"
	SCodeUnauthenticateGtTolerationTime     string = "401-20006"
	SCodeBadRequestWithPathParamErr         string = "400-20007"
	SCodeBadRequestWithQueryParamErr        string = "400-20008"
	SCodeBadRequestWithAZEmpty              string = "400-20009"
	SCodeBadRequestWithHostIP               string = "400-20010"
	SCodeBadRequestWithIP4_6                string = "400-20011"
	SCodeBadRequestWithIPLen                string = "400-20012"
	SCodeBadRequestWithIPCalc               string = "400-20013"
	SCodeBadRequestWithIPErr                string = "400-20014"
	SCodeBadRequestWithIPEmpty              string = "400-20015"
	SCodeBadRequestWithNameRe               string = "400-20016"
	SCodeBadRequestWithNameSlash            string = "400-20017"
	SCodeBadRequestWithValueEmpty           string = "400-20018"
	SCodeBadRequestWithUserNameEmailPhoneRe string = "400-20019"
	SCodeBadRequestWithGroupName            string = "400-20020"
	SCodeBadRequestWithUserPwdRe            string = "400-20021"
	SCodeBadRequestWithSvcChildNodeEmpty    string = "400-20022"
	SCodeBadRequestWithParentIDEmpty        string = "400-20023"
	SCodeBadRequestWithArrayEmpty           string = "400-20024"
	SCodeBadRequestWithDomainRe             string = "400-20025"
	SCodeBadRequestWithDNSRecordNotEmpty    string = "400-20026"
	SCodeInternalServerErrorWithDao         string = "500-30001"
	SCodeInternalServerErrorWithPayloadBind string = "500-30002"
	SCodeInternalServerErrorWithGenJwtToken string = "500-30004"
	SCodeInternalServerErrorWithIPNetParse  string = "500-30003"
	SCodeUnknow                             string = "500-40001"
)

// Msg code 对应的 message
var Msg = map[string]string{
	SCodeOK:                                 "请求处理成功",
	SCodeNotFoundWithDao:                    "数据库记录没找到",
	SCodeForbidden:                          "没权限, 访问被拒绝",
	SCodeUnauthenticateWithLogin:            "登录失败",
	SCodeUnauthenticateWithNotExpired:       "JWT/Token没过期,但发生了错误, 请联系管理员",
	SCodeUnauthenticateNotFoundClaims:       "JWT/Token Claims没找到, 请联系管理员",
	SCodeUnauthenticateGtTolerationTime:     "JWT/Token过期时间已经大于容忍时间",
	SCodeBadRequestWithPathParamErr:         "请求路径参数错误",
	SCodeBadRequestWithQueryParamErr:        "请求查询参数错误",
	SCodeBadRequestWithAZEmpty:              "可用区不能为空",
	SCodeBadRequestWithHostIP:               "给定的主机IP已被使用或者IP可用区不匹配",
	SCodeBadRequestWithIP4_6:                "IPv4和IPv6不能同时存在",
	SCodeBadRequestWithIPLen:                "CIDR地址错误",
	SCodeBadRequestWithIPCalc:               "IPv4范围太大, 计算失败",
	SCodeBadRequestWithIPErr:                "IP地址无效",
	SCodeBadRequestWithIPEmpty:              "必须指定一个IP地址(v4或者v6)",
	SCodeBadRequestWithNameRe:               "name字段不符合正则表达式规则 ^[a-zA-Z][-._a-zA-Z0-9]*$ ",
	SCodeBadRequestWithNameSlash:            "name字段中不能包含'/'",
	SCodeBadRequestWithValueEmpty:           "value字段不能为空",
	SCodeBadRequestWithUserNameEmailPhoneRe: "用户名(^[a-zA-Z0-9_-]{4,16}$)/邮箱/Phone不符合正则表达式规则",
	SCodeBadRequestWithGroupName:            "组名不符合正则表达式规则 ^[a-zA-Z0-9\\p{L}_-]{2,16}$ ",
	SCodeBadRequestWithUserPwdRe:            "用户密码不符合规范",
	SCodeBadRequestWithDomainRe:             "域名不符合RFC规范,请参考RFC-1123",
	SCodeBadRequestWithDNSRecordNotEmpty:    "该域名下的解析记录不为空",
	SCodeBadRequestWithSvcChildNodeEmpty:    "服务树子节点不为空",
	SCodeBadRequestWithParentIDEmpty:        "父ID不能为空",
	SCodeBadRequestWithArrayEmpty:           "数组为空, 或者数组中无可用对象",
	SCodeInternalServerErrorWithDao:         "dao操作失败, 请联系管理员",
	SCodeInternalServerErrorWithPayloadBind: "JSON数据解析失败",
	SCodeInternalServerErrorWithGenJwtToken: "生成JWT或者Token失败",
	SCodeInternalServerErrorWithIPNetParse:  "CIDR解析错误",
	SCodeUnknow:                             "未知错误, 请稍后重试",
}
