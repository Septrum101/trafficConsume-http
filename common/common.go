package common

import (
	"os"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetDownloadUrl() ([]string, error) {
	f, err := os.ReadFile("urls.txt")
	if err != nil {
		return nil, err
	}

	u := strings.Split(strings.TrimSpace(string(f)), "\n")
	u = unique(u)
	return u, nil
}

func unique(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		val = strings.TrimSpace(val)
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	sort.Strings(u)
	if err := os.WriteFile("urls.txt", []byte(strings.Join(u, "\n")), 0644); err != nil {
		logrus.Errorln(err)
	}
	return u
}
