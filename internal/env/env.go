package env

import "github.com/Dcarbon/go-shared/libs/utils"

var ServerHost = utils.StringEnv("SERVER_HOST", "localhost:4001")
var ServerScheme = utils.StringEnv("SERVER_SCHEME", "http")

var StorageHost = utils.StringEnv("STORAGE_HOST", "")
