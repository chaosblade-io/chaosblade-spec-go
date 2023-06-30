package util

func Split(path string) (dir, file string) {
	i := lastSlash(path)
	return path[:i+1], path[i+1:]
}
func lastSlash(s string) int {
	i := len(s) - 1
	for i >= 0 && s[i] != '\\' {
		i--
	}
	return i
}
