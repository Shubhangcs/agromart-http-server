// Package mailer sends transactional email. Production uses Amazon SES (same AWS
// credentials as S3); when EMAIL_FROM is not configured a log-only mailer is used so
// local development never needs SES access.
package mailer

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/shubhangcs/agromart-server/internal/env"
)

type Mailer interface {
	Send(ctx context.Context, to, subject, textBody, htmlBody string) error
}

// New returns an SES mailer when EMAIL_FROM is set, otherwise a logging mailer.
func New(logger *slog.Logger) (Mailer, error) {
	from := env.GetString("EMAIL_FROM", "")
	if from == "" {
		logger.Warn("EMAIL_FROM not set: emails will be logged, not sent")
		return &logMailer{logger: logger}, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(env.GetString("REGION", "ap-south-1")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			env.GetString("ACCESS_KEY", ""), env.GetString("SECRET_KEY", ""), "")),
	)
	if err != nil {
		return nil, err
	}
	return &sesMailer{client: sesv2.NewFromConfig(awsCfg), from: from, logger: logger}, nil
}

type sesMailer struct {
	client *sesv2.Client
	from   string
	logger *slog.Logger
}

func (m *sesMailer) Send(ctx context.Context, to, subject, textBody, htmlBody string) error {
	_, err := m.client.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(m.from),
		Destination:      &types.Destination{ToAddresses: []string{to}},
		Content: &types.EmailContent{Simple: &types.Message{
			Subject: &types.Content{Data: aws.String(subject), Charset: aws.String("UTF-8")},
			Body: &types.Body{
				Text: &types.Content{Data: aws.String(textBody), Charset: aws.String("UTF-8")},
				Html: &types.Content{Data: aws.String(htmlBody), Charset: aws.String("UTF-8")},
			},
		}},
	})
	if err != nil {
		m.logger.Error("ses send failed", "to", to, "subject", subject, "error", err)
	}
	return err
}

type logMailer struct{ logger *slog.Logger }

func (m *logMailer) Send(_ context.Context, to, subject, textBody, _ string) error {
	m.logger.Info("email (not sent: EMAIL_FROM unset)", "to", to, "subject", subject, "body", textBody)
	return nil
}
