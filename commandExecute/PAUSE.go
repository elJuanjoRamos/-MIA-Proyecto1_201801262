package commandExecute

import (
	"bufio"
	"fmt"
	"os"
)

func ExecutePause() {
	fmt.Println("==================================")
	fmt.Println("	Program in Pause		   ")
	fmt.Println("   Press 'Enter' to continue...   ")
	fmt.Println("==================================")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
