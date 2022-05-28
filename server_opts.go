package sequence

type serverOpts struct {
	webhookSecretKey []byte
}

type ServerOpt func(*serverOpts) error

func WithWebhookSecret(secret string) ServerOpt {
	return func(so *serverOpts) error {
		so.webhookSecretKey = []byte(secret)
		return nil
	}
}
