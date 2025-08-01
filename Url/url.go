package Url

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Url struct {
	//full_name string
	Domain string `json:"domain"`
	Protocol string `json:"protocol"`
	Path string `json:"path"`
	Target bool `json:"target"`
}

var Urls []Url

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

			new_url, err := CreateUrl(line)
			if(err != nil) { return }
			Urls = append(Urls, new_url)
		}

		for i, url := range Urls {
			fmt.Printf("URL %d: %s://%s/%s\n", i+1, url.Protocol, url.Domain, url.Path)
		}
}

func GetTargetURLs() []Url{
	return Urls
}

func GetTargetURL(index int) Url{
	return Urls[index]
}