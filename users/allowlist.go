package users

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

var allowlist = []string{"bia_rm"}

func IsInAllowList(username string) bool {
	return contains(allowlist, username)
}
