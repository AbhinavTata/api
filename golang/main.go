package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/araalinetworks/api/golang/araalictl"
	"gopkg.in/yaml.v2"
)

func main() {
	// Set araalictl path to a differnt value.
	araalictl.SetAraalictlPath("./araalictl")
	for {
		fmt.Printf("\nEnter command to run:\n")
		fmt.Printf("\t0: quit\n")
		fmt.Printf("\t1: authorize\n")
		fmt.Printf("\t2: deauthorize\n")
		fmt.Printf("\t3: zones_apps\n")
		fmt.Printf("\t4: zones_apps_links\n")
		fmt.Printf("\t5: summary\n")
		fmt.Printf("\t6: alert_card\n")
		fmt.Printf("\t7: alerts\n")
		fmt.Printf("\t8: compute\n")
		fmt.Printf("\t9: all_alerts\n")

		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		text = strings.TrimSpace(text)

		if text == "0" {
			return
		}

		if text == "1" {
			registeredEmailId := "**********" // set the email registered with araali
			token := "**********"             // set to the token value obtained for api access.
			araalictl.Authorize(registeredEmailId, token, false)
		}

		if text == "2" {
			araalictl.DeAuthorize(false)
		}

		if text == "3" {
			zones, err := araalictl.GetZones(false, "")
			if err != nil {
				fmt.Println(err)
			} else {
				for _, zone := range zones {
					fmt.Printf("%s:\n", zone.ZoneName)
					for _, app := range zone.Apps {
						fmt.Printf("\t%s\n", app.AppName)
					}
					fmt.Println()
				}
			}
			continue
		}

		if text == "4" {
			zones, err := araalictl.GetZones(false, "")
			if err != nil {
				fmt.Println(err)
			} else {
				for _, zone := range zones {
					fmt.Printf("%s:\n", zone.ZoneName)
					for _, app := range zone.Apps {
						fmt.Printf("\t%s\n", app.AppName)
						for _, link := range app.Links {
							fmt.Printf("\t\t%+v\n", link)
						}
					}
					fmt.Println()
				}
			}
			continue
		}

		if text == "5" {
			summary := make(map[string]int)
			zones, err := araalictl.GetZones(false, "")
			if err != nil {
				fmt.Println(err)
			} else {
				for _, zone := range zones {
					for _, app := range zone.Apps {
						for _, link := range app.Links {
							summary["type."+link.Type] += 1
							summary["state."+link.State] += 1
						}
					}
				}
				fmt.Printf("\nSummary:\n")
				for k, v := range summary {
					fmt.Printf("\t%-30s %d\n", k, v)
				}
			}
			continue
		}

		if text == "6" {
			alertCard, err := araalictl.GetAlertCard("")
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("%v\n", alertCard)
			}
		}

		if text == "7" {
			startTime := time.Now().Add(-(3 * araalictl.ONE_DAY)).Unix()
			alertPage, err := araalictl.GetAlerts("", startTime, 0, 25, false)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Fetched %d alerts.\n", len(alertPage.Alerts))
				for {
					if !alertPage.HasNext() {
						fmt.Println("Done fetching!")
						break
					}
					alertPage.NextPage()
					fmt.Printf("Fetched %d alerts.\n", len(alertPage.Alerts))
				}
			}
		}

		if text == "8" {
			computeList, err := araalictl.GetCompute("", "", "")
			if err != nil {
				fmt.Println(err)
			} else {
				computeL, _ := yaml.Marshal(computeList)
				fmt.Println(string(computeL))
			}
		}
		if text == "9" {
			startTime := time.Now().Add(-(3 * araalictl.ONE_DAY)).Unix()
			alertPage, err := araalictl.GetAlerts("", startTime, 0, 25, true)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Fetched %d alerts.\n", len(alertPage.Alerts))
				for {
					if !alertPage.HasNext() {
						fmt.Println("Done fetching!")
						break
					}
					alertPage.NextPage()
					fmt.Printf("Fetched %d alerts.\n", len(alertPage.Alerts))
				}
			}
		}
	}

	// unreachable code, left for sample reasons
	links, _ := araalictl.GetLinks("azure3", "istio-system", "amk")
	fmt.Printf("%+v\n", links)
}
