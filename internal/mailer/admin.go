package mailer

import "fmt"

func AdminLoginEmail(token string) Email {
	return Email{
		Subject: "Admin Login Link",
		Body: fmt.Sprintf(`your login link:
			%s/admin/verify?token=%s

			this link expires in 15 minutes.
			if you didn't request this, ignore it.

			---
			lepidoptera`, baseURL, token),
	}
}
