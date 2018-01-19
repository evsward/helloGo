// 要想run，必须package main，默认是文件夹目录名，要更改一下
package pipeline

import "fmt"

// 设想一个拥有三个阶段的管道

/*
 * First Stage: gen
 * params: 一个以逗号分隔的整数列表，数量不限
 * return: 一个通道，包含参数中整数列表的通道
 */
func Gen(nums ... int) <-chan int {
	out := make(chan int)
	// 通过一个goroutine来将参数中的每个整数发送到通道中去。
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out) // close方法作为上面的for循环的终止条件，不能省略。
	}()
	return out
}

/*
 * Second Stage: sq
 * params: 一个包含参数中整数列表的通道
 * return: 一个包含将参数通道中每个整数平方后的列表的通道
 * note: 因为参数和返回值的类型都是相同的整型通道，所以可以反复嵌套该方法。
 */
func Sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n // 平方
		}
		close(out)
	}()
	return out
}

/*
 * Final Stage: main
 * 是一个main函数，没有参数也没有返回值，它相当于客户端调用
 */
func main() {
	c := Gen(2, 3) // 建立通道
	out := Sq(c)   // 通道处理
	// 上面传入两个值2和3，那么这里就要对应的消费两次输出
	fmt.Println(<-out)
	fmt.Println(<-out)
	// 嵌套sq
	for n := range Sq(Sq(Gen(1, 2, 4, 5))) {
		fmt.Println(n)
	}
}

// output:
//  4
//	9
//	1
//	16
//	256
//	625
