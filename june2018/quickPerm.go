package main

// quickPerm is ported from https://www.quickperm.org/
func quickPerm(in []int, callback func([]int) error) {
	n := len(in)

	// Initialize p, 0 to len(in)+1
	p := make([]int, n+1)
	for i := 0; i < n+1; i++ {
		p[i] = i
	}

	err := callback(in)
	if err != nil {
		return
	}

	i := 1
	var j int
	for i < len(in) {
		p[i]--
		if i%2 == 1 {
			j = p[i]
		} else {
			j = 0
		}
		in[j], in[i] = in[i], in[j]

		err = callback(in)
		if err != nil {
			return
		}

		i = 1
		for p[i] == 0 {
			p[i] = i
			i++
		}
	}
}
