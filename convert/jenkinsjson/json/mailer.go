package json

import (
	harness "github.com/drone/spec/dist/go"
)

// ConvertMailer creates a Harness step for nunit plugin.
func ConvertMailer(node Node, arguments map[string]interface{}) *harness.Step {
	subject, _ := arguments["subject"].(string)
	to, _ := arguments["to"].(string)
	body, _ := arguments["body"].(string)

	convertMailer := &harness.Step{
		Id:   SanitizeForId(node.SpanName, node.SpanId),
		Name: "Mailer",
		Type: "plugin",
		Spec: &harness.StepPlugin{
			Image: "plugins/email",
			With: map[string]interface{}{
				"host":         "smtp.gmail.com",
				"port":         "587",
				"username":     "<username>",
				"password":     "<password>",
				"subject":      subject,
				"body":         body,
				"recipients":   to,
				"from.address": "<from_address>",
			},
		},
	}

	return convertMailer
}
