package common

// DeleteSliceElement 删除slice中存在deleteElements的元素
func DeleteSliceElement(origin []string, deleteElements ...string) []string {
	delMap := make(map[string]interface{})
	for _, item := range deleteElements {
		delMap[item] = nil
	}
	start := 0
	for i, item := range origin {
		if _, ok := delMap[item]; !ok {
			origin[start] = origin[i]
			start += 1
		}
	}
	return origin[0:start]
}
