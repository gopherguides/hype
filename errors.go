package hype

import (
	"fmt"
	"strings"
)

const (
	ErrNilFigure = ErrIsNil("figure")
)

type ErrAttrNotFound string

func (e ErrAttrNotFound) Error() string {
	return fmt.Sprintf("attribute %q not found", string(e))
}

type ErrAttrEmpty string

func (e ErrAttrEmpty) Error() string {
	return fmt.Sprintf("attribute %q is empty", string(e))
}

type ErrIsNil string

func (e ErrIsNil) Error() string {
	return fmt.Sprintf("%s is nil", string(e))
}

type PostExecuteError struct {
	Err          error
	OrigErr      error
	PostExecuter PostExecuter
}

func (e PostExecuteError) Error() string {
	var errs []string

	if e.Err != nil {
		errs = append(errs, e.Err.Error())
	}

	if e.OrigErr != nil {
		errs = append(errs, e.OrigErr.Error())
	}

	return fmt.Sprintf("post execute error: [%T]: %v", e.PostExecuter, strings.Join(errs, "; "))
}

type PostParseError struct {
	Err        error
	OrigErr    error
	PostParser PostParser
}

func (e PostParseError) Error() string {
	var errs []string

	if e.Err != nil {
		errs = append(errs, e.Err.Error())
	}

	if e.OrigErr != nil {
		errs = append(errs, e.OrigErr.Error())
	}

	return fmt.Sprintf("post parse error: [%T]: %v", e.PostParser, strings.Join(errs, "; "))
}

type PreExecuteError struct {
	Err         error
	PreExecuter PreExecuter
}

func (e PreExecuteError) Error() string {
	return fmt.Sprintf("pre execute error: [%T]: %v", e.PreExecuter, e.Err)
}

type PreParseError struct {
	Err       error
	PreParser PreParser
}

func (e PreParseError) Error() string {
	return fmt.Sprintf("pre parse error: [%T]: %v", e.PreParser, e.Err)
}

func WrapNodeErr(n Node, err error) error {
	if err == nil {
		return nil
	}

	if h, ok := n.(Tag); ok {
		return fmt.Errorf("%T: %s: %w", h, h.StartTag(), err)
	}

	return fmt.Errorf("%T: %w", n, err)
}
