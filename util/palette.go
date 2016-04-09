package util

var testPallete = []string{
	`#400000`,
	`#604000`,
	`#004000`,
	`#004060`,
	`#000060`,
	`#600060`,
}
var testPalletePoint int


func GetBackground() string {
	testPalletePoint++
	if testPalletePoint >= len(testPallete) {
		testPalletePoint=0
	}
	return testPallete[testPalletePoint]
}
