package email

import (
	"fmt"
        "log"
//        "net"
        "strconv"
	"net/smtp"
//        "net/mail"
        "crypto/tls"
	"encoding/base64"
)

var (
	e SMTPInfo
)

// SMTPInfo is the details for the SMTP server
type SMTPInfo struct {
	Username string
	Password string
	Hostname string
	Port     int
	From     string
}

// Configure adds the settings for the SMTP server
func Configure(c SMTPInfo) {
	e = c
}

// ReadConfig returns the SMTP information
func ReadConfig() SMTPInfo {
	return e
}

// SendEmail sends an email
func SendEmail(to, subject, body string) error {
        host       := e.Hostname
        servername := host+":"+strconv.Itoa(e.Port)   
        from       := e.From
	auth := smtp.PlainAuth("", e.Username, e.Password, host)

	header := make(map[string]string)
	header["From"] = from
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = `text/plain; charset="utf-8"`
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

  // TLS config
    tlsconfig := &tls.Config {
        InsecureSkipVerify: true,
        ServerName: host,
    }
      conn, err := tls.Dial("tcp", servername, tlsconfig)
    if err != nil {
        log.Panic(err)
    }
    c, err := smtp.NewClient(conn, host)
    if err != nil {
        log.Panic(err)
    }
    // Auth
    if err = c.Auth(auth); err != nil {
        log.Panic(err)
    }
    // To && From
    if err = c.Mail(from); err != nil {
        log.Panic(err)
    }
    if err = c.Rcpt(to); err != nil {
        log.Panic(err)
    }
    // Data
    w, err := c.Data()
    if err != nil {
        log.Panic(err)
    }

    _, err = w.Write([]byte(message))
    if err != nil {
        log.Panic(err)
    }

    err = w.Close()
    if err != nil {
        log.Panic(err)
    }
    c.Quit()

	// Send the email
/*
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", e.Hostname, e.Port),
		auth,
		e.From,
		[]string{to},
		[]byte(message),
	)
*/
	return err
}
