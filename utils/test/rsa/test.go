package main

import (
	"fmt"
	"go_crud/server/user/utils"
)

func main() {
	var str string
	_, err := fmt.Scanln(&str)
	if err != nil {
		return
	}
	encryptedData, err := utils.RsaEncode(str)
	if err != nil {
		fmt.Println("加密失败:", err)
	} else {
		fmt.Printf("%q\n", encryptedData)
	}
}
