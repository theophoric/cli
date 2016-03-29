package formatters

import (
	. "github.com/theophoric/cf-cli/cf/i18n"
)

func Allowed(allowed bool) string {
	if allowed {
		return T("allowed")
	}
	return T("disallowed")
}
