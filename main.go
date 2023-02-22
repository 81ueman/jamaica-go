package main

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func dice(r *rand.Rand) int {
	return r.Intn(6) + 1

}

type problem struct {
	obj  int
	vars [5]int
}

func new() *problem {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &problem{
		obj:  dice(r)*10 + dice(r),
		vars: [5]int{dice(r), dice(r), dice(r), dice(r), dice(r)},
	}
}

func (prob *problem) print() {
	fmt.Printf("obj = %d\n", prob.obj)
	fmt.Printf("vars = %d %d %d %d %d\n", prob.vars[0], prob.vars[1], prob.vars[2], prob.vars[3], prob.vars[4])
}

func eval_Expr(n ast.Expr) (int, error) {

	switch v := n.(type) {
	case *ast.BinaryExpr:
		{
			x, err := eval_Expr(v.X)
			if err != nil {
				return 0, err
			}
			y, err := eval_Expr(v.Y)
			if err != nil {
				return 0, err
			}
			switch v.Op {
			case token.ADD:
				{
					return x + y, nil
				}
			case token.SUB:
				{
					return x - y, nil
				}
			case token.MUL:
				{
					return x * y, nil
				}
			case token.QUO:
				{
					return x / y, nil
				}
			}
		}
	case *ast.BasicLit:
		{
			num, ok := strconv.Atoi(v.Value)
			if ok != nil {
				return 0, errors.New("val is not val")
			}
			return num, nil
		}
	case *ast.ParenExpr:
		{
			return eval_Expr(v.X)
		}
	}
	return 0, errors.New("some went wrong")

}

func main() {
	set := new()
	set.print()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		expr, err := parser.ParseExpr(scanner.Text())
		if err != nil {
			log.Print(err.Error())
			continue
		}
		res, err := eval_Expr(expr)
		if err != nil {
			println(err.Error())
			continue
		}
		fmt.Println(res)
		if res == set.obj {
			fmt.Println("Conguratulation!")
			break
		} else {
			fmt.Println("wrong answer")
			fmt.Printf("obj = %d\n", set.obj)
			fmt.Printf("res = %d\n", res)
		}
	}

}
