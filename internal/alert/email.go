package alert

import (
	"fmt"
	"log"

	"github.com/resend/resend-go/v2"
)

type EmailAlerter struct {
	client *resend.Client
	from   string
}

func NewEmailAlerter(apiKey string, resendEmail string) *EmailAlerter {
	return &EmailAlerter{
		client: resend.NewClient(apiKey),
		from:   resendEmail,
	}
}

func (e *EmailAlerter) SendDownAlert(to string, monitorName string, url string) error {
	log.Printf("Alert: attempting to send email to %s for monitor %s", to, monitorName)

	params := &resend.SendEmailRequest{
		From:    e.from,
		To:      []string{to},
		Subject: fmt.Sprintf("🔴 Monitor Down: %s", monitorName),
		Html: fmt.Sprintf(`
			<h2>Monitor Alert: Downtime Detected</h2>
			<p>We detected that the following monitor is currently unreachable.</p>
			<p><strong>Monitor:</strong> %s</p>
			<p><strong>URL:</strong> %s</p>
			<p>We will continue to check and notify you when service is restored.</p>
		`, monitorName, url),
	}

	resp, err := e.client.Emails.Send(params)
	if err != nil {
		log.Printf("Alert: FAILED to send email: %v", err)
		return err
	}

	log.Printf("Alert: email sent successfully, id: %s", resp.Id)
	return nil
}

func (e *EmailAlerter) SendUpAlert(to string, monitorName string, url string) error {
	params := &resend.SendEmailRequest{
		From:    e.from,
		To:      []string{to},
		Subject: fmt.Sprintf("🟢 Monitor Recovered: %s", monitorName),
		Html: fmt.Sprintf(`
			<h2>Monitor Recovery Confirmed</h2>
			<p>The following monitor is responding normally again.</p>
			<p><strong>Monitor:</strong> %s</p>
			<p><strong>URL:</strong> %s</p>
			<p>All checks are passing at this time.</p>
		`, monitorName, url),
	}

	_, err := e.client.Emails.Send(params)
	if err != nil {
		log.Printf("Alert: failed to send up email: %v", err)
		return err
	}

	log.Printf("Alert: sent up email for %s to %s", monitorName, to)
	return nil
}
