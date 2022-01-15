package correcthorse_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/treaster/correcthorse"
)

type fakeRoller struct {
	values []int
	next   int
}

func (r *fakeRoller) Roll(n int) int {
	val := r.values[r.next]
	r.next++
	return val % n
}

func TestCorrectHorse_Simple(t *testing.T) {

	words := bytes.NewReader([]byte(`dog
cat
duck
goat
horse
goose`))

	{
		words.Seek(0, io.SeekStart)
		roller := fakeRoller{values: []int{9, 7}}
		b, err := correcthorse.NewBuilderFromReader(roller.Roll, words, 3, 5)
		require.NoError(t, err)
		output := b.Build(0)
		require.Equal(t, "", output)
	}

	{
		words.Seek(0, io.SeekStart)
		roller := fakeRoller{values: []int{9, 7}}
		b, err := correcthorse.NewBuilderFromReader(roller.Roll, words, 3, 5)
		require.NoError(t, err)
		output := b.Build(1)
		require.Equal(t, "goat", output)
	}

	{
		words.Seek(0, io.SeekStart)
		roller := fakeRoller{values: []int{9, 7}}
		b, err := correcthorse.NewBuilderFromReader(roller.Roll, words, 3, 5)
		require.NoError(t, err)
		output := b.Build(2)
		require.Equal(t, "goat-cat", output)
	}
}

func TestCorrectHorse_LengthLimits(t *testing.T) {
	words := bytes.NewReader([]byte(`dog
cat
duck
goat
horse   
goose   
`))

	{
		words.Seek(0, io.SeekStart)
		roller := fakeRoller{values: []int{9, 7}}
		b, err := correcthorse.NewBuilderFromReader(roller.Roll, words, 2, 4)
		require.NoError(t, err)
		output := b.Build(2)
		require.Equal(t, "cat-goat", output)
	}

	{
		words.Seek(0, io.SeekStart)
		roller := fakeRoller{values: []int{9, 7}}
		b, err := correcthorse.NewBuilderFromReader(roller.Roll, words, 4, 6)
		require.NoError(t, err)
		output := b.Build(2)
		require.Equal(t, "goat-goose", output)
	}

	{
		words.Seek(0, io.SeekStart)
		roller := fakeRoller{values: []int{9, 7}}
		b, err := correcthorse.NewBuilderFromReader(roller.Roll, words, 4, 4)
		require.NoError(t, err)
		output := b.Build(2)
		require.Equal(t, "goat-goat", output)
	}
}

func TestCorrectHorse_Repeats(t *testing.T) {
	words := bytes.NewReader([]byte(`dog`))

	{
		words.Seek(0, io.SeekStart)
		roller := fakeRoller{values: []int{9, 7, 5, 3}}
		b, err := correcthorse.NewBuilderFromReader(roller.Roll, words, 2, 4)
		require.NoError(t, err)
		output := b.Build(4)
		require.Equal(t, "dog-dog-dog-dog", output)
	}
}

func TestCorrectHorse_NoisyWordList(t *testing.T) {
	words := bytes.NewReader([]byte(`
illegalquote'
illegalquote"
#illegalhashcanuseascomment
illegal space
illegal-dash
		legalleadingspace
legaltrailingspace    

legalwordafterblankline
`))

	{
		roller := fakeRoller{values: []int{0, 1, 2}}
		b, err := correcthorse.NewBuilderFromReader(roller.Roll, words, 0, 100)
		require.NoError(t, err)
		output := b.Build(3)
		require.Equal(t, "legalleadingspace-legaltrailingspace-legalwordafterblankline", output)
	}
}
