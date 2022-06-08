package audio

import (
	"bytes"
	"encoding/binary"
	"sync"

	"github.com/disgoorg/disgo/audio/opus"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

func NewPCMOpusReceiver(decoderCreateFunc func() *opus.Decoder, pcmFrameReceiver PCMFrameReceiver) voice.OpusFrameReceiver {
	if decoderCreateFunc == nil {
		decoderCreateFunc = func() *opus.Decoder {
			decoder, _ := opus.NewDecoder(48000, 2)
			return decoder
		}
	}
	return &pcmOpusReceiver{
		decoderCreateFunc: decoderCreateFunc,
		decoders:          map[snowflake.ID]*opus.Decoder{},
		pcmFrameReceiver:  pcmFrameReceiver,
		pcmBuff:           make([]int16, 960),
		opusBuff:          make([]byte, 4000),
	}
}

type pcmOpusReceiver struct {
	decoderCreateFunc func() *opus.Decoder
	decoders          map[snowflake.ID]*opus.Decoder
	decodersMu        sync.Mutex
	pcmFrameReceiver  PCMFrameReceiver
	pcmBuff           []int16
	opusBuff          []byte
}

func (r *pcmOpusReceiver) ReceiveOpusFrame(userID snowflake.ID, packet *voice.Packet) {
	r.decodersMu.Lock()
	decoder, ok := r.decoders[userID]
	if !ok {
		decoder = r.decoderCreateFunc()
		r.decoders[userID] = decoder
	}
	r.decodersMu.Unlock()

	n, err := decoder.Decode(r.opusBuff, r.pcmBuff, true)
	if err != nil {
		return
	}

	reader := bytes.NewReader(r.opusBuff[:n])
	if err = binary.Read(reader, binary.LittleEndian, r.pcmBuff); err != nil {
		return
	}

	r.pcmFrameReceiver.ReceivePCMFrame(userID, &PCMPacket{
		SSRC:      packet.SSRC,
		Sequence:  packet.Sequence,
		Timestamp: packet.Timestamp,
		PCM:       r.pcmBuff,
	})
}

func (r *pcmOpusReceiver) CleanupUser(userID snowflake.ID) {
	r.decodersMu.Lock()
	defer r.decodersMu.Unlock()
	decoder, ok := r.decoders[userID]
	if ok {
		decoder.Destroy()
		delete(r.decoders, userID)
	}
	r.pcmFrameReceiver.CleanupUser(userID)
}

func (r *pcmOpusReceiver) Close() {
	r.decodersMu.Lock()
	defer r.decodersMu.Unlock()
	for _, decoder := range r.decoders {
		decoder.Destroy()
	}
	r.pcmFrameReceiver.Close()
}

type PCMPacket struct {
	SSRC      uint32
	Sequence  uint16
	Timestamp uint32
	PCM       []int16
}