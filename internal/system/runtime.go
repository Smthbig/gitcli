package system

func IsRestrictedRuntime() bool {
	return isRestrictedEnv()
}
