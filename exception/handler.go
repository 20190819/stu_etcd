package exception

import "log"

func Handler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
