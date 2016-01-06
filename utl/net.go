package utl

func BodyOk(message string) map[string]string {
	result := map[string]string{}
	result["status"] = "OK"
	result["message"] = message
	return result
}

func BodyErr(message string) map[string]string {
	result := map[string]string{}
	result["status"] = "ERROR"
	result["message"] = message
	return result
}
