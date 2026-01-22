package main


func findNthSpace(s string, n int) int {
    c := 0
    for i := 0; i < len(s); i++ {
        if s[i] == ' ' {
			c += 1
		}
		if n == c {
			return i;
		}
    }
	return -1;
    
}

