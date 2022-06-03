// repl/repl.go

package repl

import (
  "bufio"
  "fmt"
  "io"
  "intInGo/lexer"
  "intInGo/token"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
  scanner := bufio.NewScanner(in)

  for {
    fmt.Fprintf(out, PROMPT)

    // read from input source until hitting newline
    scanned := scanner.Scan()
    if !scanned {
      return
    }

    // instantiate lexer with just read line
    line := scanner.Text()
    l := lexer.New(line)

    // print all tokens from lexer until hitting EOF
    for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
      fmt.Fprintf(out, "%+v\n", tok)
    }
  }
}
