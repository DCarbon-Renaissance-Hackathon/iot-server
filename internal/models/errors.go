package models

// type ECode int

// const (
// 	// Common error
// 	ECodeBadRequest        ECode = 40000
// 	ECodeUnauthorized      ECode = 40001
// 	ECodePermissionDenied  ECode = 40003
// 	ECodeNotExisted        ECode = 40004
// 	ECodeExisted           ECode = 40005
// 	ECodeQueryParamInvalid ECode = 40006
// 	ECodeInvalidSignature  ECode = 40007
// 	ECodeAddressIsEmpty    ECode = 40008

// 	// Project error

// 	// IOT error
// 	ECodeIOTNotAllowed      ECode = 41000
// 	ECodeIOTInvalidNonce    ECode = 41001
// 	ECodeIOTInvalidMintSign ECode = 41002

// 	// Sensor error
// 	ECodeSensorNotAllowed      ECode = 41100
// 	ECodeSensorInvalidNonce    ECode = 41101
// 	ECodeSensorInvalidMintSign ECode = 41102
// 	ECodeSensorInvalidMetric   ECode = 41103
// 	ECodeSensorInvalidType     ECode = 41104
// 	ECodeSensorHasNoAddress    ECode = 41105
// 	ECodeSensorHasAddress      ECode = 41106
// )

// const (
// 	ECodeInternal     = 50000
// 	ECodeNotImplement = 50001
// )

// var (
// 	ErrorUnauthorized     = NewError(ECodeUnauthorized, "")
// 	ErrorPermissionDenied = NewError(ECodePermissionDenied, "")
// )

// type Error struct {
// 	Code    ECode  `json:"code"`
// 	Message string `json:"message"`
// }

// func NewError(code ECode, msg string) error {
// 	var err = &Error{
// 		Code:    code,
// 		Message: msg,
// 	}
// 	return err
// }

// func (err *Error) Error() string {
// 	return fmt.Sprintf("Code:%d Message:%s", err.Code, err.Message)
// }

// func (err Error) String() string {
// 	return fmt.Sprintf("Code:%d Message:%s", err.Code, err.Message)
// }

// func ParsePostgresError(label string, err error) error {
// 	if nil == err {
// 		return nil
// 	}
// 	log.Println("Postgres error: ", err)
// 	if err == gorm.ErrRecordNotFound {
// 		return NewError(
// 			ECodeNotExisted,
// 			label+" is not existed",
// 		)
// 	}

// 	if strings.Contains(err.Error(), "duplicate") {
// 		return NewError(
// 			ECodeExisted,
// 			label+" is existed",
// 		)
// 	}
// 	return ErrInternal(err)
// }

// // ErrInternal :
// func ErrInternal(err error) error {
// 	if nil == err {
// 		return nil
// 	}
// 	log.Println("Internal error: ", err)
// 	return NewError(ECodeInternal, "internal error")
// }

// // ErrInternal :
// func ErrNotImplement() error {
// 	return NewError(ECodeNotImplement, "not implement")
// }

// // ErrInternal :
// func ErrBadRequest(msg string) error {
// 	// log.Println("Bad request error: ", err)
// 	return NewError(ECodeBadRequest, msg)
// }

// // ErrInternal :
// func ErrQueryParam(msg string) error {
// 	return NewError(ECodeQueryParamInvalid, msg)
// }
