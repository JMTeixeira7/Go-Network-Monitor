package main

import (
	//"fmt"
	"applications_manager/Url"
	"fmt"
	"net/http"
)

func main() {

	storage := Url.FileStorage{Filename: "data/urls.json"}

	Url.SetTargetURLs()
	err := storage.Save(Url.GetTargetURLs())
	if err != nil {
		fmt.Println("Failed to save URLs: ", err)
	}

	urls, err := storage.Load()
	if err != nil {
		fmt.Println("Failed to load URLs: ", err)
	}

	fmt.Println("Loaded URLs:", urls)


}
