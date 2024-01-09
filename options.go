package mailscheduler

// MailBodyType 限制了用户能够发送的邮件类型
type MailBodyType string

const (
	MailBodyPlain MailBodyType = "text/plain"
	MailBodyHtml  MailBodyType = "text/html"
)

// MailBodyBuildFunc 是邮件内容构造函数
// 邮件应该发送什么内容由调用方决定，且不一定唯一
type MailBodyBuildFunc func() (string, error)

// MailFallbackBodyBuildFunc 是失败邮件内容构造函数
// 当定时发送邮件失败时，会通过这个函数生成一封失败的邮件发送给收件人
type MailFallbackBodyBuildFunc func() string

// MailOptions 允许用户配置邮件的相关参数
type MailOptions struct {
	From                  string                    // 发件人邮箱
	To                    []string                  // 收件人邮箱，可以指定多个收件人
	Subject               string                    // 主题
	BodyType              MailBodyType              // 邮件内容类型
	BodyBuildFunc         MailBodyBuildFunc         // 邮件内容构造器
	FallbackBodyBuildFunc MailFallbackBodyBuildFunc // 失败邮件内容构造器
}

// SmtpOptions 是 smtp 服务器器的配置
type SmtpOptions struct {
	Host       string // 服务器主机地址
	Port       int    // 端口号
	Username   string // 登录用户名
	Credential string // 登录密钥（通行证）
}
