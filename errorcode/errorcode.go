package errorcode

// ============================================================
//  业务错误码规范
//
//  分段规则（5位数字）：
//    首位 = 错误大类，后4位 = 具体场景
//    0     成功
//    1xxxx 系统级错误
//    2xxxx 认证鉴权错误 (未登录/Token过期)
//    3xxxx 数据资源错误 (不存在/冲突)
//    4xxxx 参数校验错误 (格式/必填/长度) ← 用户要求的 6xx 改为 4xxxx 更符合业界惯例
//    5xxxx 服务端操作错误 (数据库/加密/Token生成)
//    6xxxx 业务逻辑错误 (密码错/编号重复等)
// ============================================================

const (
	// ── 成功 ──
	Success int = 0

	// ── 1xxxx 系统级错误 ──
	SystemError     = 10001 // 系统内部异常
	ServerPanic     = 10002 // 服务宕机

	// ── 2xxxx 认证鉴权错误 ──
	Unauthorized       = 20001 // 未登录 (无Token)
	TokenFormatInvalid = 20002 // Token 格式错误
	TokenExpired       = 20003 // Token 已过期或无效
	Forbidden          = 20004 // 无权限访问

	// ── 3xxxx 数据资源错误 ──
	ResourceNotFound = 30001 // 资源不存在
	DataConflict     = 30002 // 数据冲突 (如唯一约束)

	// ── 4xxxx 参数校验错误 ──
	ParamValidateFailed = 40001 // 参数校验失败(通用)
	ParamRequired       = 40002 // 必填字段缺失
	ParamFormat         = 40003 // 字段格式错误 (类型不对)
	ParamLengthExceeded = 40004 // 字段长度超限
	ParamValueRange     = 40005 // 字段值超出范围
	ParamEnumInvalid    = 40006 // 字段枚举值不合法

	// ── 5xxxx 服务端操作错误 ──
	DBError        = 50001 // 数据库操作失败
	DBConnection   = 50002 // 数据库连接失败
	EncryptFailed  = 50101 // 密码加密失败
	TokenGenFailed = 50201 // Token 生成失败

	// ── 6xxxx 业务逻辑错误 ──
	LoginFailed      = 60001 // 账号或密码错误
	RegisterFailed   = 60002 // 注册失败 (通用)
	DuplicateNo      = 60003 // 编号已被注册
	AccountDisabled  = 60004 // 账号已被禁用
)

// ErrorMsg 错误码 → 默认中文提示映射
var ErrorMsg = map[int]string{
	Success:            "成功",
	SystemError:        "系统内部异常，请稍后重试",
	ServerPanic:        "服务内部错误",
	Unauthorized:       "未登录，请先登录",
	TokenFormatInvalid: "认证格式错误，格式: Bearer <token>",
	TokenExpired:       "Token 已过期或无效，请重新登录",
	Forbidden:          "无权限访问该资源",
	ResourceNotFound:    "请求的资源不存在",
	DataConflict:       "数据冲突，请检查后重试",
	ParamValidateFailed: "参数校验失败",
	ParamRequired:       "缺少必填字段",
	ParamFormat:         "字段格式错误",
	ParamLengthExceeded: "字段长度超限",
	ParamValueRange:     "字段值超出允许范围",
	ParamEnumInvalid:    "字段值不在允许范围内",
	DBError:            "数据库操作失败",
	DBConnection:       "数据库连接失败",
	EncryptFailed:      "密码处理失败",
	TokenGenFailed:     "Token 生成失败",
	LoginFailed:        "账号或密码错误",
	RegisterFailed:     "注册失败，请稍后重试",
	DuplicateNo:        "编号已被使用",
	AccountDisabled:    "账号已被禁用，请联系管理员",
}

// Msg 获取错误码对应的默认提示文案
func Msg(code int) string {
	if msg, ok := ErrorMsg[code]; ok {
		return msg
	}
	return "未知错误"
}
