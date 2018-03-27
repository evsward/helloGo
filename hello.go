package main

import (
	"fmt"
	"math"
	"time"
)

const ab string = "const here"

func main() {
	fmt.Println("hello, world!\n")
	fmt.Println("7.5/2.5=", 7.5/2.5)
	fmt.Println(true && false)
	fmt.Println(true || false)
	fmt.Println(!true)

	var a string = "initial"
	fmt.Println(a)

	var b, c int = 1, 2
	fmt.Println(b + c)

	var d = true
	fmt.Println(d)

	var e int
	fmt.Println(e)

	f := "short"
	fmt.Println(f)

	fmt.Println(ab)
	const n = 500000000
	const g = 3e20 / n
	fmt.Println(g)

	fmt.Println(int64(g))

	fmt.Println(math.Sin(n))

	i := 1
	for i <= 3 {
		fmt.Println(i)
		i = i + 1
	}

	for j := 7; j <= 9; j++ {
		fmt.Println(j)
	}

	for {
		fmt.Println("loop")
		break
	}

	for n := 0; n <= 5; n++ {
		if n%2 == 0 {
			continue
		}
		fmt.Println(n)
	}

	switch time.Now().Weekday() {
	case time.Friday:
		fmt.Println("friday")
	case time.Monday, time.Tuesday, time.Wednesday:
		fmt.Println("workday")
	default:
		fmt.Println("day")
	}

	t := time.Now()

	switch {
	case t.Hour() > 12:
		fmt.Println("good afternoon")
	case t.Hour() < 12:
		fmt.Println("good morning")
	default:
		fmt.Println("hello")
	}

	whoami := func(i interface{}) {
		switch t := i.(type) {
		case bool:
			fmt.Println("I'm a bool")
		case int:
			fmt.Println("I'm a int")
		default:
			fmt.Printf("don't know type %T%s\n", t,i) // %T shows the type of variable.
		}
	}
	whoami(12)
	whoami("heyhey")
	whoami(true)

	var arr [5]int
	fmt.Println("emp", arr)

	arr[4] = 100
	fmt.Println("set", arr)
	fmt.Println("get", arr[4])
	fmt.Println("len", len(arr))
	brr := [5]int{1, 2, 3, 4, 5}
	fmt.Println("brr=", brr)

	var twoD [2][3]int
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("twoD=", twoD)

	s := make([]string, 3)
	fmt.Println("emp", s)
	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("emp", s)

	s = append(s, "d", "e")
	fmt.Println("emp", s)

	ck := make([]string, len(s))
	copy(ck, s)
	fmt.Println("ckemp", ck)
	fmt.Println("sl1:", s[2:4])
	fmt.Println("sl2:", s[:6])
	fmt.Println("sl3:", s[2:])

	tt := []string{"g", "h", "i"}
	fmt.Println("ttdcl:", tt)
}
