package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/codegangsta/cli"
	"github.com/mgutz/ansi"
	"github.com/remind101/empctl/pkg/awsutil/elbutil"
)

func init() {
	AddCommand(cli.Command{
		Name:   "elb",
		Usage:  "Check the ELB status for an app",
		Flags:  commonFlags,
		Action: cmdELBInfo,
	})
}

func cmdELBInfo(c *cli.Context) {
	client := elbutil.New(elb.New(nil), ec2.New(nil))
	lbs, err := client.ListByTags(map[string]string{"App": c.String("app")})
	if err != nil {
		log.Fatal(err)
	}

	for _, lb := range lbs {
		fmt.Printf("Name:       %s\n", lb.Name)
		fmt.Printf("DNSName:    %s\n", lb.DNSName)
		fmt.Printf("SG:         %s\n", strings.Join(lb.SecurityGroups, ", "))

		ss := make([]string, len(lb.InstanceStates))
		for i, s := range lb.InstanceStates {
			if s.State == "InService" {
				ss[i] = ansi.Color(s.InstanceID, "green")
			} else {
				ss[i] = ansi.Color(s.InstanceID, "red")
			}
		}
		fmt.Printf("Instances:  %s\n", strings.Join(ss, ", "))
		fmt.Println()
	}
}
