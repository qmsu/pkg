package musicxml

import "testing"

func TestDecode(t *testing.T) {
	musicxmlFile := "1.musicxml"
	res, err := Decode(musicxmlFile)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(res)
}
