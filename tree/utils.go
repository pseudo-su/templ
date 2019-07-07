package tree

func JoinStr(args []string) string {
	str := ""
	for _, arg := range args {
		str += arg
	}
	return str
}
