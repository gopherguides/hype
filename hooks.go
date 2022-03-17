package hype

import "fmt"

type PreParseHook func(p *Parser) error

type PreParseHooks []PreParseHook

type PostParseHook func(p *Parser, err error) error

type PostParseHooks []PostParseHook

type HookError struct {
	ParseErr error
	HookErr  error
}

func (h HookError) Unwrap() error {
	if h.ParseErr != nil {
		return h.ParseErr
	}
	return h.HookErr
}

func (h HookError) Error() string {
	if h.ParseErr != nil {
		return h.ParseErr.Error()
	}

	return fmt.Sprintf("%s: %s", h.ParseErr, h.HookErr)
}

func (hs PreParseHooks) Run(p *Parser) error {
	for _, h := range hs {
		if err := h(p); err != nil {
			return err
		}
	}

	return nil
}

func (hs PostParseHooks) Run(p *Parser, perr error) error {
	for _, h := range hs {
		if err := h(p, perr); err != nil {
			return HookError{
				ParseErr: perr,
				HookErr:  err,
			}
		}
	}

	return nil
}
