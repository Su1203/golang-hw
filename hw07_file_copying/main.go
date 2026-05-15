package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	// валидация
	if from == "" || to == "" {
		fmt.Println("Error: -from and -to are required")
		os.Exit(1)
	}

	fmt.Printf("Copying from %s to %s (offset: %d, limit: %d)\n", from, to, offset, limit)

	if err := Copy(from, to, offset, limit); err != nil {
		fmt.Printf("Ошибка при копировании: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Копирование успешно завершено.")
}
