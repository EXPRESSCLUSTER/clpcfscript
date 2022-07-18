/*
 * Convert clp.conf (XML file) to bash script
 */
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type conf struct {
	All       all       `xml:"all"`
	Cluster   cluster   `xml:"cluster"`
	Server    []server  `xml:"server"`
	Heartbeat heartbeat `xml:"heartbeat"`
	Group     []group   `xml:"group"`
	Resource  resource  `xml:"resource"`
}
type all struct {
	Charset  string `xml:"charset"`
	ServerOS string `xml:"serveros"`
	Encode   string `xml:"encode"`
}
type cluster struct {
	Name string `xml:"name"`
}
type server struct {
	Name     string   `xml:"name,attr"`
	Priority int      `xml:"priority"`
	Device   []device `xml:"device"`
}
type device struct {
	ID   int    `xml:"id,attr"`
	Type string `xml:"type"`
	Info string `xml:"info"`
}
type heartbeat struct {
	Types  []types  `xml:"types"`
	LANHB  []lanhb  `xml:"lanhb"`
	LANKHB []lankhb `xml:"lankhb"`
}
type types struct {
	Name string `xml:"name,attr"`
}
type lanhb struct {
	Name     string `xml:"name,attr"`
	Priority int    `xml:"priority"`
	Device   int    `xml:"id"`
}
type lankhb struct {
	Name     string `xml:"name,attr"`
	Priority int    `xml:"priority"`
	Device   int    `xml:"id"`
}
type group struct {
	Name          string        `xml:"name,attr"`
	Groupresource []grpResource `xml:"resource"`
}
type grpResource struct {
	Name string `xml:"name,attr"`
}
type resource struct {
	MD   []md   `xml:"md"`
	EXEC []exec `xml:"exec"`
}
type md struct {
	Name       string        `xml:"name,attr"`
	Parameters mdParameeters `xml:"parameters"`
}
type mdParameeters struct {
	Netdev  []netdev `xml:"netdev"`
	NMPPath string   `xml:"nmppath"`
	Mount   mount    `xml:"mount"`
	Diskdev diskdev  `xml:"diskdev"`
	FS      string   `xml:"fs"`
}
type netdev struct {
	ID       int    `xml:"id,attr"`
	Priority int    `xml:"priority"`
	Device   int    `xml:"device"`
	MDCName  string `xml:"mdcname"`
}
type mount struct {
	Point string `xml:"point"`
}
type diskdev struct {
	DPPath string `xml:"dppath"`
	CPPath string `xml:"cppath"`
}
type exec struct {
	Name string `xml:"name,attr"`
}

