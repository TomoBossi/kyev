package kyev

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	bufferSize     = 24 * 128
	inputEventSize = 24
	keyboardEV     = 0x120013
)

type InputEvent struct {
	Time  Timeval
	Type  uint16
	Code  uint16
	Value uint32
} // 24 bytes

type Timeval struct {
	Sec  uint64
	Usec uint64
}

type Keyboard struct {
	Name    string
	Phys    string
	devNode *os.File
	buffer  []byte
}

func newKeyboard(name, phys string, devNode *os.File) *Keyboard {
	return &Keyboard{
		Name:    name,
		Phys:    phys,
		devNode: devNode,
		buffer:  make([]byte, bufferSize),
	}
}

func discoverKeyboards(nameMatch, physMatch string) ([]*Keyboard, error) {
	var keyboards []*Keyboard

	devices, err := os.ReadFile("/proc/bus/input/devices")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`(?s).*N: Name="(.*)".*P: Phys=([^\n]*).*H: Handlers=.*event(\d+).*B: EV=([^\n]*)`)
	for dev := range strings.SplitSeq(string(devices), "\nI: ") {

		matches := re.FindStringSubmatch(dev)
		if matches == nil {
			continue
		}

		name := matches[1]
		phys := matches[2]
		devNodePath := fmt.Sprintf("/dev/input/event%s", matches[3])

		ev, err := strconv.ParseUint(matches[4], 16, 64)
		if err != nil {
			return nil, err
		}

		isKeyboard := ev&keyboardEV == keyboardEV
		nameMatched := strings.Contains(strings.ToLower(name), strings.ToLower(nameMatch))
		physMatched := strings.Contains(strings.ToLower(phys), strings.ToLower(physMatch))

		if isKeyboard && nameMatched && physMatched {

			devNode, err := os.Open(devNodePath) // requires sudo
			if err != nil {
				return nil, err
			}

			keyboard := newKeyboard(
				name,
				phys,
				devNode,
			)

			keyboards = append(keyboards, keyboard)
		}
	}

	return keyboards, nil
}

func GetKeyboard(nameMatch, physMatch string) (*Keyboard, error) {
	keyboards, err := discoverKeyboards(nameMatch, physMatch)
	if err != nil {
		return nil, err
	}

	if len(keyboards) > 0 {
		sort.Slice(keyboards, func(i, j int) bool {
			return len(keyboards[i].Name) < len(keyboards[j].Name)
		})
		return keyboards[0], nil
	}

	return nil, fmt.Errorf("no keyboards were discovered")
}

func (k *Keyboard) GetEvents() ([]InputEvent, error) {
	numBytes, err := k.devNode.Read(k.buffer)
	if err != nil {
		return nil, err
	}

	eventCount := numBytes / inputEventSize
	if eventCount == 0 {
		return nil, nil
	}

	events := make([]InputEvent, eventCount)
	buffer := bytes.NewReader(k.buffer[:numBytes])
	err = binary.Read(buffer, binary.LittleEndian, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (k *Keyboard) GetKeyPresses() ([]InputEvent, error) {
	events, err := k.GetEvents()
	if err != nil {
		return nil, err
	}

	var keyPresses []InputEvent
	for _, event := range events {
		if event.Type == EV_KEY && event.Value == BTN_PRESSED {
			keyPresses = append(keyPresses, event)
		}
	}

	return keyPresses, nil
}

func (k *Keyboard) Close() error {
	if k.devNode == nil {
		return nil
	}
	err := k.devNode.Close()
	k.devNode = nil
	return err
}
