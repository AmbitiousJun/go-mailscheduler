package mailscheduler

import (
	"errors"
	"log"

	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
)

// Scheduler 是核心的任务执行器
type Scheduler struct {
	Cron string // 任务的执行时间 cron 表达式

	MailOptions *MailOptions // 邮件选项

	SmtpOptions *SmtpOptions // smtp 服务器选项

	message *gomail.Message // 根据用户传入的选项构造邮件报文

	dialer *gomail.Dialer // 根据用户传入的选项初始化 smtp 连接

	cronExec *cron.Cron   // cron 表达式执行器
	entryId  cron.EntryID // 初始化 cron 执行器时获取到的执行 id，用于停止定时任务
}

// New 用于初始化一个定时邮件执行器
func New(cron string, mOpt *MailOptions, sOpt *SmtpOptions) (*Scheduler, error) {
	if mOpt == nil || sOpt == nil {
		return nil, errors.New("mail or smtp options should not be nil")
	}

	s := Scheduler{
		Cron:        cron,
		SmtpOptions: sOpt,
		MailOptions: mOpt,
	}

	if err := s.genCronExecutor(); err != nil {
		return nil, errors.Join(err, errors.New("generate cron executor failed"))
	}

	s.genMailMessage()

	s.genSmtpDialer()

	return &s, nil
}

// Start 开启定时器
func (s *Scheduler) Start() {
	s.cronExec.Start()
}

// Stop 停止定时器
func (s *Scheduler) Stop() {
	s.cronExec.Remove(s.entryId)
	s.cronExec.Stop().Done()
}

// Send 手动发送邮件
// fallback == true 时，将会发送失败邮件
func (s *Scheduler) Send(fallback bool) error {
	if fallback {
		return s.sendFallback()
	}
	return s.sendNormal()
}

// sendFallback 发送一封失败邮件给收件人
func (s *Scheduler) sendFallback() error {
	f := s.MailOptions.FallbackBodyBuildFunc
	if f == nil {
		return errors.New("no fallback body build function")
	}

	s.message.SetBody(string(s.MailOptions.BodyType), f())

	return s.sendMessage()
}

// sendNormal 发送一封正常的邮件给收件人
func (s *Scheduler) sendNormal() error {
	f := s.MailOptions.BodyBuildFunc
	if f == nil {
		return errors.New("no mail body build function")
	}

	body, err := f()
	if err != nil {
		return errors.Join(err, errors.New("build mail body failed"))
	}

	s.message.SetBody(string(s.MailOptions.BodyType), body)

	return s.sendMessage()
}

// sendMessage 将封装好的邮件报文发送出去
func (s *Scheduler) sendMessage() error {
	if err := s.dialer.DialAndSend(s.message); err != nil {
		return errors.Join(err, errors.New("send message failed"))
	}
	return nil
}

// genCronExecutor 初始化 cron 执行器
func (s *Scheduler) genCronExecutor() error {
	s.cronExec = cron.New()

	id, err := s.cronExec.AddFunc(s.Cron, func() {

		sendSuccess := false

		// 发送正常邮件（重试 3 次）
		for i := 1; i <= 3; i++ {
			log.Printf("try to send mail to: %v, current try: %d\n", s.MailOptions.To, i)

			err := s.sendNormal()
			if err != nil {
				log.Printf("send failed, error: %v", err)
				continue
			}

			sendSuccess = true
			break
		}

		if sendSuccess {
			log.Println("mail send success!")
			return
		}

		log.Println("mail send failed, try to send fallback mail...")
		if err := s.sendFallback(); err != nil {
			log.Printf("send fallback mail failed, error: %v", err)
		}
	})

	if err != nil {
		return err
	}

	s.entryId = id
	return nil
}

// genSmtpDialer 根据用户传递的配置初始化一个 smtp 连接
func (s *Scheduler) genSmtpDialer() {
	s.dialer = gomail.NewDialer(
		s.SmtpOptions.Host,
		s.SmtpOptions.Port,
		s.SmtpOptions.Username,
		s.SmtpOptions.Credential,
	)
}

// genMailMessage 生成一个 Message 对象，设置好固定信息后存放在 Scheduler 对象中
// 该方法不会调用 BodyBuildFunc 函数，而是等到定时器触发时再生成最新的邮件内容
func (s *Scheduler) genMailMessage() {
	m := gomail.NewMessage()
	m.SetHeader("From", s.MailOptions.From)
	m.SetHeader("To", s.MailOptions.To...)
	m.SetHeader("Subject", s.MailOptions.Subject)
	s.message = m
}
