package utils

func StringsToAny(items ...string) []interface{} {
	var interfaces []interface{} = make([]interface{}, len(items))
	for index, arg := range items {
		interfaces[index] = arg
	}
	return interfaces
}
