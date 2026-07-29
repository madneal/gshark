package service

import "testing"

func TestValidateFilterClass(t *testing.T) {
	t.Parallel()

	for _, class := range []string{"extension", "keyword"} {
		if err := validateFilterClass(class); err != nil {
			t.Fatalf("expected %q to be supported: %v", class, err)
		}
	}

	if err := validateFilterClass("sec_keyword"); err == nil {
		t.Fatal("expected removed secondary keyword filter class to be rejected")
	}
}
