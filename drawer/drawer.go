package drawer

import (
	"DrawingInterpreter/lexer"
	"DrawingInterpreter/node"
	"DrawingInterpreter/parser"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"time"
)

type param struct {
	OriginX float64
	OriginY float64
	Angle   float64
	ScaleX  float64
	ScaleY  float64
}

// Draw ：开始画图
func Draw(statements []parser.Statement) string {
	thisParam := param{0, 0, 0, 1, 1}
	pic := image.NewRGBA(image.Rect(0, 0, 500, 500))
	for _, statement := range statements {
		if statement["statement"].(string) == lexer.ORIGIN {
			thisParam.OriginX = statement["x"].(*node.Node).GetValue(0)
			thisParam.OriginY = statement["y"].(*node.Node).GetValue(0)

		} else if statement["statement"].(string) == lexer.ROT {
			thisParam.Angle = statement["angle"].(*node.Node).GetValue(0)

		} else if statement["statement"].(string) == lexer.SCALE {
			thisParam.ScaleX = statement["x"].(*node.Node).GetValue(0)
			thisParam.ScaleY = statement["y"].(*node.Node).GetValue(0)

		} else if statement["statement"].(string) == lexer.FOR {
			now := statement["begin"].(*node.Node).GetValue(0)
			end := statement["end"].(*node.Node).GetValue(0)
			step := statement["step"].(*node.Node).GetValue(0)
			for now <= end {
				// 获取初值
				x := statement["x"].(*node.Node).GetValue(now)
				y := statement["y"].(*node.Node).GetValue(now)
				// 比例变换
				x, y = x*thisParam.ScaleX, y*thisParam.ScaleY
				// 旋转变换
				x, y = x*math.Cos(thisParam.Angle)+y*math.Sin(thisParam.Angle), y*math.Cos(thisParam.Angle)-x*math.Sin(thisParam.Angle)
				// 平移变换
				x, y = x+thisParam.OriginX, y+thisParam.OriginY
				// 取整
				x1 := int(x)
				y1 := int(y)
				// 画点
				draw.Draw(pic, image.Rect(x1, y1, x1+1, y1+1), &image.Uniform{color.RGBA{0, 0, 255, 255}}, image.ZP, draw.Src)

				now += step
			}

		} else {
			panic("Semantic error: Unexpected statement.")
		}
	}
	filename := fmt.Sprintf("%d.png", time.Now().UnixNano())
	picFile, _ := os.Create("./pic/" + filename)
	png.Encode(picFile, pic)
	picFile.Close()
	return filename
}
