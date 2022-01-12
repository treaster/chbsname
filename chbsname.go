package chbsname

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
	Generate(int) string
}

func NewBuilder(words io.Reader, minLength int, maxLength int) (Builder, error) {
	finalWords := []string{}
	scanner := bufio.NewScanner(words)
	for scanner.Scan() {
		word := scanner.Text()
		fmt.Println("consider", word)
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
			fmt.Println("discard", word)
			continue
		}
		finalWords = append(finalWords, word)
	}

	if len(finalWords) == 0 {
		return nil, errors.New("No valid words found in words input")
	}

	return builder{
		finalWords,
	}, nil
}

func NewBuilderFromFile(wordsFilePath string, minLength int, maxLength int) (Builder, error) {
	f, err := os.Open(wordsFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading words file %q: %s", wordsFilePath, err.Error())
	}
	defer func() { _ = f.Close() }()

	return NewBuilder(f, minLength, maxLength)
}

type builder struct {
	words []string
}

func (b builder) Generate(count int) string {
	components := make([]string, count)
	for i := 0; i < count; i++ {
		r := rand.Intn(len(b.words))
		components[i] = b.words[r]
	}

	return strings.Join(components, "-")
}
