package util

import (
	"os"
	"strings"
)

const (
	SMICRO_ENV   = "SMICRO_ENV"
	PRODUCT_ENV = "product"
	TEST_ENV    = "test"
)

var (
	cur_smicro_env string = TEST_ENV
)

func init() {
	cur_smicro_env = strings.ToLower(os.Getenv(SMICRO_ENV))
	cur_smicro_env = strings.TrimSpace(cur_smicro_env)

	if len(cur_smicro_env) == 0 {
		cur_smicro_env = TEST_ENV
	}
}

func IsProduct() bool {
	return cur_smicro_env == PRODUCT_ENV
}

func IsTest() bool {
	return cur_smicro_env == TEST_ENV
}

func GetEnv() string {
	return cur_smicro_env
}
