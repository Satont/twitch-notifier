package i18n

func GetNested[T any](v any, keys ...string) (T, bool) {
	res := v
	for _, key := range keys {
		mp, ok := res.(map[string]any)
		if !ok {
			var e T
			return e, false
		}
		res = mp[key]
	}
	a, ok := res.(T)
	return a, ok
}
