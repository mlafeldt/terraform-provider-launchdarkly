package launchdarkly

func stringPtr(v string) *string { return &v }

func stringList(v []interface{}) []string {
	list := make([]string, len(v))
	for i, elem := range v {
		list[i] = elem.(string)
	}
	return list
}
