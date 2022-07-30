package parser

import "log"

func VerifyRequestFileObject(reqFile RequestsFile) bool {
	log.Println("Verifying requests file")
	if reqFile.Address == "" {
		log.Printf("No host given")
		return false
	}
	return true
}
