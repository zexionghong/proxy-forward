package utils

func BelongsToList(lookup string, list []string) bool {
	for _, val := range list {
		if val == lookup {
			return true
		}
	}
	return false
}
