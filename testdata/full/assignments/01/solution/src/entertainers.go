package assignment04

import "fmt"

var _ Entertainer = &Band{}
var _ Setuper = &Band{}
var _ Teardowner = &Band{}

type Band struct {
	IsSetup     bool
	IsTorndown  bool
	MinAudience int
	PlayedFor   int
}

func (b Band) Name() string {
	return "The Thunders of Rock"
}

func (b Band) Validate(v Venue) error {
	if v.Audience < b.MinAudience {
		return fmt.Errorf("we don't play small gigs")
	}
	return nil
}

func (b *Band) Perform(v Venue) error {
	if err := b.Validate(v); err != nil {
		return err
	}

	b.PlayedFor = v.Audience

	return nil
}

func (b *Band) Setup(v Venue) error {
	if b.IsSetup {
		return fmt.Errorf("we already setup our gear")
	}

	if err := b.Validate(v); err != nil {
		return err
	}
	b.IsSetup = true
	return nil
}

func (b *Band) Teardown(v Venue) error {
	if b.IsTorndown {
		return fmt.Errorf("we already tore down our gear")
	}

	b.IsTorndown = true
	return nil
}

var _ Entertainer = Poet{}

type Poet struct{}

func (p Poet) Name() string {
	return "Maybelle Marie"
}

func (p Poet) Perform(v Venue) error {
	if v.Audience == 1 {
		return fmt.Errorf("i'm not playing for just the bartender")
	}
	return nil
}
