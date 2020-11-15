package util

func InStringArray(array []string, item string) bool {
	for _, ele := range array {
		if item == ele {
			return true
		}
	}

	return false
}
