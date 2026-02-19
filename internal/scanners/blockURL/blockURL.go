package blockURL

import (
	"net/http"

)

type Block struct{
	//I need a service here (use double dispatch at controller)? --> Change design for double dispatch 
}

func New() *Block {
	return &Block{}
}

func (b *Block) Scan(r *http.Request) (res bool, reasons []string) {
	panic("not yet implemented")
}