package main

func indexOf(list []string, item string) (int, bool) {
	for idx, elm := range list {
		if item == elm {
			return idx, false
		}
	}

	return -1, true
}