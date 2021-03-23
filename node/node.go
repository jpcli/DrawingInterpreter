package node

import (
	"DrawingInterpreter/lexer"
	"math"
)

// Node ：树结点结构体
type Node struct {
	Token  lexer.Token
	Lchild *Node
	Rchild *Node
}

// NewNode ：新建一个树结点
func NewNode(token lexer.Token) *Node {
	return &Node{
		Token:  token,
		Lchild: nil,
		Rchild: nil,
	}
}

// GetValue ：取树的值
func (thisNode *Node) GetValue(t float64) float64 {
	if thisNode.Lchild == nil {
		switch thisNode.Token.TokenType {
		case lexer.CONST_ID:
			return thisNode.Token.Value
		case lexer.T:
			return t
		}

	} else if thisNode.Rchild == nil {
		switch thisNode.Token.TokenType {
		case lexer.PLUS:
			return thisNode.Lchild.GetValue(t)
		case lexer.MINUS:
			return -1 * thisNode.Lchild.GetValue(t)
		case lexer.FUNC:
			return thisNode.Token.FuncPtr(thisNode.Lchild.GetValue(t))
		}
	} else {
		switch thisNode.Token.TokenType {
		case lexer.PLUS:
			return thisNode.Lchild.GetValue(t) + thisNode.Rchild.GetValue(t)
		case lexer.MINUS:
			return thisNode.Lchild.GetValue(t) - thisNode.Rchild.GetValue(t)
		case lexer.MUL:
			return thisNode.Lchild.GetValue(t) * thisNode.Rchild.GetValue(t)
		case lexer.DIV:
			rValue := thisNode.Rchild.GetValue(t)
			if rValue == 0 {
				panic("Semantic error: The divisor cannot be 0.")
			}
			return thisNode.Lchild.GetValue(t) / rValue
		case lexer.POWER:
			return math.Pow(thisNode.Lchild.GetValue(t), thisNode.Rchild.GetValue(t))
		}
	}
	panic("Semantic error: Cannot build the tree.")
}

// GetTree ：获取树结构
func (thisNode *Node) GetTree() map[string]interface{} {
	var ret = map[string]interface{}{"name": thisNode.Token.Lexeme}
	var chi []map[string]interface{}
	if thisNode.Lchild != nil {
		chi = append(chi, thisNode.Lchild.GetTree())
	}
	if thisNode.Rchild != nil {
		chi = append(chi, thisNode.Rchild.GetTree())
	}
	ret["children"] = chi
	return ret
}
