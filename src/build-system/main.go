package main

import "os"

func help() {
    println("Usage: go run main.go <command>")
    println("help - help")
    println("build - build and zip all lambdas")
    println("ffmpeg - build and zip ffmpeg transocder")
    println("<lambda-name> - build and zip a specific lambda")
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
    } else if (command == "ffmpeg") {
    } else {
        // build a specific lambda
    }
}
