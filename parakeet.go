package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mvdan.cc/xurls/v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed "channel.gohtml"
var tplStr string

// Channel represent an IRC channel
type Channel struct {
	Name     string
	Messages []Message
}

// Message represent an IRC message
type Message struct {
	Time    time.Time
	Sender  string
	Content string
}

// Context is the application context
type Context struct {
	Colors []string
	Users  map[string]string
}

func main() {
	inputFileFlag := flag.String("input", "", "IRC log input file")
	outputFileFlag := flag.String("output", "", "Where the HTML should be outputted")

	flag.Parse()

	inputFile, err := os.Open(*inputFileFlag)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	ch, err := parseLog(trimSuffixes(filepath.Base(*inputFileFlag), []string{".txt", ".log"}), inputFile)
	if err != nil {
		panic(err)
	}
	log.Printf("Successfully parsed %d messages from %s", len(ch.Messages), ch.Name)

	outputFile, err := os.Create(*outputFileFlag)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	if err := generateHTML(ch, outputFile); err != nil {
		panic(err)
	}
	log.Printf("Sucessfully generated HTML template at: %s", *outputFileFlag)
}

func generateHTML(ch *Channel, writer io.Writer) error {
	ctx := Context{
		Colors: []string{
			"red", "green", "blue", "violet", "turquoise",
			"coral", "brown", "crimson", "darkblue",
			"fuschia", "indigo", "maroon", "navy",
		},
		Users: map[string]string{},
	}

	tpl, err := template.New("channel").
		Funcs(map[string]interface{}{
			"colorUsername": ctx.colorUsername,
		}).
		Parse(tplStr)
	if err != nil {
		return err
	}

	return tpl.Execute(writer, ch)
}

func parseLog(name string, reader io.Reader) (*Channel, error) {
	ch := &Channel{
		Name:     name,
		Messages: []Message{},
	}

	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Replace all URLs in one pass
	content := applyURLs(string(b))

	for _, line := range strings.Split(content, "\n") {
		// Only keep 'message' line i.e which contains something like '] <username>'
		if !strings.Contains(line, "] <") || !strings.Contains(line, ">") {
			continue
		}

		// Approximate line parsing
		date := line[1:strings.Index(line, "] <")]
		username := line[strings.Index(line, "<")+1 : strings.Index(line, ">")]
		content := line[strings.Index(line, "> ")+2:]

		t, err := time.Parse(time.RFC3339, date)
		if err != nil {
			break
		}

		ch.Messages = append(ch.Messages, Message{
			Time:    t,
			Sender:  username,
			Content: content,
		})
	}

	return ch, nil
}

func trimSuffixes(s string, suffixes []string) string {
	for _, suffix := range suffixes {
		s = strings.TrimSuffix(s, suffix)
	}

	return s
}

func applyURLs(s string) string {
	rxStrict := xurls.Strict()
	return rxStrict.ReplaceAllStringFunc(s, func(s string) string {
		return fmt.Sprintf("<a href=\"%s\">%s</a>", s, s)
	})
}

func (c *Context) colorUsername(s string) template.HTML {
	// check if username color is not yet applied
	color, exist := c.Users[s]
	if !exist {
		// pick-up new random color for the username
		color = c.Colors[rand.Intn(len(c.Colors)-1)]
		c.Users[s] = color
	}

	return template.HTML(fmt.Sprintf("&lt;<span style=\"color: %s; font-weight: bold;\">%s</span>&gt;", color, s))
}
