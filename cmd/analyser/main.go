package main

import (
	"os"

	"github.com/tenebresus/djensen_server/pkg/analyser"
)

func main() {

    // Get input arguments
    args := os.Args
    audio_file := args[1]

    // call run with the path of the audio file to analyse
    analyser.Run(audio_file)

}
