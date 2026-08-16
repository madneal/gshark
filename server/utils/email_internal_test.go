package utils

import "testing"

func TestSplitRecipients(t *testing.T) {
	got := splitRecipients("first@example.com, second@example.com;third@example.com\nfourth@example.com")
	want := []string{"first@example.com", "second@example.com", "third@example.com", "fourth@example.com"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestSplitRecipientsIgnoresEmptyValues(t *testing.T) {
	got := splitRecipients(" ; first@example.com,, ")
	if len(got) != 1 || got[0] != "first@example.com" {
		t.Fatalf("got %v", got)
	}
}
