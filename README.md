# go-mailscheduler

A simple smtp mail scheduler build by gomail and cron.

## Install
```shell

go get -u github.com/AmbitiousJun/go-mailscheduler@v1.0.0

```

## Usage

```go

sOpt := mailscheduler.SmtpOptions{
    Host:       "smtp.qq.com",
    Port:       587,
    Username:   "your_mail@qq.com",
    Credential: "your credential",
}
mOpt := mailscheduler.MailOptions{
    From:     "your_mail@qq.com",
    To:       []string{"receiver1@qq.com", "receiver2@qq.com"},
    Subject:  "Test Subject",
    // Use text/html mail body
    BodyType: mailscheduler.MailBodyHtml,
    BodyBuildFunc: func() (string, error) {
        return "<html><body><h1>Hello World</h1></body></html>", nil
    },
    FallbackBodyBuildFunc: func() string {
        return "<html><body><h1>This is a fallback mail</h1></body></html>"
    },
}

// Send mail every 2 minutes
ms, err := mailscheduler.New("*/2 * * * *", &mOpt, &sOpt)
if err != nil {
    t.Error(err)
    return
}

// Send normal mail manually
ms.Send(false)

// Send fallback mail manually
ms.Send(true)

// Start scheduler
ms.Start()

// Wait 3 minutes
<-time.After(time.Minute * 3)

ms.Stop()
log.Println("scheduler stop.")

```
