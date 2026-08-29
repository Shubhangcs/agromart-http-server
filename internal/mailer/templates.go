package mailer

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

const brand = "South Canara Agro Mart"

func wrap(inner string) string {
	return `<div style="font-family:-apple-system,Segoe UI,Roboto,sans-serif;max-width:520px;margin:0 auto;padding:24px;color:#0F172A">` +
		`<h2 style="color:#0060B8;margin:0 0 16px">` + brand + `</h2>` + inner +
		`<p style="color:#64748B;font-size:12px;margin-top:32px">You received this email because you have an account on ` + brand + `.</p></div>`
}

// SendWelcome is fire-and-forget: a failed welcome email must never fail signup.
func SendWelcome(m Mailer, logger *slog.Logger, to, firstName string) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		subject := "Welcome to " + brand
		text := fmt.Sprintf("Hi %s,\n\nWelcome to %s — India's B2B marketplace for dry fruits and agri produce.\n\nBrowse products, follow sellers, raise RFQs and chat directly with businesses. When you're ready to sell, tap “Become a Seller” in the app.\n\nHappy trading!", firstName, brand)
		html := wrap(fmt.Sprintf(`<p>Hi %s,</p><p>Welcome to <b>%s</b> — India's B2B marketplace for dry fruits and agri produce.</p><p>Browse products, follow sellers, raise RFQs and chat directly with businesses. When you're ready to sell, tap <b>Become a Seller</b> in the app.</p><p>Happy trading!</p>`, firstName, brand))
		if err := m.Send(ctx, to, subject, text, html); err != nil {
			logger.Warn("welcome email failed", "to", to, "error", err)
		}
	}()
}

func ResetCodeEmail(firstName, code string) (subject, text, html string) {
	subject = "Your " + brand + " password reset code"
	text = fmt.Sprintf("Hi %s,\n\nYour password reset code is %s. It expires in 15 minutes.\n\nIf you did not request this, you can ignore this email.", firstName, code)
	html = wrap(fmt.Sprintf(`<p>Hi %s,</p><p>Your password reset code is:</p><p style="font-size:32px;font-weight:700;letter-spacing:8px;color:#0060B8">%s</p><p>It expires in 15 minutes. If you did not request this, you can ignore this email.</p>`, firstName, code))
	return
}
