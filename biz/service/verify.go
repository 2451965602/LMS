package service

import (
	"github.com/2451965602/LMS/pkg/errno"
	"regexp"
	"strconv"
)

func CheckMobile(phone string) bool {
	// 匹配规则
	// ^1第一位为一
	// [345789]{1} 后接一位345789 的数字
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1[345789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)

}

func isValidISBN10(isbn string) bool {
	sum := 0
	for i, char := range isbn {
		var digit int
		if i == 9 && char == 'X' {
			digit = 10
		} else {
			digit, _ = strconv.Atoi(string(char))
			if digit < 0 || digit > 9 {
				return false
			}
		}
		sum += digit * (10 - i)
	}
	return sum%11 == 0
}

func isValidISBN13(isbn string) bool {
	sum := 0
	for i, char := range isbn {
		digit, _ := strconv.Atoi(string(char))
		if digit < 0 || digit > 9 {
			return false
		}
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}
	return sum%10 == 0
}

func RegisterCheck(username, phone string) error {
	if len(username) != 9 {
		return errno.Errorf(errno.ServiceInvalidUsername, "username must be between 3 and 20 characters") // 用户名长度不符合要求，返回错误
	}
	if !CheckMobile(phone) {
		return errno.Errorf(errno.ServiceInvalidPhone, "invalid phone number") // 手机号码格式不正确，返回错误
	}
	return nil
}

func IsValidISBN(isbn string) bool {
	length := len(isbn)

	switch length {
	case 10:
		return isValidISBN10(isbn)
	case 13:
		return isValidISBN13(isbn)
	default:
		return false
	}
}

func CheckAuthor(author string) bool {
	regRuler := `^(?:\[\p{Han}+\]\s+)?[\p{Han}\p{Latin}\s·]+$`
	reg := regexp.MustCompile(regRuler)

	return reg.MatchString(author)
}
