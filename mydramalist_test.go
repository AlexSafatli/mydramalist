package mydramalist

import (
	"testing"
)

func TestClient_Search(t *testing.T) {
	s := NewClient()
	dramas, err := s.Search("Eternal Love")
	if err != nil {
		t.Errorf("Error after searching: %+v", err)
	}
	if len(dramas) == 0 {
		t.Error("Dramas was empty slice")
	}
	if dramas[0].Country != "China" {
		t.Errorf("Drama did not have country China, had %s", dramas[0].Country)
	}
	if !stringsContains(dramas[0].Genres, "Martial Arts") {
		t.Errorf("%+v did not have Martial Arts", dramas[0].Genres)
	}
}

func stringsContains(a []string, target string) bool {
	for _, s := range a {
		if s == target {
			return true
		}
	}
	return false
}
