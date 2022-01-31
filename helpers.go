package main

import (
	"fmt"
	"reflect"
	"regexp"
)

func isValidUrl(s string)(bool) {
	urlRegex := regexp.MustCompile("((http|https)://)(www.)?[a-zA-Z0-9@:%._\\+~#?&//=]{2,256}\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%._\\+~#?&//=]*)")
	return urlRegex.MatchString(s)
}

func validateMaintainers(maintainers []Maintainers)(bool) {
	for _, m := range maintainers {
		emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

		if (!emailRegex.MatchString(m.Email) || len(m.Name) == 0) {
			return false
		}
	}

	return true
}

func areMetaDataFieldsCorrect(p Application)(bool) {
	_, licenseErr := regexp.MatchString("[A-Z][a-zA-Z]-[0-9].[0-9]", p.License)

	if (len(p.Title) == 0 || len(p.Company) == 0 || len(p.Description) == 0) {
		return false
	} else if (len(p.Website) == 0 || !isValidUrl(p.Website) || len(p.Source) == 0 || !isValidUrl(p.Source)) {
		return false
	} else if (licenseErr != nil || len(p.License) == 0) {
		return false
	}

	return validateMaintainers(p.Maintainers)
}

func filterApplications(apps *[]Application, key string, value string) {
	for _, app := range applications {
		// using reflects to access fields dynamically
		r := reflect.ValueOf(app)
		fmt.Println(r)
		f := reflect.Indirect(r).FieldByName(key)

		if f.String() == value {
			*apps = append(*apps, app)
		}
	}
}