func main() {
	var str string
	var (
		m = flag.Int("m", 0, "Debug flag (0: Normal, 1: Debug Mode)")
	)

	flag.Parse()

	/* Read clp.conf */
	//	clpconf, err := ioutil.ReadFile("./clp.conf")
	clpconf, err := ioutil.ReadFile("./clp.conf")
	if err != nil {
		fmt.Println("ReadFile() failed: ", err)
		os.Exit(1)
	}
	//	fmt.Println(string(clpconf))
	data := new(conf)
	if err := xml.Unmarshal(clpconf, data); err != nil {
		fmt.Println("Unmarshal() failed: ", err)
		os.Exit(1)
	}
	os.Mkdir("conf", 0755)
	fp, err := os.Create("conf/create-cluster.sh")
	/* TODO: Error handling */
	os.Chmod("conf/create-cluster.sh", 755)
	/* TODO: Error handling */

	/* Initialize the file */
	if *m == 1 {
		fmt.Println(data.Cluster.Name)
		fmt.Println(data.All.Charset)
		fmt.Println(data.All.ServerOS)
	}
	str = fmt.Sprintf("clpcfset create %s %s %s \n", data.Cluster.Name, data.All.Charset, data.All.ServerOS)
	fp.WriteString(str)

	/* Add server */
	for i := 0; i < len(data.Server); i++ {
		if *m == 1 {
			fmt.Println(data.Server[i].Name)
			fmt.Println(data.Server[i].Priority)
		}
		str = fmt.Sprintf("clpcfset add srv %s %d \n", data.Server[i].Name, data.Server[i].Priority)
		fp.WriteString(str)
		for j := 0; j < len(data.Server[i].Device); j++ {
			if data.Server[i].Device[j].ID >= 10700 {
				// fmt.Println("http")
			} else if data.Server[i].Device[j].ID >= 400 {
				/* Add device (mdc) */
				if *m == 1 {
					fmt.Println(data.Server[i].Device[j].Type)
					fmt.Println(data.Server[i].Device[j].Info)
				}
				str = fmt.Sprintf("clpcfset add device %s mdc %d %s \n", data.Server[i].Name, (data.Server[i].Device[j].ID - 400), data.Server[i].Device[j].Info)
				fp.WriteString(str)
			} else {
				/* Add device (lan) */
				if *m == 1 {
					fmt.Println(data.Server[i].Device[j].Info)
				}
				str = fmt.Sprintf("clpcfset add device %s lan %d %s \n", data.Server[i].Name, data.Server[i].Device[j].ID, data.Server[i].Device[j].Info)
				fp.WriteString(str)
			}
		}

	}
	/* Add heartbeat (lankhb) */
	for i := 0; i < len(data.Heartbeat.LANKHB); i++ {
		if *m == 1 {
			fmt.Println(data.Heartbeat.LANKHB[i].Name)
			fmt.Println(data.Heartbeat.LANKHB[i].Device)
			fmt.Println(data.Heartbeat.LANKHB[i].Priority)
		}
		str = fmt.Sprintf("clpcfset add hb lankhb %d %d\n", data.Heartbeat.LANKHB[i].Device, data.Heartbeat.LANKHB[i].Priority)
		fp.WriteString(str)
	}
	/* Add group (failover) */
	/* TODO: Think about Management Group */
	for i := 0; i < len(data.Group); i++ {
		str = fmt.Sprintf("clpcfset add grp failover %s\n", data.Group[i].Name)
		fp.WriteString(str)
	}
	/* Add resource (md) */
	for i := 0; i < len(data.Resource.MD); i++ {
		for j := 0; j < len(data.Group); j++ {
			for k := 0; k < len(data.Group[j].Groupresource); k++ {
				if strings.Contains(data.Group[j].Groupresource[k].Name, data.Resource.MD[i].Name) {
					//fmt.Println(data.Group[j].Groupresource[k].Name)
					str = fmt.Sprintf("clpcfset add rsc %s md %s \n", data.Group[j].Name, data.Resource.MD[i].Name)
					fp.WriteString(str)
					for l := 0; l < len(data.Resource.MD[i].Parameters.Netdev); l++ {
						str = fmt.Sprintf("clpcfset add rscparam md %s parameters/netdev@%d/priority %d\n", data.Resource.MD[i].Name, data.Resource.MD[i].Parameters.Netdev[l].ID, data.Resource.MD[i].Parameters.Netdev[l].Priority)
						fp.WriteString(str)
						str = fmt.Sprintf("clpcfset add rscparam md %s parameters/netdev@%d/device %d\n", data.Resource.MD[i].Name, data.Resource.MD[i].Parameters.Netdev[l].ID, data.Resource.MD[i].Parameters.Netdev[l].Device-400)
						fp.WriteString(str)
						str = fmt.Sprintf("clpcfset add rscparam md %s parameters/netdev@%d/mdcname %s\n", data.Resource.MD[i].Name, data.Resource.MD[i].Parameters.Netdev[l].ID, data.Resource.MD[i].Parameters.Netdev[l].MDCName)
						fp.WriteString(str)
					}
					str = fmt.Sprintf("clpcfset add rscparam md %s parameters/nmppath %s\n", data.Resource.MD[i].Name, data.Resource.MD[i].Parameters.NMPPath)
					fp.WriteString(str)
					str = fmt.Sprintf("clpcfset add rscparam md %s parameters/mount/point %s\n", data.Resource.MD[i].Name, data.Resource.MD[i].Parameters.Mount.Point)
					fp.WriteString(str)
					str = fmt.Sprintf("clpcfset add rscparam md %s parameters/diskdev/dppath %s\n", data.Resource.MD[i].Name, data.Resource.MD[i].Parameters.Diskdev.DPPath)
					fp.WriteString(str)
					str = fmt.Sprintf("clpcfset add rscparam md %s parameters/diskdev/cppath %s\n", data.Resource.MD[i].Name, data.Resource.MD[i].Parameters.Diskdev.CPPath)
					fp.WriteString(str)
					str = fmt.Sprintf("clpcfset add rscparam md %s parameters/fs %s\n", data.Resource.MD[i].Name, data.Resource.MD[i].Parameters.FS)
					fp.WriteString(str)
				}
			}
		}
	}
	/*
		// Add resource (exec)
		for i := 0; i < len(data.Resource.EXEC); i++ {
			for j := 0; j < len(data.Group); j++ {
				for k := 0; k < len(data.Group[j].Groupresource); k++ {
					if strings.Contains(data.Group[j].Groupresource[k].Name, data.Resource.EXEC[i].Name) {
						//fmt.Println(data.Group[j].Groupresource[k].Name)
						str = fmt.Sprintf("clpcfset add rsc %s exec %s \n", data.Group[j].Name, data.Resource.EXEC[i].Name)
						fp.WriteString(str)
					}
				}
			}
		}
	*/
}
