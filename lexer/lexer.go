package lexer

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

const (
	ORIGIN    = "ORIGIN"
	SCALE     = "SCALE"
	ROT       = "ROT"
	IS        = "IS"
	TO        = "TO"
	STEP      = "STEP"
	DRAW      = "DRAW"
	FOR       = "FOR"
	FROM      = "FROM"
	T         = "T"
	SEMICO    = ";"
	L_BRACKET = "("
	R_BRACKET = ")"
	COMMA     = ","
	PLUS      = "+"
	MINUS     = "-"
	MUL       = "*"
	DIV       = "/"
	POWER     = "**"
	FUNC      = "FUNC"
	CONST_ID  = "CONST_ID"
	NONTOKEN  = "NONTOKEN"
	ERRTOKEN  = "ERRTOKEN"
)

type Token struct {
	TokenType string
	Lexeme    string
	Value     float64
	FuncPtr   func(x float64) float64
}

var TokenTypeDict = map[string]Token{
	"PI":     Token{CONST_ID, "PI", math.Pi, nil},
	"E":      Token{CONST_ID, "E", math.E, nil},
	"T":      Token{T, "T", 0.0, nil},
	"ORIGIN": Token{ORIGIN, "ORIGIN", 0.0, nil},
	"SCALE":  Token{SCALE, "SCALE", 0.0, nil},
	"ROT":    Token{ROT, "ROT", 0.0, nil},
	"IS":     Token{IS, "IS", 0.0, nil},
	"FOR":    Token{FOR, "FOR", 0.0, nil},
	"FROM":   Token{FROM, "FROM", 0.0, nil},
	"TO":     Token{TO, "TO", 0.0, nil},
	"STEP":   Token{STEP, "STEP", 0.0, nil},
	"DRAW":   Token{DRAW, "DRAW", 0.0, nil},
	"SIN":    Token{FUNC, "SIN", 0.0, math.Sin},
	"COS":    Token{FUNC, "COS", 0.0, math.Cos},
	"TAN":    Token{FUNC, "TAN", 0.0, math.Tan},
	"LN":     Token{FUNC, "LN", 0.0, math.Log},
	"EXP":    Token{FUNC, "EXP", 0.0, math.Exp},
	"SQRT":   Token{FUNC, "SQRT", 0.0, math.Sqrt},
}

func getChar(str string, i int) string {
	if i < len(str) {
		return string(str[i])
	}
	return ""

}

func isAlpha(char string) bool {
	if char != "" && unicode.IsLetter([]rune(char)[0]) {
		return true
	}
	return false
}

func isDigit(char string) bool {
	if char != "" && unicode.IsDigit([]rune(char)[0]) {
		return true
	}
	return false
}

// Lexer ： 词法分析
func Lexer(str string) []Token {
	str = strings.ToUpper(str)

	var tokens []Token
	lineNum := 1
	i := 0
	for {
		char := getChar(str, i)
		if char == "" {
			tokens = append(tokens, Token{NONTOKEN, "", 0.0, nil})
			break
		} else if char == "\n" {
			lineNum++
			i++
			continue
		} else if char == " " || char == "\t" || char == "\r" {
			i++
			continue
		}

		var token Token

		if isAlpha(char) { // 解析关键词
			temp := char
			for {
				i++
				char = getChar(str, i)
				if isAlpha(char) {
					temp += char
				} else {
					i--
					break
				}
			}

			// 不是存在于map中的关键词，就返回错误
			if value, ok := TokenTypeDict[temp]; ok {
				token = value
			} else {
				token = Token{ERRTOKEN, "ERR", 0.0, nil}
				panic(fmt.Sprintf("Lexer error: Line%d unexpected token '%s'.", lineNum, temp))
			}

		} else if isDigit(char) { // 解析常量数字
			temp := char
			for {
				i++
				char = getChar(str, i)
				if isDigit(char) {
					temp += char
				} else if char == "." {
					temp += "."
					for {
						i++
						char = getChar(str, i)
						if isDigit(char) {
							temp += char
						} else {
							i--
							break
						}
					}
				} else {
					i--
					break
				}
			}
			num, _ := strconv.ParseFloat(temp, 64)
			token = Token{CONST_ID, temp, num, nil}

		} else {
			switch char {
			case ";":
				token = Token{SEMICO, char, 0.0, nil}
			case "(":
				token = Token{L_BRACKET, char, 0.0, nil}
			case ")":
				token = Token{R_BRACKET, char, 0.0, nil}
			case ",":
				token = Token{COMMA, char, 0.0, nil}
			case "+":
				token = Token{PLUS, char, 0.0, nil}
			case "*":
				i++
				char = getChar(str, i)
				if char == "*" {
					token = Token{POWER, "**", 0.0, nil}
				} else {
					i--
					token = Token{MUL, "*", 0.0, nil}
				}

			case "-":
				i++
				char = getChar(str, i)
				if char == "-" { // 注释
					for getChar(str, i+1) != "\n" && getChar(str, i+1) != "" {
						i++
					}
				} else {
					i--
					token = Token{MINUS, "-", 0.0, nil}
				}

			case "/":
				i++
				char = getChar(str, i)
				if char == "/" { // 注释
					for getChar(str, i+1) != "\n" && getChar(str, i+1) != "" {
						i++
					}
				} else {
					i--
					token = Token{DIV, "/", 0.0, nil}
				}

			default:
				token = Token{ERRTOKEN, "ERR", 0.0, nil}
			}
		}

		if token.TokenType != "" {
			tokens = append(tokens, token)
		}

		i++
	}

	return tokens
}
