package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"strings"
	"tour1/internal/word"
)

const (
	ModeUpper = iota + 1
	ModeLower
	ModeUnderscoreToUpperCamelcase
	ModeUnderscoreToLowerCamelcase
	ModeCamelcaseToUnderscore
)

var desc = strings.Join([]string{
	"支持单词转换格式转换",
	"1： 全部转大写",
	"2： 全部转小写",
	"3： 下划线转大写驼峰",
	"4： 下划线转小写驼峰",
	"5： 驼峰转下划线",
}, "\n")

var str string
var mode int8

var wordCmd = &cobra.Command{
	Use:   "word",
	Short: "单词格式转换",
	Long:  desc,
	Run: func(cmd *cobra.Command, args []string) {
		var content string
		switch mode {
		case ModeUpper:
			content = word.ToUpper(str)
		case ModeLower:
			content = word.ToLower(str)
		case ModeUnderscoreToUpperCamelcase:
			content = word.UnderscoreToUpperCamelCase(str)
		case ModeUnderscoreToLowerCamelcase:
			content = word.UnderscoreToLowerCamelCase(str)
		case ModeCamelcaseToUnderscore:
			content = word.CameCaseToUnderscore(str)
		default:
			log.Fatalf("not support")
		}
		log.Printf("out: %s\n", content)
	},
}

func init() {
	wordCmd.Flags().StringVarP(&str, "str", "s", "", "word")
	wordCmd.Flags().Int8VarP(&mode, "mode", "m", 0, "mode")
}
