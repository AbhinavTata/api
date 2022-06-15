package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/araalinetworks/araali/third_party/api/golang/v2/araalictl"
)

func main() {

	var op, tenantID, tenantName, userEmail, userName string
	flag.StringVar(&op, "op", "ADD", "specify op(ADD/DEL/ADD-USER/DEL-USER)")
	flag.StringVar(&tenantID, "id", "", "specify tenant")
	flag.StringVar(&tenantName, "name", "", "specify tenant name")
	flag.StringVar(&userEmail, "user-email", "", "specify user email")
	flag.StringVar(&userName, "user-name", "", "specify user name")

	flag.Parse()

	araalictl.SetBackend("nightly.aws.araalinetworks.com")
	araalictl.SetToken("eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjdXN0b21lcklEIjoibWV0YS10YXAiLCJlbWFpbCI6InZhbXNpQGFyYWFsaW5ldHdvcmtzLmNvbSIsImV4cCI6MTY4Njc5MjI5MiwiaWF0IjoxNjU1MjU2MjkyLCJqdGkiOiI4d2k1YUsyUTgydjdneWRSV0xBUXpRIiwibmJmIjoxNjU1MjU2MjkyLCJybGlzdCI6ImRvY3MsaW50ZXJhY3QiLCJyb2xlIjoiQXJhYWxpQWRtaW4ifQ.Q_Aj_DxJZwE-eL4UwIcQQWxX3WMcipQZL70Fzha37DZ9KeZqImjbq8drYEBxwYr81AuDluptEC7QQkZ6zHF18ICdb4ItBldxHXxQ-yYqGyIYoUmrP_RzFapM9DbWmmnNhCJgK1W_phzcJQ8FH4jlXXhKBTvn7mXtysMv06clnumqoIXtJrPu9X8j78yPVMScelTACoqmrLWUMavJYnSGe1O0vHEXVe9yjkIyUijILKInnmBbXfX8aJBqa0HD8P5KyKkazOMkbkkW87dC3KxcbfH54-2qKtF7jhW5jzvd5ozD86qIZ9PFvxtJ3mEoc2YB63Ec7pcuNSrZaAby7otrXBZAEOwEt83UFbh9AAlnV8KQ_HBVhxrwCvagUdDtY0SrSXtYCLEBQd2dKat2cahKRc8RIOAvS__PkDwb9Zd7RH4Uo6bxT4BaLkcY5keO9RGy7AU0wNYHf2pkGAs7rPaFnv0ABBgSkpZADa3WcyKF_cZxqGp1B1MXtciAD-3dHP0s9iKMqwYNaZyI7qG6gB3IH5ju5o_OtrOycNM2CNH_yraxX_20gOqkc5FtP-7TSisBJMziVBzGVvz-8u0oR9zZHyx4tAWUjCCrsFs2pYYk23pqLrtnsF6T2E4xz81Tl8ZyiYghfqn9kWPrkDpOa6JUi4XLMPNaZ7v8dp12bpE3ysk")

	if op == "ADD" {
		if userEmail == "" || userName == "" {
			fmt.Println("-user-email and -user-name must be specified when op=ADD")
			os.Exit(1)
		}
		fmt.Println(araalictl.TenantCreate(tenantName, userName, userEmail, true))
	} else if op == "DEL" {
		if tenantID == "" {
			fmt.Println("-id must be specified when op=DEL")
			os.Exit(1)
		}
		fmt.Println(araalictl.TenantDelete(tenantID))
	} else if op == "ADD-USER" {
		if tenantID == "" || userEmail == "" || userName == "" {
			fmt.Println("-id, -user-email and -user-name must be specified when op=ADD-USER")
			os.Exit(1)
		}
		fmt.Println(araalictl.UserAdd(tenantID, userName, userEmail, "ADMIN"))
	} else if op == "DEL-USER" {
		if tenantID == "" || userEmail == "" || userName == "" {
			fmt.Println("-id, -user-email and -user-name must be specified when op=DEL-USER")
			os.Exit(1)
		}
		fmt.Println(araalictl.UserDelete(tenantID, userEmail))
	}
}
