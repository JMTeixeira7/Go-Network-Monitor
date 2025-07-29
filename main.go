package main

import (
	"fmt"
	"applications_manager/Url"
)

func main() {
	Url.SetTargetURLs()

	urls := Url.GetTargetURLs()
	fmt.Println("Collected URLs:")
	for _, u := range urls {
		fmt.Println("-", u)
	}
		// search for open browsers

		// search for open tabs inside each browser

		// read defined target tabs

		// shutdown tab if it is a target tab

		// save the open tabs on data base every time it scans

}
