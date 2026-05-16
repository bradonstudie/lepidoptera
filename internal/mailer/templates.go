package mailer

import (
	"fmt"
	"time"
)

const baseURL = "https://lepidoptera.com"

type ShowEmailData struct {
	Title       string
	Venue       string
	Date        time.Time
	Slug        string
	Description string
	TicketURL   string
	Bands       []BandEmailData
}

type BandEmailData struct {
	Name        string
	IsHeadliner bool
}

func ConfirmSubscriptionEmail(token string) Email {
	return Email{
		Subject: "confirm your subscription",
		Body: fmt.Sprintf(`hey,

click the link below to confirm and start getting show alerts.

%s/confirm?token=%s

if you didn't sign up for this, just ignore it.

---
lepidoptera.com`, baseURL, token),
	}
}

func SubscriptionConfirmedEmail(unsubToken string) Email {
	return Email{
		Subject: "you're subscribed",
		Body: fmt.Sprintf(`you're in.

you'll get an email when new shows are posted and a reminder the day before.

---
unsubscribe: %s/unsubscribe?token=%s`, baseURL, unsubToken),
	}
}

func ShowAnnouncedEmail(show ShowEmailData, unsubToken string) Email {
	return Email{
		Subject: fmt.Sprintf("show alert: %s @ %s — %s",
			show.Title, show.Venue, show.Date.Format("Mon Jan 2")),
		Body: fmt.Sprintf(`%s is playing %s on %s.

%s

on the bill:
%s
%s
---
more info: %s/shows/%s

unsubscribe: %s/unsubscribe?token=%s`,
			show.Title,
			show.Venue,
			show.Date.Format("Monday January 2 at 3:04pm"),
			show.Description,
			formatBands(show.Bands),
			formatTicketURL(show.TicketURL),
			baseURL, show.Slug,
			baseURL, unsubToken),
	}
}

func ReminderEmail(show ShowEmailData, unsubToken string) Email {
	return Email{
		Subject: fmt.Sprintf("tomorrow: %s @ %s", show.Title, show.Venue),
		Body: fmt.Sprintf(`just a reminder — %s is tomorrow at %s.

%s
%s
---
more info: %s/shows/%s

unsubscribe: %s/unsubscribe?token=%s`,
			show.Title,
			show.Venue,
			show.Date.Format("3:04pm"),
			formatTicketURL(show.TicketURL),
			baseURL, show.Slug,
			baseURL, unsubToken),
	}
}

func formatBands(bands []BandEmailData) string {
	out := ""
	for _, b := range bands {
		if b.IsHeadliner {
			out += fmt.Sprintf("- %s (headliner)\n", b.Name)
		} else {
			out += fmt.Sprintf("- %s\n", b.Name)
		}
	}
	return out
}

func formatTicketURL(url string) string {
	if url == "" {
		return ""
	}
	return fmt.Sprintf("tickets: %s\n", url)
}
