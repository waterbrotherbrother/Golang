package main

import (
	"fmt"
	"os"
)

func main() {

	tmp := []string{"line 1 \n\f", "line 2 \n\f", "line 3 \n\f", "line 4 \n\f", "line 5 \n\f", "line 6 \n\f", "line 7 \n\f", "line 8 \n\f", "line 9 \n\f", "line 10 \n\f"}
	file, error := os.Create("data_f.txt")
	if error != nil {
		fmt.Println(error)
	}

	for i := 0; i < 10; i++ {
		file.WriteString(tmp[i])
	}

	file.Close()
}
