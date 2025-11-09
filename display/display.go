package display

import (
	"log/slog"

	"github.com/mrechtien/mixgo/config"
)

type DisplayType int

const (
	BISTATE DisplayType = iota
	TEMPORARY
)

type DisplayEvent struct {
	Control uint8
	Value   uint8
}

type DisplayHolder struct {
	DisplayType DisplayType
	Display     any
}

type BiState interface {
	On()
	Off()
}

type Temporary interface {
	Trigger()
}

type Display interface {
	CreateBiState(mapping config.Mapping) BiState
	CreateTemporary(mapping config.Mapping) Temporary
}

type DisplayCreator func() Display

var displayRegistry = map[string]DisplayCreator{}
var displays = map[uint8]*DisplayHolder{}

func AddDisplay(inputType string, creator DisplayCreator) {
	displayRegistry[inputType] = creator
}

func createDisplay(display string) Display {
	return displayRegistry[display]()
}

func newDisplayHolder(displayType DisplayType, display any) *DisplayHolder {
	return &DisplayHolder{
		DisplayType: displayType,
		Display:     display,
	}
}

func CreateDisplayEvent(control uint8, value uint8) *DisplayEvent {
	return &DisplayEvent{
		Control: control,
		Value:   value,
	}
}

func SetupDisplay(display *config.Display, mappings []config.Mapping, displayEvents chan *DisplayEvent) {

	stateDisplay := createDisplay(display.Name)
	for _, mapping := range mappings {
		switch mapping.Name {
		case config.MUTE_GROUP, config.MUTE_CHANNEL:
			display := stateDisplay.CreateBiState(mapping)
			displays[mapping.Target] = newDisplayHolder(BISTATE, display)
		case config.TAP_DELAY:
			displays[mapping.Target] = newDisplayHolder(TEMPORARY, stateDisplay.CreateTemporary(mapping))
		}
	}

	go listenToDisplayEvents(displayEvents)
}

func listenToDisplayEvents(displayEvents chan *DisplayEvent) {
	for event := range displayEvents {
		displayHolder := displays[event.Control]
		if displayHolder == nil {
			slog.Warn("unmapped displayHolder", slog.Any("control", event.Control))
			continue
		}
		slog.Debug("mapping display event", slog.Any("type", displayHolder.DisplayType), slog.Any("event", event))
		switch displayHolder.DisplayType {
		case BISTATE:
			notifyBiStateDisplay(displayHolder.Display.(BiState), event)
		case TEMPORARY:
			notifyTemporaryDisplay(displayHolder.Display.(Temporary), event)
		}
	}
}

func notifyBiStateDisplay(biState BiState, event *DisplayEvent) {
	if event.Value > 0 {
		biState.On()
	} else {
		biState.Off()
	}
}

func notifyTemporaryDisplay(temporary Temporary, event *DisplayEvent) {
	temporary.Trigger()
}
