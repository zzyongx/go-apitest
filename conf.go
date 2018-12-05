package apitest

import (
	"fmt"
	ini "gopkg.in/ini.v1"
	"io/ioutil"
	"strings"
)

type IniCnf struct {
	defSection *ini.Section
	devSection *ini.Section
}

func NewIniCnf(file string) (*IniCnf, error) {
	txt, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var lines []string
	for _, line := range strings.Split(string(txt), "\n") {
		if !strings.HasPrefix(line, "include") {
			lines = append(lines, line)
		}
	}
	txt = []byte(strings.Join(lines, "\n"))

	cfg, err := ini.LoadSources(ini.LoadOptions{UnescapeValueDoubleQuotes: true}, txt)
	if err != nil {
		return nil, err
	}

	return &IniCnf{
		defSection: cfg.Section(""),
		devSection: cfg.Section("dev"),
	}, nil
}

func MustNewIniCnf(file string) *IniCnf {

	if cnf, err := NewIniCnf(file); err != nil {
		panic(fmt.Sprintf("parse ini file %s error: %s", file, err))
	} else {
		return cnf
	}
}

func (this *IniCnf) GetString(key string) string {
	if this.devSection.HasKey(key) {
		return this.devSection.Key(key).Value()
	} else {
		return this.defSection.Key(key).Value()
	}
}
