package config

import (
	"fmt"
	"os"
	"strconv"
)

func GetString(k string) (string, error) {
	v, ok := os.LookupEnv(k)
	if !ok {
		return "", fmt.Errorf("'%s' is not found in os env var, please set it", k)
	}

	return v, nil
}

func GetInt(k string) (int, error) {
	vStr, err := GetString(k)
	if err != nil {
		return 0, fmt.Errorf("error get string:: %w", err)
	}

	v, err := strconv.Atoi(vStr)
	if err != nil {
		return 0, fmt.Errorf("error convert string '%s' to int:: %w", vStr, err)
	}

	return v, nil
}

func GetBool(k string) (bool, error) {
	vStr, err := GetString(k)
	if err != nil {
		return false, fmt.Errorf("error get string:: %w", err)
	}

	v, err := strconv.ParseBool(vStr)
	if err != nil {
		return false, fmt.Errorf("error convert string '%s' to bool:: %w", vStr, err)
	}

	return v, nil
}

func GetServiceName() (string, error) {
	return GetString("X_OTELTRACE_APP_SERVICE_NAME")
}

func GetAppVersion() (string, error) {
	return GetString("X_OTELTRACE_APP_VERSION")
}

func GetAppEnvironment() (string, error) {
	return GetString("X_OTELTRACE_APP_ENVIRONMENT")
}

func GetOtelOTLPNewrelicHost() (string, error) {
	return GetString("X_OTELTRACE_OTEL_OTLP_NEWRELIC_HOST")
}

func GetOtelOTLPNewrelicHeaderAPIKey() (string, error) {
	return GetString("X_OTELTRACE_OTEL_OTLP_NEWRELIC_HEADER_API_KEY")
}
