package analyser

import (
	"bytes"
	"encoding/base64"
	"image/color"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/tenebresus/djensen_server/pkg/db"
	"golang.org/x/image/bmp"
)

const TMPDIR string = "/tmp/djensen_server"
const SPLIT_DURATION int = 30

type Base struct {
    r uint32
    g uint32
    b uint32
}

func (ba Base) RGBA() (r uint32, g uint32, b uint32, a uint32) {
    return ba.r, ba.g, ba.b, 65535
}

func Run(audio_file string) {

    song_name := getSongName(audio_file)
    encoded_timestamps := getEncodedTimestamps(audio_file)
    audio_data := getAudioData(audio_file)
    
    db.Insert(song_name, encoded_timestamps, audio_data)

}

func getSongName(audio_file string) string {

    slash_splits := strings.Split(audio_file, "/")
    dot_splits := strings.Split(slash_splits[len(slash_splits) - 1], ".")
    return dot_splits[0]

}

func getAudioData(audio_file string) []byte {

    audio_data, err := os.ReadFile(audio_file)

    if err != nil {
        log.Fatal(err)
    }

    return audio_data

}

func getEncodedTimestamps(audio_file string) string {
    
    os.Mkdir(TMPDIR, 0777)
    defer os.RemoveAll(TMPDIR)

    timestamps := getTimestamps(audio_file)
    timestamps_string := strings.Join(timestamps, ",")

    return base64.StdEncoding.EncodeToString([]byte(timestamps_string))

}

func getTimestamps(audio_file string) []string {

    audio_length := getAudioFileLength(audio_file)
    audio_split_amount := math.Floor(float64(audio_length) / 30)

    var ret []string

    i := 0
    start_time := 0
    for i < int(audio_split_amount) {

        split_name := "split_" + strconv.Itoa(i) + ".wav"
        spectopic_name := "specto_" + strconv.Itoa(i) + ".bmp"

        createAudioSplit(start_time, audio_file, split_name)
        createSpectrumBMP(TMPDIR + "/" + split_name, spectopic_name)

        spectrum_timestamps := getTimestampsFromSpectrum(TMPDIR + "/" + spectopic_name)
        ret = append(ret, spectrum_timestamps...)

        i++
        start_time += SPLIT_DURATION

    }

    return ret

}

func getTimestampsFromSpectrum(spectopic_path string) []string {

    // R: 38036 G: 257 B: 23130
    base_colour := Base{r: 38036, g: 257, b: 23130}
    base_luminance := getLuminance(base_colour)

    file, err := os.ReadFile(spectopic_path)
    if err != nil {
        log.Fatal(err)
    }

    bmp_image, err := bmp.Decode(bytes.NewReader(file))
    if err != nil {
        log.Fatal(err)
    }

    var ret []string
    elapsed_ms := 0

    i := 0
    for i < bmp_image.Bounds().Max.X {

        // TODO: make sure the elapsed_ms makes sense; revisit later
        elapsed_ms += 100

        c := bmp_image.At(i, 0)
        luminance := getLuminance(c)

        if luminance > base_luminance {
            ret = append(ret, strconv.Itoa(elapsed_ms))
            elapsed_ms = 0
        }

        i++

    }

    return ret

} 

// ffmpeg -ss 30 -i input.wmv -c copy -t 10 output.wmv

func createAudioSplit(start_time int, audio_file string, split_name string) {

    audio_file_length_cmd := exec.Command("/opt/homebrew/bin/ffmpeg", "-i", audio_file, "-ss", strconv.Itoa(start_time), "-t", strconv.Itoa(SPLIT_DURATION), TMPDIR + "/" + split_name)

    err := audio_file_length_cmd.Run()
    if err != nil {
        log.Println("Something went wrong in splitting the audio!")
        log.Fatal(err)
    }

}

// ffmpeg -i output2.wav -lavfi "showspectrumpic=s=300x512:start=16875:stop=16876:legend=disabled" spectrogram.bmp
// Start freq = 16875 en freq = 16876
func createSpectrumBMP(split_file string, spectopic_name string) {

    audio_file_length_cmd := exec.Command("/opt/homebrew/bin/ffmpeg", "-i", split_file, "-lavfi", "showspectrumpic=s=300x512:start=16875:stop=16876:legend=disabled", TMPDIR +  "/" + spectopic_name)

    err := audio_file_length_cmd.Run()
    if err != nil {
        log.Println("Something went wrong in splitting the audio!")
        log.Fatal(err)
    }

}

func getAudioFileLength(audio_file string) int {

    audio_file_length_cmd := exec.Command("/opt/homebrew/bin/ffprobe", "-i", audio_file, "-show_entries", "format=duration", "-v", "quiet", "-of", "csv=p=0")

    output, err := audio_file_length_cmd.Output()
    if err != nil {
        log.Println("Something went wrong!")
        log.Fatal(err)
    }

    audio_length_string := string(output[:len(output)-1])
    audio_length := 0

    if strings.Contains(audio_length_string, ".") {

        split := strings.Split(audio_length_string, ".")
        audio_length, _ = strconv.Atoi(split[0])

    }

    return audio_length

}

// Luminance (Y) = 0.2126 * R' + 0.7152 * G' + 0.0722 * B'

// R' = ((R / 255) <= 0.03928) ? (R / 255) / 12.92 : ((R / 255 + 0.055) / 1.055) ^ 2.4
// G' = ((G / 255) <= 0.03928) ? (G / 255) / 12.92 : ((G / 255 + 0.055) / 1.055) ^ 2.4
// B' = ((B / 255) <= 0.03928) ? (B / 255) / 12.92 : ((B / 255 + 0.055) / 1.055) ^ 2.4

func getLuminance(c color.Color) float64 {

    r, g, b, _ := c.RGBA()

    r_linear := getLinearGrade(float64(r))
    g_linear := getLinearGrade(float64(g))
    b_linear := getLinearGrade(float64(b))

    return 0.2126 * r_linear + 0.7152 * g_linear + 0.0722 * b_linear

}

func getLinearGrade(grade float64) float64 {

    var ret float64 = 0
    if grade / 255 <= 0.03928 {
        ret = (grade / 255) / 12.92
    } else {
        ret = math.Pow((grade / 255 + 0.055) / 1.055, 2.4)
    }

    return ret

}
