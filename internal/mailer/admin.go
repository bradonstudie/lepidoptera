package mailer

import "fmt"

func AdminLoginEmail(token string) Email {
	return Email{
		Subject: "Admin Login Link",
		Body: fmt.Sprintf(`your login link:
			%s/admin/verify?token=%s

			This link expires in 15 minutes.
			If you didn't request this, it's probably time to go check the logs.
			`, baseURL, token),
	}
}
