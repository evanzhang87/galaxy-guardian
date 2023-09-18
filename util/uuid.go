package util

import "time"

func FetchUuid() (string, error) {

	var uuid string

	err := ReadCommandTimeout(time.Second*5, func(line string) error {
		uuid = line
		return nil

	}, "bash", "-c", `sudo dmidecode | grep UUID`)

	return uuid, err
}
