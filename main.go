package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/dihedron/stima/version"
	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Input *string `short:"i" long:"input" description:"The inline chunk of text to be used as input or a @file path." optional:"yes"`
}

type category struct {
	key     string
	count   int
	regexes []*regexp.Regexp
	color   func(a ...any) string
}

func main() {

	if len(os.Args) == 2 && valueIn(os.Args[1], "version", "-version", "--version", "ver", "-ver", "--ver", "v", "-v", "--v") {
		version.Print(os.Stdout)
		return
	}

	options := Options{}
	// Parse flags from `args'. Note that here we use flags.ParseArgs for
	// the sake of making a working example. Normally, you would simply use
	// flags.Parse(&opts) which uses os.Args
	_, err := flags.Parse(&options)
	if err != nil {
		slog.Error("error parsing command line", "error", err)
		panic(err)
	}

	slog.Debug("command line parsed", "options", options)

	var input io.Reader
	if options.Input == nil {
		input = os.Stdin
	} else {
		if strings.HasPrefix(*options.Input, "@") {
			*options.Input = strings.TrimPrefix(*options.Input, "@")
		}
		if file, err := os.Open(*options.Input); err != nil {
			panic(err)
		} else {
			defer file.Close()
			input = file
		}
	}

	d, err := io.ReadAll(input)
	if err != nil {
		panic(err)
	}

	data := string(d)

	categories := []category{
		{
			key: "Contributo scarso",
			regexes: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\sridott(a|e|o|i)`),
				regexp.MustCompile(`(?i)\sscars(a|e|o|i)`),
				regexp.MustCompile(`(?i)\snon apprezzabil(e|i)`),
				regexp.MustCompile(`(?i)\srar(a|e|o|i)`),
			},
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
				regexp.MustCompile(`(?i)\sapprezzabil(e|i)`),
			},
			color: color.New(color.FgGreen).SprintFunc(),
		},
		{
			key: "Impegno, competenza, managerialità",
			regexes: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\sdot[e|i] managerial(e|i)`),
				regexp.MustCompile(`(?i)\scapacit[à|a'] organizzativ(a|e)`),
				regexp.MustCompile(`(?i)\scompetenz(a|e)] tecnic[a|he]`),
			},
			color: color.New(color.FgBlue).SprintFunc(),
		},
	}

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
	for i := len(categories) - 1; i >= 0; i-- {
		fmt.Printf("%s: %d\n", categories[i].color(categories[i].key), categories[i].count)
	}
	fmt.Printf("----------------------------------------------------------------\n")
	fmt.Printf("%s", data)
}

func valueIn(value string, values ...string) bool {
	for _, v := range values {
		if value == v {
			return true
		}
	}
	return false
}
