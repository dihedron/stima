package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/dihedron/stima/version"
	"github.com/fatih/color"
)

type Options struct {
}

type category struct {
	key     string
	count   int
	regexes []*regexp.Regexp
	help    []string
	color   func(a ...any) string
}

var categories = []category{
	{
		key: "Contributo scarso",
		regexes: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\sridott(a|e|o|i)`),
			regexp.MustCompile(`(?i)\sscars(a|e|o|i)`),
			regexp.MustCompile(`(?i)\snon apprezzabil(e|i)`),
			regexp.MustCompile(`(?i)\srar(a|e|o|i)`),
		},
		help:  []string{"ridott*", "scars*", "non appezzabil*", "rar*"},
		color: color.New(color.FgMagenta).SprintFunc(),
	},
	{
		key: "Contributo limitato",
		regexes: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\slimitat(o|a|i|e)`),
			regexp.MustCompile(`(?i)\smoderat(o|a|i|e)`),
			regexp.MustCompile(`(?i)\sparzial(e|i)`),
			regexp.MustCompile(`(?i)\sliev(e|i)`),
			regexp.MustCompile(`(?i)\spoc(o|a|hi|he)`),
			regexp.MustCompile(`(?i)\scontenut(o|a|i|e)`),
			regexp.MustCompile(`(?i)\snon consistent(i|e)`),
		},
		help:  []string{"limitat*", "moderat*", "parzial*", "liev*", "poc*", "contenut*", "non consistent*"},
		color: color.New(color.FgRed).SprintFunc(),
	},
	{
		key: "Contributo apprezzabile",
		regexes: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\sdiscret(o|a|i|e)`),
			regexp.MustCompile(`(?i)\sapprezzabil(e|i)`),
			regexp.MustCompile(`(?i)\smaggior parte`),
			regexp.MustCompile(`(?i)\sper lo pi(ù|u')`),
			regexp.MustCompile(`(?i)\sconsistent(e|i)`),
			regexp.MustCompile(`(?i)\sbuon(o|a|i|e)`),
			regexp.MustCompile(`(?i)\sspesso`),
			regexp.MustCompile(`(?i)\sattent(o|a|i|e)`),
		},
		help:  []string{"discret*", "apprezzabil*", "maggior parte", "per lo più", "consistent*", "buon*", "spesso", "attent*"},
		color: color.New(color.FgYellow).SprintFunc(),
	},
	{
		key: "Contributo significativo",
		regexes: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\ssignificativ(o|a|i|e)`),
			regexp.MustCompile(`(?i)\snotevol(e|i)`),
			regexp.MustCompile(`(?i)\sconsiderevol(e|i)`),
			regexp.MustCompile(`(?i)\sfort(e|i)`),
			regexp.MustCompile(`(?i)\selevat(o|a|i|e)`),
			regexp.MustCompile(`(?i)\srilevant(e|i)`),
			regexp.MustCompile(`(?i)\sottim(o|a|i|e)`),
			regexp.MustCompile(`(?i)\sfondamental(e|i)`),
			regexp.MustCompile(`(?i)\sdeterminant(e|i)`),
			//regexp.MustCompile(`(?i)\sapprezzabil(e|i)`),
		},
		help:  []string{"significativ*", "notevol*", "considerevol*", "fort*", "elevat*", "rilevant*", "ottim*", "fondamental*", "determinant*"},
		color: color.New(color.FgGreen).SprintFunc(),
	},
	{
		key: "Impegno, competenza, managerialità",
		regexes: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\sdot(e|i) managerial(e|i)`),
			regexp.MustCompile(`(?i)\scapacit(à|a') organizzativ(a|e)`),
			regexp.MustCompile(`(?i)\scompetenz(a|e) tecnic(a|he)`),
		},
		help:  []string{"dot* managerial*", "capacità organizzativ*", "competenz* tecnic*"},
		color: color.New(color.FgBlue).SprintFunc(),
	},
}

func main() {

	if len(os.Args) == 2 {
		if valueIn(os.Args[1], "version", "-version", "--version", "ver", "-ver", "--ver", "v", "-v", "--v") {
			version.Print(os.Stdout)
			return
		} else if valueIn(os.Args[1], "explain", "-explain", "--explain", "e", "-e", "--e", "show", "-show", "--show", "s", "-s", "--s") {
			for _, category := range categories {
				fmt.Printf("%-44s: ", category.color(category.key))
				for _, h := range category.help {
					fmt.Printf("%-35s", category.color(h))
				}
				fmt.Println()
			}
			return
		}
	}

	var inputs []io.Reader
	if len(os.Args) == 1 {
		inputs = append(inputs, os.Stdin)
	} else {
		for _, arg := range os.Args[1:] {
			if file, err := os.Open(arg); err != nil {
				fmt.Fprintf(os.Stderr, "cannot open %s: %v\n", arg, err)
				os.Exit(1)
			} else {
				defer file.Close()
				inputs = append(inputs, file)
			}
		}

		for _, input := range inputs {
			d, err := io.ReadAll(input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
				os.Exit(1)
			}

			data := strings.TrimSpace(string(d))

			var buffer bytes.Buffer
			for i, category := range categories {
				for _, regex := range category.regexes {
					matches := regex.FindAllStringIndex(data, -1)
					index := 0
					if matches != nil {
						for _, match := range matches {
							categories[i].count = categories[i].count + 1
							buffer.WriteString(data[index:match[0]])
							if strings.Contains(data[match[0]:match[1]], " ") {
								tokens := strings.Split(data[match[0]:match[1]], " ")
								for i, token := range tokens {
									tokens[i] = category.color(token)
								}
								buffer.WriteString(strings.Join(tokens, " "))
							} else {
								buffer.WriteString(category.color(data[match[0]:match[1]]))
							}
							index = match[1]
						}
						buffer.WriteString(data[index:])
						data = buffer.String()
						buffer.Reset()
					}
				}
			}
			fmt.Println("----------------------------------------------------------------")
			fmt.Println()
			fmt.Printf("%s", data)
			fmt.Println()
			fmt.Println()
			for i := len(categories) - 1; i >= 0; i-- {
				fmt.Printf("  %s: %d\n", categories[i].color(categories[i].key), categories[i].count)
				categories[i].count = 0
			}
			fmt.Println()
		}
		fmt.Println("----------------------------------------------------------------")
	}
}

func valueIn(value string, values ...string) bool {
	for _, v := range values {
		if value == v {
			return true
		}
	}
	return false
}
