package apitest

import (
	"fmt"
	ini "gopkg.in/ini.v1"
)

type IniCnf struct {
	defSection *ini.Section
	devSection *ini.Section
}

func NewIniCnf(file string) (*IniCnf, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{UnescapeValueDoubleQuotes: true}, file)
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
