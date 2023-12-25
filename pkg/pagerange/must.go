package pagerange

func must1[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}
