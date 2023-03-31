package models

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

const (
	// Common error
	ECodeBadRequest        = 40000
	ECodeUnauthorized      = 40001
	ECodePermissionDenied  = 40003
	ECodeNotExisted        = 40004
	ECodeExisted           = 40005
	ECodeQueryParamInvalid = 40006
	ECodeInvalidSignature  = 40007
	ECodeAddressIsEmpty    = 40008

	// Project error

	// IOT error
	ECodeIOTNotAllowed      = 41000
	ECodeIOTInvalidNonce    = 41001
	ECodeIOTInvalidMintSign = 41002

	// Sensor error
	ECodeSensorNotAllowed      = 41100
	ECodeSensorInvalidNonce    = 41101
	ECodeSensorInvalidMintSign = 41102
	ECodeSensorInvalidMetric   = 41103
	ECodeSensorHasNoAddress    = 41104
	ECodeSensorHasAddress      = 41105
)

const (
	ECodeInternal     = 50000
	ECodeNotImplement = 50001
)

var (
	ErrorUnauthorized     = NewError(ECodeUnauthorized, "")
	ErrorPermissionDenied = NewError(ECodePermissionDenied, "")
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, msg string) error {
	var err = &Error{
		Code:    code,
		Message: msg,
	}
	return err
}

func (err *Error) Error() string {
	return fmt.Sprintf("Code:%d Message:%s", err.Code, err.Message)
}

func (err Error) String() string {
	return fmt.Sprintf("Code:%d Message:%s", err.Code, err.Message)
}

func ParsePostgresError(label string, err error) error {
	if nil == err {
		return nil
	}
	log.Println("Postgres error: ", err)
	if err == gorm.ErrRecordNotFound {
		return NewError(
			ECodeNotExisted,
			label+" is not existed",
		)
	}

	if strings.Contains(err.Error(), "duplicate") {
		return NewError(
			ECodeExisted,
			label+" is existed",
		)
	}
	return ErrInternal(err)
}

// ErrInternal :
func ErrInternal(err error) error {
	if nil == err {
		return nil
	}
	log.Println("Internal error: ", err)
	return NewError(ECodeInternal, "internal error")
}

// ErrInternal :
func ErrNotImplement() error {
	return NewError(ECodeNotImplement, "not implement")
}

// ErrInternal :
func ErrBadRequest(msg string) error {
	// log.Println("Bad request error: ", err)
	return NewError(ECodeBadRequest, msg)
}

// ErrInternal :
func ErrQueryParam(msg string) error {
	return NewError(ECodeQueryParamInvalid, msg)
}
