package main

import "fmt"

func bsearch0(arr []int,target int){
	l,r := 0,len(arr)-1
	for l<=r {
		mid := l+(r-l)>>1
		if arr[mid]<=target{
			if mid == len(arr) || arr[mid+1] > target{
				fmt.Println(arr[mid],mid)
				break
			}else {
				l = mid+1
			}
		}else{
			r= mid-1
		}
	}
}

func bsearch1(arr []int,target int) {
	l,r := 0 ,len(arr)-1
	for l<=r {
		mid := l+(r-l)>>1
		if arr[mid]>=target {
			if mid == 0 || arr[mid-1] < target{
				fmt.Println(arr[mid],mid)
				break
			}else{
				r = mid-1
			}
		}else{
			//肯定在右边区域
			l = mid+1
		}
	}
}

func bsearch2(arr []int,target int) {
	l,r := 0,len(arr)-1
	for l<=r {
		mid := l+(r-l)>>1
		if arr[mid] == target {
			if mid==0 || arr[mid-1] < target {
				fmt.Println(arr[mid],mid)
				break
			}else {
				r = mid-1
			}
		}else if arr[mid]<target{
			l = mid+1
		}else {
			r = mid-1
		}
	}
}

func bsearch3(arr []int,target int){
	l,r := 0,len(arr)-1
	for l<=r {
		mid := l+(r-l)>>1
		if arr[mid] == target {
			if mid == len(arr)-1 || arr[mid+1] > target{
				fmt.Println(arr[mid],mid)
				break
			}else {
				l = mid+1
			}
		}else if arr[mid] > target{
			l = mid+1
		}else {
			r= mid-1
		}
	}
}

func bsearch4(arr []int,target int){
	l,r := 0,len(arr)-1
	for l<=r {
		mid := l+(r-l)>>1
		if arr[mid] <= target {
			l = mid+1
		}else {
			r = mid-1
		}
	}
	fmt.Println(l)
}

func bsearch5(arr []int,target int){
	l,r := 0,len(arr)-1
	for l<=r {
		mid := l+(r-l)>>1
		if arr[mid] >= target {
			r = mid-1
		}else {
			l = mid+1
		}
	}
	fmt.Println(r)
}

func main() {
	exp := []int{1,2,4,4,5}
	bsearch4(exp,6)
}