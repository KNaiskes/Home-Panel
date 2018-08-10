package devices


type LedStrip struct {
	DisplayName string
	Name	    string
	State	    string
	Color       string
	Topic	    string
}

type Lights struct {
	DisplayName string
	Name	    string
	State       string
	Topic       string
}

func GetLedstrips() []LedStrip {
	//Add any new led strips below
	bedroomLedstrip := LedStrip{"Bedroom", "bedroom_ledstrip", "false",
				   "white", "ledStrip"}

	MyledStrips := []LedStrip{bedroomLedstrip}

	return MyledStrips
}

func GetLights() []Lights {
	//Add any new light below
	officeLamp := Lights{"Office Lamp", "office_lamp", "true", "officeLamp"}
	DeskLamp := Lights{"Desk Lamp", "desk_lamp", "false", "deskLamp"}
	MyLights := []Lights{officeLamp, DeskLamp}

	return MyLights
}


