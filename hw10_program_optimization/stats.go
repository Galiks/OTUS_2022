package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
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

var (
	ErrReaderIsNil   = errors.New("reader is nil")
	ErrDomainIsEmpty = errors.New("domain is empty")
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	if r == nil {
		return nil, ErrReaderIsNil
	}

	if domain == "" {
		return nil, ErrDomainIsEmpty
	}
	bf := bufio.NewReader(r)
	stats := make(DomainStat)

	for {
		line, _, err := bf.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		email, err := jsonparser.GetUnsafeString(line, "Email")
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(email, domain) {
			emailDomain := email[strings.IndexRune(email, '@')+1:]
			stats[strings.ToLower(emailDomain)]++
		}
	}
	return stats, nil
}

func GetDomainStatOld(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched, err := regexp.Match("\\."+domain, []byte(user.Email))
		if err != nil {
			return nil, err
		}

		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
