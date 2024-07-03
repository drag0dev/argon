package main

import (
	"fmt"
	"os"
	"strings"
)

func help() {
    println("Usage: go run main.go <command>")
    println("help - help")
    println("build - build and zip all lambdas")
    println("ffmpeg - build and zip ffmpeg transocder")
    println("<lambda-name> - build and zip a specific lambda")
}

func getAllLambdas() []string {
    dir, err := os.Open("../")
    if (err != nil) {
        fmt.Printf("Cannot open src directory: %v\n", err)
        os.Exit(-1)
    }
    defer dir.Close()

    entries, err := dir.Readdir(0)
    if (err != nil) {
        fmt.Printf("Error reading src directory contents: %v\n", err)
        os.Exit(-1)
    }

    var lambdas []string

    for _, entry := range entries {
        fileName := entry.Name()
        if (strings.HasPrefix(fileName, "lambda-")) {
            lambdas = append(lambdas, fileName)
        }
    }
    return lambdas
}

func main() {
    if (len(os.Args) < 2) {
        help()
        return
    }

    command := os.Args[1]
    if (command == "help") {
        help()
    } else if (command == "build") {
        lambdas := getAllLambdas()
    } else if (command == "ffmpeg") {
    } else {
        // build a specific lambda
    }
}
