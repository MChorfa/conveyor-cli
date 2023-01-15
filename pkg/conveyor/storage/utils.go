package storage

import "log"

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
