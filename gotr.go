package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/sparrovv/gotr/googletranslate"
)

const translateToPath string = "/tmp/gotr.speech.file.to.mpg"

var usage string = `google translate in terminal

Usage:
  gotr <from> <to> <phrase>
  gotr [-s] <from> <to> <phrase>
  gotr (-h | --help)
  gotr (-l | --list)

Options:
  -h, --help     Show help
  -l, --list     List available languages
  -s, --speech   Enable speech synthesis
`

func main() {
	args, _ := docopt.Parse(usage, nil, true, "Gotr", false)

	if args["--list"].(bool) == true {
		fmt.Println(`Supported languages:`)
		fmt.Println(googletranslate.ListLanguages())
		os.Exit(1)
	}

	from := args["<from>"].(string)
	to := args["<to>"].(string)
	term := args["<phrase>"].(string)

	if len(from) == 0 || len(to) == 0 || len(term) == 0 {
		fmt.Println("  Usage: " + usage)
		os.Exit(1)
	}

	phrase, err := googletranslate.Translate(from, to, term)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(phrase.Translation)
	fmt.Println(strings.Join(phrase.ExtraMeanings, ", "))

	if args["--speech"].(bool) == true {
		err := googletranslate.FetchSoundFile(to, phrase.Translation, translateToPath)
		if err == nil {
			playSound(translateToPath)
		} else {
			fmt.Println(err)
		}
	}
}

var audioPlayers []Player = []Player{
	Player{Name: "afplay"},
	Player{Name: "mpg123"},
}

func playSound(path string) {
	for _, player := range audioPlayers {
		_, err := player.Play(path)
		if err == nil {
			return
		}
	}

	fmt.Println("Can't find compatible audio player")
	os.Exit(1)
}
