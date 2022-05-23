package assignment04

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func Test_Venue_Logs(t *testing.T) {
	t.Parallel()

	bb := &bytes.Buffer{}

	v := &Venue{Log: bb}

	err := v.Entertain(100, &Band{MinAudience: 50})

	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}

	act := strings.TrimSpace(bb.String())
	exp := `The Thunders of Rock has completed setup.
The Thunders of Rock has performed for 100 people.
The Thunders of Rock has completed teardown.`

	if act != exp {
		t.Fatalf("expected %q, got %q", exp, act)
	}

}

func Test_Venue_Entertain(t *testing.T) {
	t.Parallel()

	const allgood = `The Thunders of Rock has completed setup.
The Thunders of Rock has performed for 42 people.
The Thunders of Rock has completed teardown.
Maybelle Marie has performed for 42 people.`

	table := []struct {
		acts []Entertainer
		aud  int
		err  bool
		exp  string
		name string
	}{
		{
			name: "setup error",
			acts: []Entertainer{
				&Band{IsSetup: true},
			},
			exp: "we already setup our gear",
			err: true,
		},
		{
			name: "teardown error",
			acts: []Entertainer{
				&Band{IsTorndown: true},
			},
			exp: "we already tore down our gear",
			err: true,
		},
		{
			name: "play error",
			acts: []Entertainer{
				Poet{},
			},
			exp: "i'm not playing for just the bartender",
			err: true,
			aud: 1,
		},
		{
			name: "no acts",
			err:  true,
			exp:  "there are no entertainers to perform",
		},
		{
			name: "all good",
			acts: []Entertainer{
				&Band{},
				Poet{},
			},
			aud: 42,
			exp: allgood,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			bb := &bytes.Buffer{}

			v := &Venue{Log: bb}

			err := v.Entertain(tt.aud, tt.acts...)

			if tt.err {
				if err == nil {
					log.Fatalf("expected error, got none")
				}

				act := err.Error()
				if !strings.Contains(act, tt.exp) {
					log.Fatalf("expected %q to contain %q", act, tt.exp)
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %s", err)
			}

			act := bb.String()

			if !strings.Contains(act, tt.exp) {
				log.Fatalf("expected %q to contain %q", act, tt.exp)
			}

		})
	}
}
