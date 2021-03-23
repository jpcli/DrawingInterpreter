package parser

import (
	"DrawingInterpreter/lexer"
	"DrawingInterpreter/node"
	"fmt"
)

var tokens []lexer.Token
var nowToken lexer.Token
var numToken int // Token数
var iToken int   // 当前下标

func getToken() lexer.Token {
	if iToken < numToken {
		nowToken = tokens[iToken]
		iToken++
		return nowToken
	}
	panic("Parser error: Token index out of range.")

}

func matchToken(tokenType string) {
	if nowToken.TokenType == tokenType {
		getToken()
	} else {
		panic(fmt.Sprintf("Parser error: Expected ' %s ' , but not found.", tokenType))
	}

}

// 处理二元 + -
func e() *node.Node {
	left := t()
	for nowToken.TokenType == lexer.PLUS || nowToken.TokenType == lexer.MINUS {
		root := node.NewNode(nowToken)
		matchToken(nowToken.TokenType)
		right := t()
		root.Lchild = left
		root.Rchild = right
		left = root
	}
	return left
}

// 处理二元 * /
func t() *node.Node {
	left := f()
	for nowToken.TokenType == lexer.MUL || nowToken.TokenType == lexer.DIV {
		root := node.NewNode(nowToken)
		matchToken(nowToken.TokenType)
		right := f()
		root.Lchild = left
		root.Rchild = right
		left = root
	}
	return left
}

// 处理一元 + -
func f() *node.Node {
	if nowToken.TokenType == lexer.PLUS || nowToken.TokenType == lexer.MINUS {
		root := node.NewNode(nowToken)
		matchToken(nowToken.TokenType)
		child := f()
		root.Lchild = child
		return root
	}
	return c()

}

// 处理 **
func c() *node.Node {
	left := a()
	if nowToken.TokenType == lexer.POWER {
		root := node.NewNode(nowToken)
		matchToken(nowToken.TokenType)
		right := c()
		root.Lchild = left
		root.Rchild = right
		left = root
	}
	return left
}

// 处理终结符
func a() *node.Node {
	if nowToken.TokenType == lexer.CONST_ID || nowToken.TokenType == lexer.T {
		root := node.NewNode(nowToken)
		matchToken(nowToken.TokenType)
		return root
	} else if nowToken.TokenType == lexer.FUNC {
		root := node.NewNode(nowToken)
		matchToken(nowToken.TokenType)
		matchToken(lexer.L_BRACKET)
		child := e()
		matchToken(lexer.R_BRACKET)
		root.Lchild = child
		return root
	} else if nowToken.TokenType == lexer.L_BRACKET {
		matchToken(lexer.L_BRACKET)
		root := e()
		matchToken(lexer.R_BRACKET)
		return root
	}
	panic(fmt.Sprintf("Parser error: Unexpected terminal symbol ' %s '.", nowToken.Lexeme))
}

// Statement ：语句描述
type Statement map[string]interface{}

// 分析程序，识别一条条语句
func parseProgram() []Statement {
	var statements []Statement

	for nowToken.TokenType != lexer.NONTOKEN {
		temp := parseStatement()
		matchToken(lexer.SEMICO)
		statements = append(statements, temp)
	}
	return statements
}

// 分析语句，用于语义分析
func parseStatement() Statement {
	var statement Statement
	if nowToken.TokenType == lexer.ORIGIN {
		statement = originStatement()
	} else if nowToken.TokenType == lexer.ROT {
		statement = rotStatement()
	} else if nowToken.TokenType == lexer.SCALE {
		statement = scaleStatement()
	} else if nowToken.TokenType == lexer.FOR {
		statement = forStatement()
	} else {
		panic("Parser error: Unexpected statement.")
	}
	return statement
}

// 坐标平移语句
func originStatement() Statement {
	var statement = Statement{"statement": lexer.ORIGIN}
	matchToken(lexer.ORIGIN)
	matchToken(lexer.IS)
	matchToken(lexer.L_BRACKET)
	statement["x"] = e()
	matchToken(lexer.COMMA)
	statement["y"] = e()
	matchToken(lexer.R_BRACKET)
	return statement
}

// 角度旋转语句
func rotStatement() Statement {
	var statement = Statement{"statement": lexer.ROT}
	matchToken(lexer.ROT)
	matchToken(lexer.IS)
	statement["angle"] = e()
	return statement
}

// 比例缩放语句
func scaleStatement() Statement {
	var statement = Statement{"statement": lexer.SCALE}
	matchToken(lexer.SCALE)
	matchToken(lexer.IS)
	matchToken(lexer.L_BRACKET)
	statement["x"] = e()
	matchToken(lexer.COMMA)
	statement["y"] = e()
	matchToken(lexer.R_BRACKET)
	return statement
}

// 循环绘图语句
func forStatement() Statement {
	var statement = Statement{"statement": lexer.FOR}
	matchToken(lexer.FOR)
	matchToken(lexer.T)
	matchToken(lexer.FROM)
	statement["begin"] = e()
	matchToken(lexer.TO)
	statement["end"] = e()
	matchToken(lexer.STEP)
	statement["step"] = e()
	matchToken(lexer.DRAW)
	matchToken(lexer.L_BRACKET)
	statement["x"] = e()
	matchToken(lexer.COMMA)
	statement["y"] = e()
	matchToken(lexer.R_BRACKET)
	return statement
}

// Parse ：语法分析
func Parse(t []lexer.Token) []Statement {
	tokens = t
	numToken = len(tokens)
	iToken = 0
	getToken()
	return parseProgram()
}
