package mailer

import (
	"bytes"
	"embed"
	"time"

	"github.com/wneessen/go-mail"

	// Import the html/template and text/template. Because these share the same
	// package name ("template") we need to disambiguate them and alias them to
	// ht and tt respectively.

	ht "html/template"
	tt "text/template"
)

// The comment directive in the format `//go:embed <path>` indicates
// to Go that we want to store the contents of the ./templates directory
// in the templatesFS embedded file system variable.

//go:embed "templates"
var templateFS embed.FS

// sender is the information for your emails ("Alice Smith <alice@example.com>")
type Mailer struct {
	client *mail.Client
	sender string
}

func New(host string, port int, username, password, sender string) (*Mailer, error) {
	client, err := mail.NewClient(
		host,
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithPort(port),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(5*time.Second),
	)

	if err != nil {
		return nil, err
	}

	mailer := &Mailer{
		client: client,
		sender: sender,
	}

	return mailer, nil
}

func (m *Mailer) Send(recipient string, templateFile string, data any) error {

	// Use the ParseFS() method from text/template to parse the required template file
	// from the embedded file system.

	textTmpl, err := tt.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Execute the named template "subject", passing in the dynamic data and storing the
	// result in a bytes.Buffer variable
	subject := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(subject, "subject", data)

	if err != nil {
		return err
	}

	// The same pattern goes for the "plainBody" template
	plainBody := new(bytes.Buffer)
	err = textTmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// Parsing the html
	htmlTmpl, err := ht.New("").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = htmlTmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMsg()

	err = msg.To(recipient)
	if err != nil {
		return err
	}

	err = msg.From(m.sender)
	if err != nil {
		return err
	}

	msg.Subject(subject.String())
	msg.SetBodyString(mail.TypeTextPlain, plainBody.String())
	msg.AddAlternativeString(mail.TypeTextHTML, htmlBody.String())

	for i := 1; i <= 3; i++ {

		// opens a connection to the SMTP server
		// sends the message
		// closes the connection
		err = m.client.DialAndSend(msg)
		if err == nil {
			return nil
		}

		// If it didn't work, sleep for a short time and retr
		if i != 3 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return err
}
