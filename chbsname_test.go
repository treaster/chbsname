package chbsname_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/treaster/chbsname"
)

func TestChbs_Simple(t *testing.T) {
	words := bytes.NewReader([]byte(`dog
cat
duck
goat
horse
goose`))

	{
		words.Seek(0, io.SeekStart)
		b, err := chbsname.NewBuilder(words, 3, 5)
		require.NoError(t, err)
		output := b.Generate(0)
		require.Equal(t, "", output)
	}

	{
		words.Seek(0, io.SeekStart)
		b, err := chbsname.NewBuilder(words, 3, 5)
		require.NoError(t, err)
		output := b.Generate(1)
		regex := "^(dog|cat|duck|goat|horse|goose)"
		require.Regexp(t, regex, output)
	}

	{
		words.Seek(0, io.SeekStart)
		b, err := chbsname.NewBuilder(words, 3, 5)
		require.NoError(t, err)
		output := b.Generate(2)
		regex := "^(dog|cat|duck|goat|horse|goose)-(dog|cat|duck|goat|horse|goose)$"
		require.Regexp(t, regex, output)
	}
}

func TestChbs_LengthLimits(t *testing.T) {
	words := bytes.NewReader([]byte(`dog
cat
duck
goat
horse   
goose   
`))

	{
		words.Seek(0, io.SeekStart)
		b, err := chbsname.NewBuilder(words, 2, 4)
		require.NoError(t, err)
		output := b.Generate(2)
		regex := "^(dog|cat|duck|goat)-(dog|cat|duck|goat)$"
		require.Regexp(t, regex, output)
	}

	{
		words.Seek(0, io.SeekStart)
		b, err := chbsname.NewBuilder(words, 4, 6)
		require.NoError(t, err)
		output := b.Generate(2)
		regex := "^(duck|goat|horse|goose)-(duck|goat|horse|goose)$"
		require.Regexp(t, regex, output)
	}

	{
		words.Seek(0, io.SeekStart)
		b, err := chbsname.NewBuilder(words, 4, 4)
		require.NoError(t, err)
		output := b.Generate(2)
		regex := "^(duck|goat)-(duck|goat)$"
		require.Regexp(t, regex, output)
	}
}

func TestChbs_Repeats(t *testing.T) {
	words := bytes.NewReader([]byte(`dog`))

	{
		words.Seek(0, io.SeekStart)
		b, err := chbsname.NewBuilder(words, 2, 4)
		require.NoError(t, err)
		output := b.Generate(4)
		require.Equal(t, "dog-dog-dog-dog", output)
	}
}

func TestChbs_NoisyWordList(t *testing.T) {
	words := bytes.NewReader([]byte(`
illegalquote'
illegalquote"
illegal space
illegal-dash
		legalleadingspace
legaltrailingspace    

legalwordafterblankline
`))

	{
		b, err := chbsname.NewBuilder(words, 0, 100)
		require.NoError(t, err)
		output := b.Generate(1)
		regex := "^(legalleadingspace|legaltrailingspace|legalwordafterblankline)$"
		require.Regexp(t, regex, output)
	}
}
