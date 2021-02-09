/*******************************************************************************
 * This is not my code
 * It's bwmarrin's code
 *
 */

package youtube

import (
	"bufio"
	"encoding/binary"
	"io"
	"os/exec"
	"strconv"
	"sync"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

const (
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // max size of opus data
)

var (
	speakers    map[uint32]*gopus.Decoder
	opusEncoder *gopus.Encoder
	mu          sync.Mutex
)

func sendPCM(v *discordgo.VoiceConnection, pcm <-chan []int16) {
	if pcm == nil {
		return
	}

	var err error

	opusEncoder, err = gopus.NewEncoder(frameRate, channels, gopus.Audio)

	if err != nil {
		return
	}

	for {
		recv, ok := <-pcm
		if !ok {
			return
		}
		opus, err := opusEncoder.Encode(recv, frameSize, maxBytes)
		if err != nil {
			return
		}

		if v.Ready == false || v.OpusSend == nil {
			return
		}
		v.OpusSend <- opus
	}
}

func receivePCM(v *discordgo.VoiceConnection, c chan *discordgo.Packet) {
	if c == nil {
		return
	}

	var err error

	for {
		if v.Ready == false || v.OpusRecv == nil {
			return
		}

		p, ok := <-v.OpusRecv
		if !ok {
			return
		}

		if speakers == nil {
			speakers = make(map[uint32]*gopus.Decoder)
		}

		_, ok = speakers[p.SSRC]
		if !ok {
			speakers[p.SSRC], err = gopus.NewDecoder(48000, 2)
			if err != nil {
				continue
			}
		}

		p.PCM, err = speakers[p.SSRC].Decode(p.Opus, 960, false)
		if err != nil {
			continue
		}

		c <- p
	}
}

func playAudioFile(v *discordgo.VoiceConnection, filename string, stop <-chan bool) {

	// Create a shell command "object" to run.
	run := exec.Command("ffmpeg", "-i", filename, "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	ffmpegout, err := run.StdoutPipe()
	if err != nil {
		return
	}

	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)
	err = run.Start()
	if err != nil {
		return
	}

	go func() {
		<-stop
		err = run.Process.Kill()
	}()

	err = v.Speaking(true)
	if err != nil {
	}
	defer func() {
		err := v.Speaking(false)
		if err != nil {
		}
	}()

	send := make(chan []int16, 2)
	defer close(send)

	close := make(chan bool)
	go func() {
		sendPCM(v, send)
		close <- true
	}()

	for running == true {
		audiobuf := make([]int16, frameSize*channels)
		err = binary.Read(ffmpegbuf, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		if err != nil {
			return
		}

		select {
		case send <- audiobuf:
		case <-close:
			return
		}
	}
}
