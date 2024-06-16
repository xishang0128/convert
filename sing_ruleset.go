package main

import (
	"bytes"
	"os"

	"github.com/sagernet/sing-box/common/srs"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/json"
)

func SaveRuleSet(rules []option.DefaultHeadlessRule, outputPath string) error {
	plainRuleSet := option.PlainRuleSetCompat{
		Version: 1,
		Options: option.PlainRuleSet{
			Rules: common.Map(rules, func(it option.DefaultHeadlessRule) option.HeadlessRule {
				return option.HeadlessRule{
					Type:           C.RuleTypeDefault,
					DefaultOptions: it,
				}
			}),
		},
	}
	if err := saveSourceRuleSet(&plainRuleSet, outputPath+".json"); err != nil {
		return err
	}
	if err := saveBinaryRuleSet(&plainRuleSet, outputPath+".srs"); err != nil {
		return err
	}
	return nil
}

func saveSourceRuleSet(ruleset *option.PlainRuleSetCompat, outputPath string) error {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(ruleset); err != nil {
		return E.Cause(err, "encode config")
	}
	output, err := os.Create(outputPath)
	if err != nil {
		return E.Cause(err, "open output")
	}
	_, err = output.Write(buffer.Bytes())
	output.Close()
	if err != nil {
		return E.Cause(err, "write output")
	}
	return nil
}

func saveBinaryRuleSet(ruleset *option.PlainRuleSetCompat, outputPath string) error {
	ruleSet := ruleset.Upgrade()
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	err = srs.Write(outputFile, ruleSet)
	if err != nil {
		outputFile.Close()
		os.Remove(outputPath)
		return err
	}
	outputFile.Close()
	return nil
}
