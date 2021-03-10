package audioplayer

import (
	"errors"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"
	"github.com/imroc/req"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const channelName = "com.kech.audioplayer/playaudio"

type AudioPlayer struct{}

var _ flutter.Plugin = &AudioPlayer{} // compile-time type check

var ctrl *beep.Ctrl

func (p *AudioPlayer) InitPlugin(messenger plugin.BinaryMessenger) error {
	channel := plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})
	channel.HandleFunc("playAudio", handlePlayAudio)
	return nil // no error
}

func handlePlayAudio(arguments interface{}) (reply interface{}, err error) {
	argsMap := arguments.(map[interface{}]interface{})
	url := argsMap["url"].(string)
	resp, err := playAudio(url)
	return resp, err
}

func playAudio(url string) (bool, error) {
	//const url1 = "https://web-mall-cdn.ams3.digitaloceanspaces.com/8180d349-9269-4958-8929-fb025041be90?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=Q6MAXU4EJUULDWY2OYJT%2F20210307%2Fams3%2Fs3%2Faws4_request&X-Amz-Date=20210307T181541Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Signature=92514ecbbff1cc1457946bf2777968ba2a566ccdc8ec099afa682b66a92248eb"
	//f, err := os.Open("gunshot.mp3")
	r, err := req.Get(url)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	if r.Response().StatusCode != 200 {
		return false, errors.New(string(r.Response().StatusCode))
	}

	dname, err := ioutil.TempDir("", "yscmessages")

	fname := filepath.Join(dname, "message.wav")

	err = r.ToFile(fname)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	f, err := os.Open(fname)

	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/5))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done

	//defer os.RemoveAll(dname)

	return false, nil
}
