package protocol

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Validator 请求参数验证器
type Validator struct{}

// NewValidator 创建验证器实例
func NewValidator() *Validator {
	return &Validator{}
}

// Validate 验证结构体
func (v *Validator) Validate(data interface{}) error {
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return fmt.Errorf("validation requires a struct")
	}

	return v.validateStruct(value)
}

// validateStruct 验证结构体字段
func (v *Validator) validateStruct(value reflect.Value) error {
	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := valueType.Field(i)

		// 跳过嵌入字段（如BaseRequest）
		if fieldType.Anonymous {
			if field.Kind() == reflect.Struct {
				if err := v.validateStruct(field); err != nil {
					return err
				}
			}
			continue
		}

		// 获取validate标签
		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		// 验证字段
		if err := v.validateField(field, fieldType.Name, validateTag); err != nil {
			return err
		}
	}

	return nil
}

// validateField 验证单个字段
func (v *Validator) validateField(field reflect.Value, fieldName, validateTag string) error {
	rules := strings.Split(validateTag, ",")

	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		if rule == "" {
			continue
		}

		if err := v.applyRule(field, fieldName, rule); err != nil {
			return err
		}
	}

	return nil
}

// applyRule 应用验证规则
func (v *Validator) applyRule(field reflect.Value, fieldName, rule string) error {
	switch {
	case rule == "required":
		return v.validateRequired(field, fieldName)
	case strings.HasPrefix(rule, "min="):
		return v.validateMin(field, fieldName, rule)
	case strings.HasPrefix(rule, "max="):
		return v.validateMax(field, fieldName, rule)
	case strings.HasPrefix(rule, "len="):
		return v.validateLen(field, fieldName, rule)
	case rule == "email":
		return v.validateEmail(field, fieldName)
	case rule == "alphanum":
		return v.validateAlphanum(field, fieldName)
	default:
		// 忽略未知规则
		return nil
	}
}

// validateRequired 验证必填
func (v *Validator) validateRequired(field reflect.Value, fieldName string) error {
	switch field.Kind() {
	case reflect.String:
		if field.String() == "" {
			return fmt.Errorf("%s不能为空", fieldName)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() == 0 {
			return fmt.Errorf("%s不能为0", fieldName)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() == 0 {
			return fmt.Errorf("%s不能为0", fieldName)
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() == 0 {
			return fmt.Errorf("%s不能为0", fieldName)
		}
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map:
		if field.IsNil() {
			return fmt.Errorf("%s不能为空", fieldName)
		}
	}
	return nil
}

// validateMin 验证最小值/长度
func (v *Validator) validateMin(field reflect.Value, fieldName, rule string) error {
	minStr := strings.TrimPrefix(rule, "min=")
	min, err := strconv.Atoi(minStr)
	if err != nil {
		return fmt.Errorf("invalid min rule: %s", rule)
	}

	switch field.Kind() {
	case reflect.String:
		if len(field.String()) < min {
			return fmt.Errorf("%s长度不能少于%d个字符", fieldName, min)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() < int64(min) {
			return fmt.Errorf("%s不能小于%d", fieldName, min)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() < uint64(min) {
			return fmt.Errorf("%s不能小于%d", fieldName, min)
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() < float64(min) {
			return fmt.Errorf("%s不能小于%d", fieldName, min)
		}
	case reflect.Slice, reflect.Map:
		if field.Len() < min {
			return fmt.Errorf("%s长度不能少于%d", fieldName, min)
		}
	}
	return nil
}

// validateMax 验证最大值/长度
func (v *Validator) validateMax(field reflect.Value, fieldName, rule string) error {
	maxStr := strings.TrimPrefix(rule, "max=")
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		return fmt.Errorf("invalid max rule: %s", rule)
	}

	switch field.Kind() {
	case reflect.String:
		if len(field.String()) > max {
			return fmt.Errorf("%s长度不能超过%d个字符", fieldName, max)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() > int64(max) {
			return fmt.Errorf("%s不能大于%d", fieldName, max)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() > uint64(max) {
			return fmt.Errorf("%s不能大于%d", fieldName, max)
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() > float64(max) {
			return fmt.Errorf("%s不能大于%d", fieldName, max)
		}
	case reflect.Slice, reflect.Map:
		if field.Len() > max {
			return fmt.Errorf("%s长度不能超过%d", fieldName, max)
		}
	}
	return nil
}

// validateLen 验证固定长度
func (v *Validator) validateLen(field reflect.Value, fieldName, rule string) error {
	lenStr := strings.TrimPrefix(rule, "len=")
	length, err := strconv.Atoi(lenStr)
	if err != nil {
		return fmt.Errorf("invalid len rule: %s", rule)
	}

	switch field.Kind() {
	case reflect.String:
		if len(field.String()) != length {
			return fmt.Errorf("%s长度必须为%d个字符", fieldName, length)
		}
	case reflect.Slice, reflect.Map:
		if field.Len() != length {
			return fmt.Errorf("%s长度必须为%d", fieldName, length)
		}
	}
	return nil
}

// validateEmail 验证邮箱格式
func (v *Validator) validateEmail(field reflect.Value, fieldName string) error {
	if field.Kind() != reflect.String {
		return nil
	}

	email := field.String()
	if email == "" {
		return nil // 空值跳过，由required规则处理
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("%s格式不正确", fieldName)
	}

	return nil
}

// validateAlphanum 验证字母数字
func (v *Validator) validateAlphanum(field reflect.Value, fieldName string) error {
	if field.Kind() != reflect.String {
		return nil
	}

	value := field.String()
	if value == "" {
		return nil // 空值跳过，由required规则处理
	}

	alphanumRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphanumRegex.MatchString(value) {
		return fmt.Errorf("%s只能包含字母和数字", fieldName)
	}

	return nil
}

// ValidateRequest 验证请求参数的便捷方法
func ValidateRequest(req interface{}) error {
	validator := NewValidator()
	return validator.Validate(req)
}