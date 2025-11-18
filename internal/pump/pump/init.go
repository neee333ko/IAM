package pump

var availablePumps map[string]Pump

func init() {
	availablePumps = map[string]Pump{
		"MongoDB Pump": &MongoPump{},
	}
}
