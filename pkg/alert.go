package pkg

import (
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail"
	m "net/mail"
	"strings"
	"time"
)

func ComposeAlertSubject(url string) string {
	return fmt.Sprintf("URL CHECK FAILURE: %s", url)
}

func ComposeTextMessage(url string, record *StoreRecord) string {
	var b = strings.Builder{}

	fmt.Fprintf(&b, "URL CHECK FAILED FOR:")
	fmt.Fprintf(&b, "%s\r\n\r\n", url)
	fmt.Fprintf(&b, "FAILURE STARTED AT: %s\r\n", record.Start)
	fmt.Fprintf(&b, "LAST FAILURE AT:    %s\r\n", record.Last)
	fmt.Fprintf(&b, "URL CHECKED %d TIMES.\r\n", record.Count)

	return b.String()
}

func RenderTemplate(template string, ctx *pongo2.Context, html bool) string {

	res := NewResources("templates")
	message, _ := res.ReadText(template)

	tmpl, err := pongo2.FromString(message)
	PanicOnError(err)

	out, err := tmpl.Execute(*ctx)
	PanicOnError(err)

	if html {
		pre, err := premailer.NewPremailerFromString(out, premailer.NewOptions())
		PanicOnError(err)

		out, err = pre.Transform()
		PanicOnError(err)
	}

	return out
}

func ComposeHtmlMessage(url string, record *StoreRecord) string {
	data := pongo2.Context{
		"url":    url,
		"record": record,
	}

	return RenderTemplate("failure-email.html", &data, true)

}

func NewSmtpServer(host string, port int, user, password string) *mail.SMTPServer {
	server := mail.NewSMTPClient()
	server.Host = host
	server.Port = port
	server.Username = user
	server.Password = password
	server.Encryption = mail.EncryptionNone

	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	return server
}

func ParseEmailAddresses(addressString string) []string {
	if addressString == "" {
		return []string{}
	}
	addresses := strings.Split(addressString, ";")
	return addresses
}

func ValidEmail(email string, single bool) bool {
	var err error
	addresses := ParseEmailAddresses(email)
	if single && len(addresses) > 1 {
		return false
	}
	for _, emailAddr := range addresses {
		_, err = m.ParseAddress(emailAddr)
		if err != nil {
			return false
		}
	}
	return true
}

func NewAlertEmail(from, to, cc string) *mail.Email {
	email := mail.NewMSG()

	console.Indent()
	// technically we should be validating there is only 1 address.
	emailFrom := ParseEmailAddresses(from)
	console.Trace("Email - set from: %s\n", emailFrom[0])
	email.SetFrom(emailFrom[0])

	for _, emailTo := range ParseEmailAddresses(to) {
		console.Trace("Email - adding to: %s\n", emailTo)
		email.AddTo(emailTo)
	}

	for _, emailCc := range ParseEmailAddresses(cc) {
		console.Trace("Email - adding cc: %s\n", emailCc)
		email.AddCc(emailCc)
	}

	console.Dedent()

	return email
}

func SendEmailAlert(server *mail.SMTPServer, email *mail.Email, url string, record *StoreRecord) {
	email.SetSubject(ComposeAlertSubject(url))
	email.SetBody(mail.TextPlain, ComposeTextMessage(url, record))
	email.SetBody(mail.TextHTML, ComposeHtmlMessage(url, record))

	console.Trace(">>>>> email >>>>>>>>>>>>>>\n")
	console.Trace(email.GetMessage())
	console.Trace("<<<<< email <<<<<<<<<<<<<<\n")

	if email.Error != nil {
		console.Print("Email contruction contains errors!\n")
		PanicOnError(email.Error)
	}

	client, err := server.Connect()
	if err != nil {
		console.Print("SMTP server connect failed!\n")
		PanicOnError(err)
	}

	err = email.Send(client)
	if err != nil {
		console.Print("Email send failed!\n")
		PanicOnError(err)
	}
}

func SendEmailReport(server *mail.SMTPServer, email *mail.Email, message *ReportMessage) {

	email.SetSubject(message.Subject())
	email.SetBody(mail.TextPlain, message.ToText())
	email.SetBody(mail.TextHTML, message.ToHtml())

	console.Trace(">>>>> email >>>>>>>>>>>>>>")
	console.Trace(email.GetMessage())
	console.Trace("<<<<< email <<<<<<<<<<<<<<")

	if email.Error != nil {
		PanicOnError(email.Error)
	}

	client, err := server.Connect()
	if err != nil {
		PanicOnError(err)
	}

	err = email.Send(client)
	if err != nil {
		PanicOnError(err)
	}
}
