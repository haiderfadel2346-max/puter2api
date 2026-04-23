package claude

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"puter2api/internal/types"
)

// ParseToolCalls 解析文本中的工具调用
// يدعم: JSON على سطر واحد أو متعدد الأسطر، spaces مختلفة
func ParseToolCalls(text string) ([]types.ParsedToolCall, string) {
	// regex أكثر مرونة: يتعامل مع whitespace و newlines متعددة
	re := regexp.MustCompile(`(?s)<tool_call>\s*([\s\S]*?)\s*</tool_call>`)
	matches := re.FindAllStringSubmatch(text, -1)

	var calls []types.ParsedToolCall
	remainingText := text

	for i, match := range matches {
		jsonStr := strings.TrimSpace(match[1])
		var call types.ParsedToolCall
		if err := json.Unmarshal([]byte(jsonStr), &call); err == nil {
			if call.ID == "" {
				call.ID = fmt.Sprintf("toolu_%d_%d", time.Now().UnixNano(), i)
			}
			// تأكد إن Input مش فاضي
			if len(call.Input) == 0 {
				call.Input = json.RawMessage("{}")
			}
			calls = append(calls, call)
		}
		// امسح الـ tag سواء نجح الـ parsing أو لأ
		remainingText = strings.Replace(remainingText, match[0], "", 1)
	}

	remainingText = strings.TrimSpace(remainingText)
	return calls, remainingText
}
