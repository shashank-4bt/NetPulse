package auth

func NewAPIKeySecret() (raw, prefix, last4, hash string) {
	raw = "npk_" + NewSecret()
	prefix = raw[:8]
	last4 = raw[len(raw)-4:]
	hash = HashSecret(raw)
	return
}

func NewWebhookSecret() (raw, hint string) {
	raw = "nwh_" + NewSecret()
	hint = raw[:8]
	return
}
