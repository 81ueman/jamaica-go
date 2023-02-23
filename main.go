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
	"sort"
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
	vars := [5]int{dice(r), dice(r), dice(r), dice(r), dice(r)}
	sort.Ints(vars[:])
	return &problem{
		obj:  dice(r)*10 + dice(r),
		vars: vars,
	}
}

func (prob *problem) print() {
	fmt.Printf("obj = %d\n", prob.obj)
	fmt.Printf("vars = %d %d %d %d %d\n", prob.vars[0], prob.vars[1], prob.vars[2], prob.vars[3], prob.vars[4])
}

func (prob *problem) correct_vars(vs []int) bool {
	if len(vs) != 5 {
		return false
	}
	sort.Ints(vs)
	for i := 0; i < 5; i++ {
		if vs[i] != prob.vars[i] {
			return false
		}
	}
	return true
}

func eval_Expr(n ast.Expr) (int, []int, error) {

	switch v := n.(type) {
	case *ast.BinaryExpr:
		{
			x, xvars, err := eval_Expr(v.X)
			if err != nil {
				return 0, []int{}, err
			}
			y, yvars, err := eval_Expr(v.Y)
			if err != nil {
				return 0, []int{}, err
			}
			switch v.Op {
			case token.ADD:
				{
					return x + y, append(xvars, yvars...), nil
				}
			case token.SUB:
				{
					return x - y, append(xvars, yvars...), nil
				}
			case token.MUL:
				{
					return x * y, append(xvars, yvars...), nil
				}
			case token.QUO:
				{
					return x / y, append(xvars, yvars...), nil
				}
			}
		}
	case *ast.BasicLit:
		{
			num, ok := strconv.Atoi(v.Value)
			if ok != nil {
				return 0, []int{}, errors.New("val is not val")
			}
			return num, []int{num}, nil
		}
	case *ast.ParenExpr:
		{
			return eval_Expr(v.X)
		}
	}
	return 0, []int{}, errors.New("some went wrong")

}

func ask_continue() bool {
	for {
		fmt.Print("continue? [y/n]:")
		var str string
		_, err := fmt.Scan(&str)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if str == "y" || str == "yes" {
			return true
		} else {
			return false
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {

		prob := new()
		prob.print()
		for {
			scanner.Scan()
			expr, err := parser.ParseExpr(scanner.Text())
			if err != nil {
				log.Print(err.Error())
				continue
			}
			res, vars, err := eval_Expr(expr)
			if err != nil {
				println(err.Error())
				continue
			}
			if !prob.correct_vars(vars) {
				fmt.Println("you should use the five vars")
				continue
			}
			fmt.Printf("Your result is : %d\n", res)
			if res == prob.obj {
				fmt.Println("Conguratulation!")
				if ask_continue() {
					break
				} else {
					return
				}
			} else {
				fmt.Println("wrong answer")
				fmt.Printf("obj = %d\n", prob.obj)
				fmt.Printf("res = %d\n", res)
			}
		}
	}

}
