package mail

import (
	"github.com/Bakhram74/small_bank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)
	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "A test email"
	content := `
<h1>Hello World<h1>
<p>This is a test message from <a href="http://techschool.guru">Tech School</a></p>
`
	to := []string{"bakhram7493@gmail.com "}
	attachFiles := []string{"../Dockerfile"}
	err = sender.SendMail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)

}
