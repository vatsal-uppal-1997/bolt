package Error

import "log"

func HandleError(err error, msg string, shouldExit bool) {
	if err == nil {
		return
	}
	if shouldExit {
		log.Fatal(msg + " : " + err.Error())
	}
	log.Println(msg + " : " + err.Error())
}
