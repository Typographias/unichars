package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type char struct {
	unicode string
	value   string
}

func (c *char) toString() string {
	s, err := strconv.Unquote(c.unicode)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func toSingleQuotedUnicodeCode(code string) string {
	return "'\\u" + code + "'"
}

func isUnicodeLine(text string) bool {
	return strings.Contains(text, "unicode = ")
}

type chars []char

func (c chars) contains(newItem char) bool {
	if len(c) == 0 {
		return false
	}

	for _, item := range c {
		if item.unicode == newItem.unicode {
			return true
		}
	}

	return false
}

func (c chars) toString() string {
	strings := ""

	for _, item := range c {
		strings += item.value
	}

	return strings
}

func main() {
	raw := flag.Bool("raw", false, "Return raw value (just strings with no linebreaks, useful for piping)")
	flag.Parse()

	filename := strings.Join(flag.Args()[0:1], "")

	if strings.HasSuffix(filename, ".glyphs") {
		if !*raw {
			fmt.Printf("\nOpening %s\n\n", filename)
		}

		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		var list chars

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := strings.Trim(scanner.Text(), " ")

			if isUnicodeLine(text) {
				newChar := char{
					unicode: toSingleQuotedUnicodeCode(strings.Trim(strings.Trim(text, "unicode = "), ";")),
				}
				newChar.value = newChar.toString()

				if !list.contains(newChar) {
					list = append(list, newChar)
				}
			}
		}

		listStrings := list.toString()

		if *raw {
			fmt.Printf("%s ", listStrings)
		} else {
			fmt.Printf("%s\n\n", listStrings)
		}

		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("We need a Glyphs App file to be able to run this.")
	}
}
