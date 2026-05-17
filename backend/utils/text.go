package utils

import (
	"github.com/abadojack/whatlanggo"
)

func IsEnglish(s string) bool {
	info := whatlanggo.Detect(s)
	return info.Lang == whatlanggo.Eng
}