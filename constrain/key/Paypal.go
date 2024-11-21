package key

import "github.com/nbvghost/dandelion/library/environments"

func BaseURL() string {
	if environments.Release() {
		return "https://api-m.paypal.com"
	} else {
		return "https://api-m.sandbox.paypal.com"
	}
}
