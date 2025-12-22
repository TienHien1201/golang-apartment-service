package xhttp

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	validate *validator.Validate

	errorMessageFuncs = map[string]errorMessageFunc{
		"required": simpleErrorMessage("%s is required"),
		"email":    simpleErrorMessage("%s must be a valid email address"),
		"min":      minMaxErrorMessage("at least"),
		"max":      minMaxErrorMessage("at most"),
		"oneof":    oneofErrorMessage("%s must be one of: %s"),
		"datetime": simpleErrorMessage("%s must be a valid datetime format: %s"),
		"uuid":     simpleErrorMessage("%s must be a valid UUID"),
		"url":      simpleErrorMessage("%s must be a valid URL"),
		"numeric":  simpleErrorMessage("%s must be a valid numeric value"),
		"alphanum": simpleErrorMessage("%s must contain only alphanumeric characters"),
		"phone":    simpleErrorMessage("%s must be a valid phone number"),
		"gt":       simpleErrorMessage("%s must be greater than %s"),
		"gte":      simpleErrorMessage("%s must be greater than or equal to %s"),
		"lt":       simpleErrorMessage("%s must be less than %s"),
		"lte":      simpleErrorMessage("%s must be less than or equal to %s"),
		"date":     simpleErrorMessage("%s must be a valid date format: YYYY-MM-DD"),
		"filetype": simpleErrorMessage("%s must be a valid file type. Allowed types: jpg, jpeg, png, gif, pdf, doc, docx"),
	}
	allowedImageExt = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}

	allowedImageMime = map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
)

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("phone", validatePhone)
	_ = validate.RegisterValidation("email", validateEmail)
	_ = validate.RegisterValidation("name", validateName)
	_ = validate.RegisterValidation("unit", validateUnit)
	_ = validate.RegisterValidation("date", validateDate)
	_ = validate.RegisterValidation("filetype", validateFileType)
}

// ReadAndValidateRequest reads the request and validates the struct.
func ReadAndValidateRequest(ctx echo.Context, req interface{}) interface{} {
	if err := ctx.Bind(req); err != nil {
		return validatorDefaultRules(err)
	}

	if err := defaults.Set(req); err != nil {
		return validatorDefaultRules(err)
	}

	if err := validate.StructCtx(ctx.Request().Context(), req); err != nil {
		return validatorDefaultRules(err)
	}

	return nil
}

func ValidateStruct(ctx context.Context, req interface{}) interface{} {
	if err := defaults.Set(req); err != nil {
		return validatorDefaultRules(err)
	}

	if err := validate.StructCtx(ctx, req); err != nil {
		return validatorDefaultRules(err)
	}

	return nil
}

func validatorDefaultRules(err error) interface{} {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errs := make([]ValidationError, 0, len(validationErrors))
		for _, e := range validationErrors {
			code := "ERR_" + strings.ToUpper(e.Tag())
			errs = append(errs, ValidationError{
				Code:    code,
				Field:   e.Field(),
				Message: getErrorMessage(e),
			})
		}
		return errs
	}

	var he *echo.HTTPError
	if errors.As(err, &he) {
		return []ValidationError{{
			Code:    "ERR_UNKNOWN",
			Message: he.Message.(string),
		}}
	}

	return []ValidationError{{
		Code:    "ERR_UNKNOWN",
		Message: err.Error(),
	}}
}

type errorMessageFunc func(fe validator.FieldError) string

func simpleErrorMessage(format string) errorMessageFunc {
	return func(fe validator.FieldError) string {
		if strings.Count(format, "%s") == 1 {
			return fmt.Sprintf(format, fe.Field())
		}
		return fmt.Sprintf(format, fe.Field(), fe.Param())
	}
}

func oneofErrorMessage(format string) errorMessageFunc {
	return func(fe validator.FieldError) string {
		return fmt.Sprintf(format, fe.Field(), strings.ReplaceAll(fe.Param(), " ", ", "))
	}
}

func minMaxErrorMessage(comparison string) errorMessageFunc {
	return func(fe validator.FieldError) string {
		if fe.Type().Kind() == reflect.String {
			return fmt.Sprintf("%s must be %s %s characters long", fe.Field(), comparison, fe.Param())
		}
		return fmt.Sprintf("%s must be %s %s", fe.Field(), comparison, fe.Param())
	}
}

