package url

import (
	"fmt"
	"net/url"
	"errors"
)

func CreateUrl(absolute_path string) (Url, error){
	parsedUrl, err :=url.Parse(absolute_path)
	if err != nil {
		fmt.Println("Error while parsing the url", err)
		return Url{}, err
	}

	newUrl := Url{Protocol: parsedUrl.Scheme, Domain: parsedUrl.Host, Path: parsedUrl.Path, Target: true}

	if(newUrl.Protocol == "" || newUrl.Domain == "") { //maybe present a lis of allowed protocols
		fmt.Println("Url presented is missing arguments")
		return Url{}, errors.New("missing protocol or domain in URL")
	}
	return newUrl, nil
}