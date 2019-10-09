package cases

func TwoSum(nums []int, target int) []int {
    res := []int{}
    for k1,v1 := range nums {
        for k2,v2 := range nums {
            if ((v1 + v2) == target) {
                res = []int{k1,k2}
                break
            }
        }
    }
    return res
}

func TwoSumNew(nums []int, target int) []int {
    m := make(map[int]int, len(nums))
    res := []int{}
    for k,v := range nums {
        if _, ok := m[target-v]; ok {
            res = []int{k, m[target-v]}
            break
        }
        m[v] = k
    }
    return res
}