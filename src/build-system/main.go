package main

import (
	"fmt"
	"os"
	"os/exec"
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

func buildLambda(lambdaName string) {
    buildCommand := exec.Command("go", "build", "-o", fmt.Sprintf("../%s/bootstrap", lambdaName), "main.go")

    buildCommand.Env = os.Environ()
    buildCommand.Env = append(buildCommand.Env, "GOOS=linux")
    buildCommand.Env = append(buildCommand.Env, "GOARCH=amd64")

    _, err := buildCommand.CombinedOutput()
    if (err != nil) {
        fmt.Printf("Error running build command: %v\n", err)
        os.Exit(-1)
    }
}

func zipLambda(lambdaName string) {
    zipPath := fmt.Sprintf("../%s/function.zip", lambdaName)
    binaryPath := fmt.Sprintf("../%s/bootstrap", lambdaName)
    zipCommand := exec.Command("zip", zipPath, binaryPath)
    _, err := zipCommand.CombinedOutput()
    if (err != nil) {
        fmt.Printf("Error running zip command: %v\n", err)
        os.Exit(-1)
    }
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
        for _, lambda := range lambdas {
            fmt.Printf("Build %s ", lambda)
            buildLambda(lambda)
            fmt.Printf("\033[32mDONE\033[0m\n")

            fmt.Printf("Zip %s ", lambda)
            zipLambda(lambda)
            fmt.Printf("\033[32mDONE\033[0m\n")
        }
    } else if (command == "ffmpeg") {
    } else {
        // build a specific lambda
    }
}
