package connections

import (
	"errors"
	"testing"
)

func TestGetConnection(t *testing.T) {
	c := Connections{
		"c1": {
			Host:     "localhost",
			Port:     3306,
			Dbname:   "test_db",
			Username: "usr",
			Password: "pass",
		},
		"c2": {
			Host:     "12.34.56.78",
			Port:     3306,
			Dbname:   "test_db",
			Username: "usr",
			Password: "pass",
		},
	}

	c1, _ := c.GetConnection("c1")

	if c1.Host != "localhost" {
		t.Errorf("wanted %q, got %q", "localhost", c1.Host)
	}

	c2, _ := c.GetConnection("c2")

	if c2.Host != "12.34.56.78" {
		t.Errorf("wanted %q, got %q", "12.34.56.78", c2.Host)
	}

	_, err := c.GetConnection("c3")
	expectedErr := errors.New("invalid connection id: 'c3'\n")
	if err.Error() != expectedErr.Error() {
		t.Errorf("wanted error %q, got %q", expectedErr, err)
	}
}
