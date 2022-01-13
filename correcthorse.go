package correcthorse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

type Builder interface {
	Build(int) string
}

func NewBuilderFromStrings(words []string, minLength int, maxLength int) (Builder, error) {
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
		if strings.ContainsAny(word, "- '\"") {
			continue
		}
		finalWords = append(finalWords, word)
	}

	if len(finalWords) == 0 {
		return nil, errors.New("No valid words found in words input")
	}

	return generator{
		finalWords,
	}, nil
}

func NewBuilderFromReader(words io.Reader, minLength int, maxLength int) (Builder, error) {
	scanner := bufio.NewScanner(words)
	wordsArr := []string{}
	for scanner.Scan() {
		word := scanner.Text()
		wordsArr = append(wordsArr, word)
	}

	return NewBuilderFromStrings(wordsArr, minLength, maxLength)
}

func NewBuilderFromFile(wordsFilePath string, minLength int, maxLength int) (Builder, error) {
	f, err := os.Open(wordsFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading words file %q: %s", wordsFilePath, err.Error())
	}
	defer func() { _ = f.Close() }()

	return NewBuilderFromReader(f, minLength, maxLength)
}

type generator struct {
	words []string
}

func (g generator) Build(count int) string {
	components := make([]string, count)
	for i := 0; i < count; i++ {
		r := rand.Intn(len(g.words))
		components[i] = g.words[r]
	}

	return strings.Join(components, "-")
}
