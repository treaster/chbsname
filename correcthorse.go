package correcthorse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// A Builder will generate correcthorse-style strings with a number of words
// equal to the argument, separated by dashes.
// e.g. builder.Build(4) -> "correct-horse-battery-staple"
type Builder interface {
	Build(int) string
}

// Generate a builder from an array of words, ensuring all words are within a
// specified length range. Also ignore words/lines containing #-'" or
// whitespace. Aim for ASCII-only words
//
// Args:
//   rollFn: the randomizer function to use. Recommend passing rand.Intn or
//     similar. Be sure to seed it first.
//   words: the word list to draw from
//   minLength: ignore words in the words list shorter than minLength.
//   maxLength: ignore words in the words list longer than maxLength.
func NewBuilderFromStrings(rollFn func(int) int, words []string, minLength int, maxLength int) (Builder, error) {
	if rollFn == nil {
		return nil, fmt.Errorf("rollFn must be defined")
	}

	finalWords := []string{}
	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) == 0 {
			continue
		}
		if len(word) < minLength {
			continue
		}
		if len(word) > maxLength {
			continue
		}
		if strings.ContainsAny(word, "#- '\"") {
			continue
		}
		finalWords = append(finalWords, word)
	}

	if len(finalWords) == 0 {
		return nil, errors.New("No valid words found in words input")
	}

	return generator{
		rollFn,
		finalWords,
	}, nil
}

// Generate a builder from an io.Reader which emits newline-separated words.
// See NewBuilderFromStrings for a description of other arguments.
func NewBuilderFromReader(rollFn func(int) int, words io.Reader, minLength int, maxLength int) (Builder, error) {
	scanner := bufio.NewScanner(words)
	wordsArr := []string{}
	for scanner.Scan() {
		word := scanner.Text()
		wordsArr = append(wordsArr, word)
	}

	return NewBuilderFromStrings(rollFn, wordsArr, minLength, maxLength)
}

// Generate a builder from text file which contains newline-separated words.
// See NewBuilderFromStrings for a description of other arguments.
func NewBuilderFromFile(rollFn func(int) int, wordsFilePath string, minLength int, maxLength int) (Builder, error) {
	f, err := os.Open(wordsFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading words file %q: %s", wordsFilePath, err.Error())
	}
	defer func() { _ = f.Close() }()

	return NewBuilderFromReader(rollFn, f, minLength, maxLength)
}

type generator struct {
	rollFn func(int) int
	words  []string
}

func (g generator) Build(count int) string {
	components := make([]string, count)
	for i := 0; i < count; i++ {
		r := g.rollFn(len(g.words))
		components[i] = g.words[r]
	}

	return strings.Join(components, "-")
}
