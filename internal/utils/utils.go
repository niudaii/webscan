package utils

func RemoveDuplicate(list []string) []string {
	var set []string
	hashSet := make(map[string]struct{})
	for _, v := range list {
		hashSet[v] = struct{}{}
	}
	for k := range hashSet {
		// 去除空字符串
		if k == "" {
			continue
		}
		set = append(set, k)
	}
	return set
}
