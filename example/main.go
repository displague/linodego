package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chiefy/linodego"
)

var linodeClient = linodego.NewClient(nil, nil)
var spendMoney = false

func main() {
	// Trigger endpoints that accrue a balance
	apiToken, apiOk := os.LookupEnv("LINODE_TOKEN")
	spendMoney = spendMoney && apiOk

	var err error
	if err != nil {
		log.Fatal(err)
	}
	linodeClient.SetDebug(false)

	if !apiOk || len(apiToken) == 0 {
		log.Fatal("Could not find LINODE_TOKEN, please assert it is set.")
	}

	// Demonstrate endpoints that require an access token
	linodeClient = linodego.NewClient(&apiToken, nil)
	if err != nil {
		log.Fatal(err)
	}

	moreExamples_authenticated()
}

func moreExamples_authenticated() {
	var linode *linodego.Instance

	linode, err := linodeClient.GetInstance(1231)
	fmt.Println("## Instance request with Invalid ID")
	fmt.Println("### Linode\n", linode, "\n### Error\n", err)

	if spendMoney {
		linode, err = linodeClient.CreateInstance(&linodego.InstanceCreateOptions{Region: "us-central", Type: "g5-nanode-1"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("## Created Instance\n", linode)
	}

	linodes, err := linodeClient.ListInstances(nil)
	fmt.Println("## List Instances")

	if len(linodes) == 0 {
		log.Println("No Linodes to inspect.")
	} else {
		// This is redundantly used for illustrative purposes
		linode, err = linodeClient.GetInstance(linodes[0].ID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("## First Linode\n", linode)

		configs, err := linodeClient.ListInstanceConfigs(linode.ID, nil)
		if err != nil {
			log.Fatal(err)
		} else if len(configs) > 0 {
			config, err := linodeClient.GetInstanceConfig(linode.ID, configs[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("### First Config:\n", config)
		} else {
			fmt.Println("### No Configs")
		}

		disks, err := linodeClient.ListInstanceDisks(linode.ID, nil)
		if err != nil {
			log.Fatal(err)
		} else if len(disks) > 0 {
			disk, err := linodeClient.GetInstanceDisk(linode.ID, disks[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("### First Disk\n", disk)
		} else {
			fmt.Println("### No Disks")
		}

		backups, err := linodeClient.GetInstanceBackups(linode.ID)
		if err != nil {
			log.Fatal(err)
		}
		if len(backups.Automatic) > 0 {
			fmt.Println("### First Auto Backup\n", backups.Automatic[0])
		} else {
			fmt.Println("### No Auto Backups")
		}
		fmt.Println("### Snapshots\n", backups.Snapshot)
		if backups.Snapshot.Current != nil {
			// snapshot fetched will be exactly the same as backups.Snapshot.Current
			// just being redundant for illustrative purposes
			if snapshot, err := linodeClient.GetInstanceSnapshot(linode.ID, backups.Snapshot.Current.ID); err == nil {
				fmt.Println("#### Current\n", snapshot)
			} else {
				fmt.Println("#### No Current Snapshot\n", err)
			}
		} else {
			fmt.Println("### No Current Snapshot")
		}

		volumes, err := linodeClient.ListInstanceVolumes(linode.ID, nil)
		if err != nil {
			log.Fatal(err)
		} else if len(volumes) > 0 {
			volume, err := linodeClient.GetInstanceVolume(linode.ID, volumes[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("### First Volume\n", volume)
		} else {
			fmt.Println("### No Volumes")
		}

		stackscripts, err := linodeClient.ListStackscripts(&linodego.ListOptions{Filter: "{\"mine\":true}"})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("## Your Stackscripts\n", stackscripts)
	}
}