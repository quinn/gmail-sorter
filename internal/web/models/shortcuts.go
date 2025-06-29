package models

import "strings"

var keysmaps = map[string]map[string][]rune{
	"querty": {
		"left":  []rune{'a', 's', 'd', 'f', 'q', 'w', 'e', 'r'},
		"right": []rune{'j', 'k', 'l', ';', 'u', 'i', 'o', 'p'},
	},
	"dvorak": {
		"left":  []rune{'a', 'o', 'e', 'u', ' ', ',', '.', 'p'},
		"right": []rune{'h', 't', 'n', 's', 'g', 'c', 'r', 'l'},
	},
}

var shortcutKeys = keysmaps["querty"]["right"]

func AssignShortcuts(links []ActionLink) []rune {
	shortcuts := make([]rune, len(links))
	used := make(map[rune]bool)

	labels := make([][]rune, len(links))
	for i, link := range links {
		labels[i] = []rune(strings.ToLower(link.Action().Label(link)))
	}

	// Step 1: Try to assign shortcuts by matching each letter position, left to right
	for pos := range 10 { // 10 is a safe upper bound for label length
		for i, label := range labels {
			if len(label) > pos && shortcuts[i] == 0 {
				ch := label[pos]
				for _, key := range shortcutKeys {
					if ch == key && !used[key] {
						shortcuts[i] = key
						used[key] = true
						break
					}
				}
			}
		}
	}

	// Step 2: Assign any remaining unused shortcut keys to unmatched links
	for i := range shortcuts {
		if shortcuts[i] == 0 {
			for _, key := range shortcutKeys {
				if !used[key] {
					shortcuts[i] = key
					used[key] = true
					break
				}
			}
		}
	}

	return shortcuts
}
