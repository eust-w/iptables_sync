package ctrl

import (
	"fmt"
	"github.com/coreos/go-iptables/iptables"
)

func GetIptablesList(ipt *iptables.IPTables, table, chain string) ([]string, error) {
	list, err := ipt.List(table, chain)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func FindIptablesPosition(ipt *iptables.IPTables, table, chain, rule string) (int, error) {
	list, err := GetIptablesList(ipt, table, chain)
	if err != nil {
		return 0, err
	}
	for i, v := range list {
		if v == rule {
			return i, nil
		}
	}
	return 0, fmt.Errorf("not found rule")
}

func InsertIptablesRuleBeforeTargetRule(ipt *iptables.IPTables, table, chain, targeRule string, rule []string) error {
	pos, err := FindIptablesPosition(ipt, table, chain, targeRule)
	if err != nil {
		return fmt.Errorf("find iptables position: %s", err)
	}
	err = ipt.Insert(table, chain, pos, rule...)
	if err != nil {
		return fmt.Errorf("insert rule: %s", err)
	}
	return nil
}

func UniqueInsertIptablesRuleBeforeTargetRule(ipt *iptables.IPTables, table, chain, targeRule string, rule []string) error {
	pos, err := FindIptablesPosition(ipt, table, chain, targeRule)
	if err != nil {
		return fmt.Errorf("find iptables position: %s", err)
	}
	err = ipt.InsertUnique(table, chain, pos, rule...)
	if err != nil {
		return fmt.Errorf("insert rule: %s", err)
	}
	return nil
}
