package fileutil

import "strings"

// ファイル名に使えない・面倒なようなものを置換する
func SanitizeReplaceName(name string) string {
	rep := strings.NewReplacer(
		"?", "？",
		"!", "！",
		"*", "＊",
		"&", "＆",
		"\n", "",
		" ", "_",
		"　", "_",
		`\`, "_",
		"/", "_",
		":", "：",
		";", "；",
		"<", "＜",
		">", "＞",
		`"`, "_",
		`'`, "_",
		"|", "_",
		"(", "_",
		")", "）",
		"+", "＋",
	)
	return rep.Replace(name)
}
