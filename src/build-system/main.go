package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func help() {
    println("Usage: go run main.go <command>")
    println("help - help")
    println("build - build and zip all lambdas")
    println("ffmpeg - prepare ffmpeg zip")
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
    buildCommand := exec.Command("go", "build", "-o", "bootstrap", "main.go")

    buildCommand.Dir = fmt.Sprintf("../%s/", lambdaName)

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
    zipCommand := exec.Command("zip", "-FS", "function.zip", "bootstrap")
    zipCommand.Dir = fmt.Sprintf("../%s/", lambdaName)

    _, err := zipCommand.CombinedOutput()
    if (err != nil) {
        fmt.Printf("Error running zip command: %v\n", err)
        os.Exit(-1)
    }
}

func lambdaFolderExist(lambdaName string) {
    _, err := os.Stat(fmt.Sprintf("../%s", lambdaName));
    if os.IsNotExist(err) {
        fmt.Printf("Lambda %s does not exist.\n", lambdaName)
        os.Exit(-1)
    }

    if (err != nil) {
        fmt.Printf("Error checking if lambda %s exists: %v\n", lambdaName, err)
        os.Exit(-1)
    }
}

func downloadFFMPEG() {
    url := "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"

    outFile, err := os.Create("./ffmpeg-binary.tar.xz")
    if (err != nil) {
        fmt.Printf("Error creating ffmpeg-binary.zip file: %v\n", err)
        os.Exit(-1)
    }
    defer outFile.Close()

    response, err := http.Get(url)
    if (err != nil) {
        fmt.Printf("Error downloading ffmpeg binary: %v\n", err)
        os.Exit(-1)
    }
    defer response.Body.Close()

    if (response.StatusCode != http.StatusOK) {
        fmt.Printf("Server returned: %s", response.Status)
        os.Exit(-1)
    }

    _, err = io.Copy(outFile, response.Body)
    if err != nil {
        fmt.Printf("Error writing binary zip: %v\n", err)
        os.Exit(-1)
    }
}

func prepFFMPEG() {
    untarCommand := exec.Command("tar", "-xf", "ffmpeg-binary.tar.xz")
    _, err := untarCommand.CombinedOutput()
    if (err != nil) {
        fmt.Printf("Error untarring ffmpeg: %v\n", err)
        os.Exit(-1)
    }

    // find the output folder
    dir, err := os.Open(".")
    if (err != nil) {
        fmt.Printf("Cannot open current directory: %v\n", err)
        os.Exit(-1)
    }
    defer dir.Close()

    entries, err := dir.Readdir(0)
    if (err != nil) {
        fmt.Printf("Error reading src directory contents: %v\n", err)
        os.Exit(-1)
    }

    untarringOutput := ""
    for _, entry := range entries {
        if (strings.HasSuffix(entry.Name(), "-amd64-static")) {
            untarringOutput = entry.Name()
            break
        }
    }

    if (untarringOutput == "") {
        fmt.Println("Cant find untarring output folder")
        os.Exit(-1)
    }

    // create layer structure
    err = os.MkdirAll(fmt.Sprintf("./%s/temp/bin", untarringOutput), 0755)
    if (err != nil) {
        fmt.Printf("Error creating dir structure for the zip: %v\n", err)
        os.Exit(-1)
    }

    // move the binary into the proper folder
    err = os.Rename(fmt.Sprintf("./%s/ffmpeg", untarringOutput), fmt.Sprintf("./%s/temp/bin/ffmpeg", untarringOutput))
    if (err != nil) {
        fmt.Printf("Error moving ffmpeg binary into the structure: %v\n", err)
        os.Exit(-1)
    }

    // zip the structure
    zipCommand := exec.Command("zip", "-r", "-X", "../ffmpeg.zip", ".")
    zipCommand.Dir = fmt.Sprintf("./%s/temp", untarringOutput)

    _, err = zipCommand.CombinedOutput()
    if (err != nil) {
        fmt.Printf("Error zipping ffmpeg binary: %v\n", err)
        os.Exit(-1)
    }

    // move the zip into the lambda
    err = os.Rename(fmt.Sprintf("./%s/ffmpeg.zip", untarringOutput), "../lambda-transcoder/ffmpeg.zip")
    if (err != nil) {
        fmt.Printf("Error moving ffmpeg.zip into the lambda: %v\n", err)
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
        fmt.Print("Downloading FFMPEG, \033[38;5;214mMIGHT TAKE A WHILE\033[0m ")
        downloadFFMPEG()
        fmt.Printf("\033[32mDONE\033[0m\n")

        fmt.Print("Making ffmpeg.zip ")
        prepFFMPEG()
        fmt.Printf("\033[32mDONE\033[0m\n")
    } else {
        lambda := os.Args[1]
        lambdaFolderExist(lambda)
        fmt.Printf("Build %s ", lambda)
        buildLambda(lambda)
        fmt.Printf("\033[32mDONE\033[0m\n")

        fmt.Printf("Zip %s ", lambda)
        zipLambda(lambda)
        fmt.Printf("\033[32mDONE\033[0m\n")
    }
}
