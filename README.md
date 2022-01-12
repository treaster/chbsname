# chbsname
Generate strings like "correct-horse-battery-staple" (https://xkcd.com/936/), suitable for secure(?) passwords or other identification.

No words list is provided. A list of words suitable for the calling application must be provided. If you need a words list, consider the Linux built-in dictionary (/usr/share/dict/words), or a Scrabble word list like https://www.wordgamedictionary.com/twl06/download/twl06.txt. Licensing on these specific files is unknown. Use your judgment.

### Usage example with words list in a strings slice
```
package main

import (
    "fmt"
    "github.com/treaster/chbsname"
)

func main() {
    words := []string{
        "dog",
        "cat",
        "duck",
        "goat",
        "horse",
        "goose",
    }

    b, err := chbsname.NewBuilder(words, 3, 4)
    if err != nil {
        return err
    }

    fmt.Println(b.Generate(3)) // Possible output = "goat-duck-duck"
}
```

### Usage example with words list in an io.Reader
```
package main

import (
    "fmt"
    "github.com/treaster/chbsname"
)

func main() {
    words := bytes.NewReader([]byte(`dog
        cat
        duck
        goat
        horse
        goose`))

    b, err := chbsname.NewBuilderFromReader(words, 3, 4)
    if err != nil {
        return err
    }

    fmt.Println(b.Generate(3)) // Possible output = "goat-duck-duck"
}
```

### Usage example with words list in a file
```
package main

import (
    "fmt"
    "github.com/treaster/chbsname"
)

func main() {
    b, err := chbsname.NewBuilderFromFile("/path/to/wordlist.txt", 3, 4)
    if err != nil {
        return err
    }

    fmt.Println(b.Generate(3)) // Possible output = "goat-duck-duck"
}
```
