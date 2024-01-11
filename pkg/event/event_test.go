package event

import "testing"

func TestScan(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		e := Event{Data: 1}
		var resp int
		if err := e.Scan(&resp); err != nil {
			t.Fatal(err)
		}

		t.Log("resp int:", resp)
	})

	t.Run("map", func(t *testing.T) {
		e := Event{Data: map[string]any{"name": "xxx", "age": 18}}
		var resp struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		if err := e.Scan(&resp); err != nil {
			t.Fatal(err)
		}

		t.Log("resp map:", resp)
	})

	t.Run("struct", func(t *testing.T) {
		type a struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		var resp a
		e := Event{Data: &a{Name: "xxx", Age: 19}}
		if err := e.Scan(&resp); err != nil {
			t.Fatal(err)
		}

		t.Log("resp map:", resp)
	})
}
