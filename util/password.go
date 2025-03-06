package util

import (
	"errors"
	"fmt"
	"github.com/alexedwards/argon2id"
	"strings"
	"unicode"
)

var (
	ErrPasswordTooShort   = errors.New("密码至少需要8个字符")
	ErrMissingUpperCase   = errors.New("密码必须包含至少一个大写字母")
	ErrMissingLowerCase   = errors.New("密码必须包含至少一个小写字母")
	ErrMissingNumber      = errors.New("密码必须包含至少一个数字")
	ErrMissingSpecialChar = errors.New("密码必须包含至少一个特殊字符(!@#$%^&*)")
)

func ValidatePasswordPolicy(pwd string) error {
	if len(pwd) < 8 {
		return ErrPasswordTooShort
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, c := range pwd {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*", c):
			hasSpecial = true
		}
	}

	var errs []error
	if !hasUpper {
		errs = append(errs, ErrMissingUpperCase)
	}
	if !hasLower {
		errs = append(errs, ErrMissingLowerCase)
	}
	if !hasNumber {
		errs = append(errs, ErrMissingNumber)
	}
	if !hasSpecial {
		errs = append(errs, ErrMissingSpecialChar)
	}

	if len(errs) > 0 {
		return fmt.Errorf("密码不符合策略: %v", errs)
	}
	return nil
}

/*var (
	upperRegex   = regexp.MustCompile(`[A-Z]`)
	lowerRegex   = regexp.MustCompile(`[a-z]`)
	numberRegex  = regexp.MustCompile(`[0-9]`)
	specialRegex = regexp.MustCompile(`[!@#$%^&*]`)
)

func ValidatePasswordPolicy(pwd string) error {
	if len(pwd) < 8 {
		return errors.New("密码至少需要8个字符")
	}

	if !(upperRegex.MatchString(pwd) && lowerRegex.MatchString(pwd) && numberRegex.MatchString(pwd) && specialRegex.MatchString(pwd)) {
		return errors.New("密码必须包含大小写字幕、数组和特殊字符")
	}
	return nil
}*/

// 注册时应生成符合PHC标准的哈希字符串
func HashPassword(password string) (string, error) {
	params := &argon2id.Params{
		Memory:      64 * 1024,
		Iterations:  3,  //迭代次数
		Parallelism: 4,  //并行度
		SaltLength:  16, //盐值长度
		KeyLength:   32, //哈希长度
	}
	hash, err := argon2id.CreateHash(password, params)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func VerifyPassword(inputPassword, storedPassword string) error {
	match, _, err := argon2id.CheckHash(inputPassword, storedPassword)
	if err != nil {
		return fmt.Errorf("密码验证失败: hash=%s error=%v", storedPassword, err)
	}
	if !match {
		return fmt.Errorf("密码不匹配: identifier=%s", inputPassword)
	}
	return nil
}
