package main

import (
	"fmt"
	"sort"
	"strconv"
)

func singleNumber(nums []int) int {

	// 使用哈希表记录每个数字出现的次数
	mapNum := make(map[int]int)

	for i, e := range nums {
		fmt.Println("i:", i, "e:", e)
		mapNum[e]++
	}
	// 查找出现次数为1的数字
	for k, v := range mapNum {
		if v == 1 {
			return k
		}
	}
	return -1
}

func isPalindrome(x int) bool {
	// 排除小于0的，以及非0的末尾为0的数
	if x < 0 || (x%10 == 0 && x != 0) {
		return false
	}
	var result int
	for x > result {
		result = result*10 + x%10
		x /= 10
	}
	return x == result || x == result/10
}

func isPalindromeForStr(x int) bool {
	if x < 0 {
		return false
	}
	xStr := strconv.Itoa(x)
	runeArr := []rune(xStr)
	for i, j := 0, len(runeArr)-1; i < j; i, j = i+1, j-1 {
		runeArr[i], runeArr[j] = runeArr[j], runeArr[i]
	}
	newStr := string(runeArr)
	return xStr == newStr
}

func isValid(s string) bool {
	stack := []rune{}
	mapKeyValue := map[rune]rune{')': '(', ']': '[', '}': '{'}
	for _, value := range s {
		if value == '(' || value == '{' || value == '[' {
			stack = append(stack, value)
		} else {
			if len(stack) == 0 {
				return false
			}
			topValue := stack[len(stack)-1]
			mapValue, okFlag := mapKeyValue[value]
			if !okFlag || topValue != mapValue {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	s := strs[0]
	for index1, value1 := range s {
		for _, value := range strs {
			if len(value) == index1 || byte(value1) != value[index1] {
				return s[:index1]
			}
		}
	}
	return ""
}

func plusOne(digits []int) []int {

	lenth := len(digits)
	index := lenth - 1
	for index >= 0 {
		value := digits[index]
		if value < 9 {
			digits[index] += 1
			return digits
		}
		digits[index] = 0
		index--
	}
	digits = append([]int{1}, digits...)
	return digits
}

func removeDuplicates(nums []int) int {
	k := 1
	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1] { // nums[i] 不是重复项
			nums[k] = nums[i] // 保留 nums[i]
			k++
		}
	}
	fmt.Print(nums)
	return k
}

func merge(intervals [][]int) [][]int {
	//intervals = [[1,3],[2,6],[8,10],[15,18]]
	lenth := len(intervals)
	if lenth == 1 || lenth == 0 {
		return intervals
	}
	//切片排序核心函数
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	fmt.Println(intervals)

	result := [][]int{intervals[0]}

	for i := 1; i < lenth; i++ {
		last := result[0]
		current := intervals[i]
		if current[0] <= last[1] {
			// 合并后的区间结束取两者的最大值
			last[1] = max(last[1], current[1])
		} else {
			// 无重叠，直接加入结果
			result = append(result, current)
		}
	}
	return result
}

func twoSum(nums []int, target int) []int {
	lenth := len(nums)
	for i := 0; i < lenth; i++ {
		for j := 1; j < lenth; j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return []int{}
}

func main() {
	nums := []int{1, 1, 2}
	result := singleNumber(nums)
	fmt.Println("The single number is:", result)

	fmt.Println(isPalindromeForStr(12321))

	fmt.Println("() is", isValid("()"))
	fmt.Println("()[]{} is", isValid("()[]{}"))
	fmt.Println("(] is", isValid("(]"))
	fmt.Println("([)] is", isValid("([)]"))
	fmt.Println("{[]} is", isValid("{[]}"))
	strs := []string{"flo", "flow", "flight"}

	fmt.Println(longestCommonPrefix(strs))

	fmt.Println(plusOne(nums))
	intervals := [][]int{{1, 3}, {8, 10}, {2, 6}, {13, 18}}
	fmt.Println(removeDuplicates(nums))
	fmt.Println(merge(intervals))
	fmt.Println(twoSum(nums, 26))
}
