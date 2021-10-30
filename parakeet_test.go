package main

import (
	"strings"
	"testing"
)

const file = `
[2021-03-01T10:29:59.215Z] *** creekorful (~creekorfu@static.8.8.8.8.clients.example.org) joined
[2021-03-01T10:29:59.215Z] *** hentiphase (~hentiphase@static.4.4.4.4.clients.example.org) joined
[2021-03-01T11:59:42.847Z] *** creekorful changed nick to creekorful_
[2021-03-01T11:59:43.847Z] *** creekorful_ changed nick to creekorful
[2021-03-01T12:00:43.164Z] *** hentiphase (~hentiphase@static.4.4.4.4.clients.example.org) quit (Ping timeout: 480 seconds)
[2021-03-01T13:00:59.215Z] *** hentiphase (~hentiphase@static.4.4.4.4.clients.example.org) joined
[2021-03-01T13:25:17.928Z] <hentiphase> creekorful: when in doubt, take it to mail. -> unblock
[2021-03-01T13:37:14.974Z] <creekorful> thanks for your feedback :)
[2021-03-01T13:59:00.732Z] <hentiphase> creekorful: https://example.org/manual :)
[2021-03-01T20:38:26.113Z] *** creekorful (~creekorfu@static.8.8.8.8.clients.example.org) quit ()
[2021-03-01T21:24:49.669Z] *** hentiphase (~hentiphase@static.4.4.4.4.clients.example.org) quit (Remote host closed the connection)
[2021-03-23T10:29:59.215Z] *** creekorful (~creekorfu@static.8.8.8.8.clients.example.org) joined
[2021-03-23T10:46:21.525Z] * creekorful sent a long message:  < something >
`

func TestParseLog(t *testing.T) {
	ch, err := parseLog("#test-channel", strings.NewReader(file))
	if err != nil {
		t.Fatal(err)
	}

	if ch.Name != "#test-channel" {
		t.Fail()
	}
	if len(ch.Messages) != 3 {
		t.Fail()
	}

	msg := ch.Messages[0]
	if msg.Sender != "hentiphase" {
		t.Fail()
	}
	if msg.Content != "creekorful: when in doubt, take it to mail. -> unblock" {
		t.Fail()
	}

	msg = ch.Messages[1]
	if msg.Sender != "creekorful" {
		t.Fail()
	}
	if msg.Content != "thanks for your feedback :)" {
		t.Fail()
	}

	msg = ch.Messages[2]
	if msg.Sender != "hentiphase" {
		t.Fail()
	}
	if msg.Content != "creekorful: https://example.org/manual :)" {
		t.Fail()
	}
}

func TestGenerateHTML(t *testing.T) {
	w := &strings.Builder{}

	ch, err := parseLog("#test-channel", strings.NewReader(file))
	if err != nil {
		t.Fatal(err)
	}

	if err := generateHTML(ch, w); err != nil {
		t.Fatal(err)
	}

	val := w.String()

	if !strings.Contains(val, "<title>#test-channel</title>") {
		t.Fail()
	}
	if !strings.Contains(val, "<a href=\"https://example.org/manual\">https://example.org/manual</a>") {
		t.Fail()
	}
	if !strings.Contains(val, "; font-weight: bold;\">creekorful</span>&gt;") {
		t.Fail()
	}
	if !strings.Contains(val, "; font-weight: bold;\">hentiphase</span>&gt;") {
		t.Fail()
	}
}

func BenchmarkGenerateHTML(b *testing.B) {
	w := &strings.Builder{}

	for n := 0; n < b.N; n++ {
		ch, err := parseLog("#test-channel", strings.NewReader(file))
		if err != nil {
			b.Fatal(err)
		}

		if err := generateHTML(ch, w); err != nil {
			b.Fatal(err)
		}
	}
}

func TestTrimSuffixes(t *testing.T) {
	if val := trimSuffixes("#test-channel.txt", []string{".txt", ".log"}); val != "#test-channel" {
		t.Fail()
	}
	if val := trimSuffixes("#test-channel.log", []string{".txt", ".log"}); val != "#test-channel" {
		t.Fail()
	}
	if val := trimSuffixes("#test-channel.bin", []string{".txt", ".log"}); val != "#test-channel.bin" {
		t.Fail()
	}
}
