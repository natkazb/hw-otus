package hw10programoptimization

import (
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type UserMin struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	var user UserMin
	decoder := json.NewDecoder(r)
	var err error
	for {
		err = decoder.Decode(&user)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(user.Email, domain) {
			if at := strings.IndexByte(user.Email, '@'); at != -1 {
				result[strings.ToLower(user.Email[at+1:])]++
			}
		}
	}

	return result, nil
}
