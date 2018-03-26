package validators

func CheckValidationErrors(err error, errors *[]string) {
	if err != nil {
		*errors = append(*errors, err.Error())
	}
}