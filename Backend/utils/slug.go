// ============================================================
// 包名: utils
// 说明: 文章 Slug 生成工具
//       从标题生成 URL 友好的唯一标识
// ============================================================
package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// GenerateSlug 从标题生成 URL slug
// 英文标题：hello world → hello-world
// 中文标题：如果没有英文则用 "post-{timestamp}" 格式
func GenerateSlug(title string) string {
	title = strings.TrimSpace(title)

	// 提取英文/数字部分
	var engWords []string
	var currentWord strings.Builder
	hasEnglish := false

	for _, r := range title {
		if unicode.IsLetter(r) && r <= unicode.MaxASCII {
			currentWord.WriteRune(unicode.ToLower(r))
			hasEnglish = true
		} else if unicode.IsDigit(r) {
			currentWord.WriteRune(r)
		} else {
			if currentWord.Len() > 0 {
				engWords = append(engWords, currentWord.String())
				currentWord.Reset()
			}
		}
	}
	if currentWord.Len() > 0 {
		engWords = append(engWords, currentWord.String())
	}

	if hasEnglish && len(engWords) > 0 {
		slug := strings.Join(engWords, "-")
		// 清理多余连字符
		slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
		slug = strings.Trim(slug, "-")
		if len(slug) > 100 {
			slug = slug[:100]
		}
		return slug
	}

	// 纯中文等情况：用时间戳生成
	return fmt.Sprintf("post-%s", time.Now().Format("20060102150405"))
}
