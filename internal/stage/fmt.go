package stage

import "fmt"

func Println(a ...any) (n int, err error) {
	if Verbose() {
		return fmt.Println(a...)
	}
	return 0, nil
}

func Printf(format string, a ...any) (n int, err error) {
	if Verbose() {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}
