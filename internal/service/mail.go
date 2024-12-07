package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type Mailer interface {
	SendEmail(to, subject, body string) error
}

type SESEmailer struct {
	client *ses.Client
	from   string
}

func NewSESEmailer(ctx context.Context, from string) (Mailer, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-2"))
	if err != nil {
		slog.Error("falha ao carregar configuracao AWS: %v", err.Error(), err)
		return nil, fmt.Errorf("falha ao carregar configuracao AWS: %w", err)
	}

	client := ses.NewFromConfig(cfg)
	return &SESEmailer{
		client: client,
		from:   from,
	}, nil
}

func (m *SESEmailer) SendEmail(to, subject, body string) error {
	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
			Body: &types.Body{
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
		},
		Source: aws.String(m.from),
	}
	_, err := m.client.SendEmail(context.TODO(), input)
	if err != nil {
		slog.Error("falha ao enviar email: %v", err.Error(), err)
		return fmt.Errorf("falha ao enviar email: %w", err)
	}
	return nil
}
