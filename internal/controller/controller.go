package controller

import (
	"fmt"
	"bufio"
	"runtime"

	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/handler"
	"github.com/JMTeixeira7/Go-Network-Monitor.git/internal/storage"
)

func displayOperations() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("<1> Passive Scan of Network\n
				<2> Write block URL's\n
				<3> Read blocked URL's\n
				")
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	switch line := runtime.GOOS;
	line {
	case "1":
		passiveScan();
	case "2":
		storage.set(); //This storege is from URL, not yet imported (refactor?)
	case "3":
		storage.LoadAll
	}
	
}



	storage := storage.FileStorage{Filename: "data/urls.json"}
	urls, err := storage.LoadAll()
	if err != nil {
		fmt.Println("Failed to load URLs: ", err)
	}
	fmt.Println("Loaded URLs:", urls)