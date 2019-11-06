package mydramalist

import "testing"

func TestClient_Search(t *testing.T) {
	s := NewClient()
	dramas, err := s.Search("Eternal Love")
	if err != nil {
		t.Errorf("Error after searching: %+v", err)
	}
	if len(dramas) == 0 {
		t.Error("Dramas returned was empty slice")
	}
}
