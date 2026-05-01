package main

import (
	"flag"
	"fmt"
	"os"
)

// 이 예제는 표준 라이브러리 `flag`만으로 서브커맨드를 구성하는 패턴입니다.
// 예:
//
//	go run ./flag/advanced/subcommands serve -port=8080
//	go run ./flag/advanced/subcommands migrate -dry-run
func main() {
	if len(os.Args) < 2 {
		printTopUsage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "serve":
		runServe(os.Args[2:])
	case "migrate":
		runMigrate(os.Args[2:])
	case "-h", "--help", "help":
		printTopUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printTopUsage()
		os.Exit(2)
	}
}

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	port := fs.Int("port", 8080, "listen port")
	env := fs.String("env", "dev", "environment name")

	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	// fs.Args() 는 "serve" 커맨드의 남은 positional args 입니다.
	if extra := fs.Args(); len(extra) > 0 {
		fmt.Fprintf(os.Stderr, "unexpected args: %v\n\n", extra)
		fs.Usage()
		os.Exit(2)
	}

	fmt.Printf("serve: env=%s port=%d\n", *env, *port)
}

func runMigrate(args []string) {
	fs := flag.NewFlagSet("migrate", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	dryRun := fs.Bool("dry-run", false, "print actions only")
	steps := fs.Int("steps", 0, "limit number of steps, 0 means all")

	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	fmt.Printf("migrate: dry_run=%v steps=%d\n", *dryRun, *steps)
	if remaining := fs.Args(); len(remaining) > 0 {
		fmt.Printf("migrate args: %v\n", remaining)
	}
}

func printTopUsage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  mycli <command> [flags]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  serve    run server")
	fmt.Fprintln(os.Stderr, "  migrate  run migrations")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "examples:")
	fmt.Fprintln(os.Stderr, "  go run ./flag/advanced/subcommands serve -port=8080")
	fmt.Fprintln(os.Stderr, "  go run ./flag/advanced/subcommands migrate -dry-run")
}
