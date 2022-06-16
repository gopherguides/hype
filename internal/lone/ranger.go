package lone

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
)

type Ranger struct {
	Start int
	End   int

	mu    sync.RWMutex
	once  sync.Once
	rx    *regexp.Regexp
	rxErr error
}

func (r *Ranger) String() string {
	if r == nil {
		return "0:0"
	}

	return fmt.Sprintf("%d:%d", r.Start, r.End)
}

// Unsigned returns an error if the range
// contains a negative number.
func (r *Ranger) Unsigned() error {
	if r == nil {
		return fmt.Errorf("ranger is nil")
	}

	if r.Start < 0 {
		return fmt.Errorf("start must be positive")
	}

	if r.End < 0 {
		return fmt.Errorf("end must be positive")
	}

	return nil
}

// Validate returns an error if
// start > end.
func (r *Ranger) Validate() error {
	if r == nil {
		return fmt.Errorf("ranger is nil")
	}

	if r.Start > r.End {
		return fmt.Errorf("start (%d) cannot be greater than end (%d)", r.Start, r.End)
	}

	return nil
}

func (r *Ranger) IsRange(s string) bool {
	rx, err := r.Regexp()
	if err != nil {
		return false
	}

	return rx.MatchString(s)
}

func (r *Ranger) Regexp() (*regexp.Regexp, error) {
	if r == nil {
		return nil, fmt.Errorf("ranger is nil")
	}

	r.once.Do(func() {
		r.mu.Lock()
		r.rx, r.rxErr = regexp.Compile(`([\+|\-]?\d*):([\+|\-]?\d*)`)
		r.mu.Unlock()
	})

	return r.rx, r.rxErr
}

func (r *Ranger) Parse(s string) error {
	rx, err := r.Regexp()
	if err != nil {
		return fmt.Errorf("failed to compile regexp: %q: %w", s, err)
	}

	res := rx.FindAllStringSubmatch(s, -1)
	if len(res) == 0 {
		return fmt.Errorf("invalid format: %q", s)
	}

	m := res[0]

	if len(m) < 3 {
		return fmt.Errorf("invalid format: %q", s)
	}

	if len(m[1]) > 0 {
		start, err := strconv.Atoi(m[1])
		if err != nil {
			return fmt.Errorf("failed to parse start: %q: %w", m[1], err)
		}
		r.Start = start
	}

	if len(m[2]) > 0 {
		end, err := strconv.Atoi(m[2])
		if err != nil {
			return fmt.Errorf("failed to parse end: %q: %w", m[2], err)
		}
		r.End = end
	}

	return nil
}

func Range(s string) (start int, end int, err error) {
	r := &Ranger{}
	err = r.Parse(s)
	return r.Start, r.End, err
}

func IsRange(s string) bool {
	r := &Ranger{}
	return r.IsRange(s)
}
