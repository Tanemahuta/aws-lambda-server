package config

import "github.com/aws/aws-sdk-go-v2/aws/arn"

// ArnAsString conversion.
func ArnAsString(arnPtr *arn.ARN) string {
	if arnPtr == nil {
		return ""
	}
	return arnPtr.String()
}