func getErrorMessage(fe validator.FieldError) string {
	if fn, ok := errorMessageFuncs[fe.Tag()]; ok {
		return fn(fe)
	}
	return fmt.Sprintf("%s failed validation: %s", fe.Field(), fe.Tag())
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	phoneRegex := `^0\d{9,10}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}

func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if email == "" {
		return false
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return emailRegex.MatchString(email)
}

func validateName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	if name == "" || len(strings.TrimSpace(name)) < 2 {
		return false
	}
	return true
}

func validateUnit(fl validator.FieldLevel) bool {
	unit := fl.Field().String()
	validUnits := map[string]bool{
		"vnd": true,
		"lct": true,
		"ltt": true,
	}
	return validUnits[unit]
}

func validateDate(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	if date == "" {
		return true // allow empty date
	}

	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func validateFileType(fl validator.FieldLevel) bool {
	// Check if the field is a pointer to multipart.FileHeader
	if fl.Field().Kind() != reflect.Ptr {
		return true // Skip validation if not a pointer
	}

	if fl.Field().IsNil() {
		return true // Skip validation if nil (let "required" tag handle this)
	}

	// Get the actual value
	val := fl.Field().Elem()
	if !val.IsValid() {
		return false
	}

	// Check if it's a multipart.FileHeader
	fileHeader, ok := val.Interface().(*multipart.FileHeader)
	if !ok {
		return false
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		return false
	}

	// Remove dot from extension for comparison
	ext = strings.TrimPrefix(ext, ".")

	// List of allowed extensions
	allowedExts := map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"gif":  true,
		"pdf":  true,
		"doc":  true,
		"docx": true,
	}

	// Check if extension is allowed
	if !allowedExts[ext] {
		return false
	}

	// Get Content-Type from the header
	contentType := fileHeader.Header.Get("Content-Type")

	// Map of valid content types for each extension
	validMimeTypes := map[string][]string{
		"jpg":  {"image/jpeg"},
		"jpeg": {"image/jpeg"},
		"png":  {"image/png"},
		"gif":  {"image/gif"},
		"pdf":  {"application/pdf"},
		"doc":  {"application/msword"},
		"docx": {"application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
	}

	// For some content types like images, browsers often set generic MIME types
	// If content type is empty, we'll rely on extension validation only
	if contentType == "" {
		return true
	}

	// If we have a content type, validate it against the expected types for this extension
	validTypes, exists := validMimeTypes[ext]
	if !exists {
		return false
	}

	// Check if content type is in the list of valid types for this extension
	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}

	// For binary/octet-stream or application/octet-stream, trust the extension
	if contentType == "application/octet-stream" || contentType == "binary/octet-stream" {
		return true
	}

	// Postman and some browsers might use incorrect MIME types
	// For less strict validation, consider returning true here

	return false
}

// validateFileField validates a single file field
func validateFileField(field reflect.Value, fieldType reflect.StructField, multipartForm *MultipartForm) error {
	if fieldType.Type != reflect.TypeOf((*multipart.FileHeader)(nil)) {
		return nil
	}

	formKey := fieldType.Tag.Get("form")
	if formKey == "" {
		formKey = strings.ToLower(fieldType.Name)
	}

	if strings.Contains(fieldType.Tag.Get("validate"), "required") {
		file, exists := multipartForm.GetFile(formKey)
		if !exists {
			return fmt.Errorf("%s is required", fieldType.Name)
		}
		field.Set(reflect.ValueOf(file))
	} else if file, exists := multipartForm.GetFile(formKey); exists {
		field.Set(reflect.ValueOf(file))
	}

	return nil
}

// validateFileFields validates all file fields in the request struct
func validateFileFields(req interface{}, multipartForm *MultipartForm) error {
	val := reflect.ValueOf(req).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		if err := validateFileField(val.Field(i), typ.Field(i), multipartForm); err != nil {
			return err
		}
	}
	return nil
}

// ReadAndValidateFormRequest reads and validates a multipart form request
func ReadAndValidateFormRequest(c echo.Context, req interface{}, maxSize int64) (interface{}, error) {
	// Parse multipart form
	multipartForm, err := NewMultipartForm(c.Request(), maxSize)
	if err != nil {
		return nil, err
	}

	// Bind form data to request struct
	if err := c.Bind(req); err != nil {
		return nil, err
	}

	// Set default values
	if err := defaults.Set(req); err != nil {
		return nil, err
	}

	// Validate struct
	if err := validate.StructCtx(c.Request().Context(), req); err != nil {
		return validatorDefaultRules(err), nil
	}

	// Validate file fields
	if err := validateFileFields(req, multipartForm); err != nil {
		return []ValidationError{{
			Code:    "ERR_REQUIRED",
			Field:   "file",
			Message: err.Error(),
		}}, err
	}

	return nil, nil
}

func ValidateImageFile(file *multipart.FileHeader, maxSize int64) error {
	if file == nil {
		return NewAppError(
			"ERR_INVALID_IMAGE",
			"avatar",
			"Invalid image file",
			http.StatusBadRequest,
		)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExt[ext] {
		return NewAppError(
			"ERR_INVALID_IMAGE_TYPE",
			"avatar",
			"Only jpg, jpeg, png, gif files are allowed",
			http.StatusBadRequest,
		)
	}

	contentType := file.Header.Get("Content-Type")
	if contentType != "" && !allowedImageMime[contentType] {
		return NewAppError(
			"ERR_INVALID_IMAGE_MIME",
			"avatar",
			"Invalid image content type",
			http.StatusBadRequest,
		)
	}

	if file.Size > maxSize {
		return NewAppError(
			"ERR_IMAGE_TOO_LARGE",
			"avatar",
			"Image size exceeds limit",
			http.StatusBadRequest,
		)
	}

	return nil
}
