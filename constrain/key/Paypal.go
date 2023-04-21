package key

const baseURL = "https://api-m.paypal.com"
const sandboxBaseURL = "https://api-m.sandbox.paypal.com"

func BaseURL() string {
	return sandboxBaseURL
}
