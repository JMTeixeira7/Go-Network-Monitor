package Url

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var urls []string

func SetTargetURLs(){
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("URL'S>");
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}

			line = strings.TrimSpace(line)

			if line == "END" { //Stops reading at the END mark (change if needed)
				break
			}

			urls = append(urls, line)
		}

		for i, line := range urls {
			fmt.Printf("Line %d: %s\n", i+1, line)
		}
}

func GetTargetURLs() []string{
	return urls
}

func GetTargetURL(index int) string{
	return urls[index]
}