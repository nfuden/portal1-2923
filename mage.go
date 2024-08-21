//go:build mage

package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/magefile/mage/mg"
)

// var Default = FullInstall

// Readme prints out a quick readme
func Readme() {
	color.Cyan("This project is run on mage. To see the available commands, run `mage -l`")
	color.Cyan("To run the full install, run `mage fullinstall`")
	color.Cyan("To run some requests and see the output run `mage demo:run2000`")
	color.Cyan("Each resource folder can be applied either via kubectl or `mage resources:apply<folder_name>`")
}

// // FullInstall is the default install for this repo
// // Run by default
// func FullInstall() {
// 	clean := AmIClean()
// 	if !clean {
// 		color.Cyan("Not clean")
// 		return

// 	}
// 	InstallGloo(false)
// 	Resources{}.Apply()
// }

func smashEnv() (license string, hasBoth bool) {
	license, hasLic := os.LookupEnv("GLOO_LICENSE_KEY")
	otherLic, hasOtherLic := os.LookupEnv("GLOO_EDGE_LICENSE_KEY")
	if hasLic && hasOtherLic {
		return "", true
	}
	if hasOtherLic {
		return otherLic, false
	}
	return license, false
}

type Gateway mg.Namespace

// InstallGloo will install edge with given values
func (g Gateway) Install() {
	toInstall := "install-gateway.sh"

	err := simpleRun(func(cmd *exec.Cmd) {

		cmd.Stderr = os.Stdout
		cmd.Stdout = os.Stdout
		cmd.Env = os.Environ()
		newLic, hasBoth := smashEnv()
		if hasBoth {
			color.Cyan("Both GLOO_LICENSE_KEY and GLOO_EDGE_LICENSE_KEY are set. Please only set one")
		} else {
			cmd.Env = append(cmd.Env, fmt.Sprintf("GLOO_EDGE_LICENSE_KEY=%s", newLic))
		}
	}, filepath.Join(".", "install", toInstall))
	if err != nil {
		color.Red(err.Error())
	}
}

type Portal mg.Namespace

// InstallGloo will install edge with given values
func (p Portal) Install() {

	toInstall := "install-portal.sh"

	err := simpleRun(func(cmd *exec.Cmd) {
		cmd.Stderr = os.Stdout
		cmd.Stdout = os.Stdout
	}, filepath.Join(".", "install", toInstall))
	if err != nil {
		color.Red(err.Error())
	}
}

type Resources mg.Namespace

func (r Resources) ApplyBasicPortal() {
	color.Cyan("Applying portal")
	r.ApplyApi()
	r.ApplyUpstreams()
	r.ApplyUserGroups()
	r.ApplyEnvironments()
	r.ApplyPortal()
	r.ApplyRouteTablesInitial()
	color.Green("Finished applying portal resources")
}

func (r Resources) applyProposedApproach() {

	r.ApplyRouteTablesBackWithAuth()
	r.ApplyVirtualServices()
	color.Green("Finished applying new routing resources")
}

// ApplyUpstreams applies upstreams
func (Resources) ApplyUpstreams() {
	color.Cyan("Applying upstreams")
	applyOrDeleteResource(false, "upstreams")
}
func (Resources) ApplyUserGroups() {
	color.Cyan("Applying users and groups")
	applyOrDeleteResource(false, "user-groups")
}
func (Resources) ApplyEnvironments() {
	color.Cyan("Applying environments")
	applyOrDeleteResource(false, "environment")
}
func (Resources) ApplyPortal() {
	color.Cyan("Applying portal")
	applyOrDeleteResource(false, "portal")
}
func (Resources) ApplyApi() {
	color.Cyan("Applying Httpbin apis")
	applyOrDeleteResource(false, "api/httpbin")
	color.Cyan("Applying Petstore apis")
	applyOrDeleteResource(false, "api/petstore")
}
func (Resources) ApplyRouteTablesInitial() {
	color.Cyan("Applying route tables in the initial state pre solutions")
	applyOrDeleteResource(false, "routetables/initial")
}

func (Resources) ApplyRouteTablesBackWithAuth() {
	color.Cyan("Applying route tables with auth fronting")
	applyOrDeleteResource(false, "routetables/backwithauth")
}
func (Resources) ApplyVirtualServices() {
	color.Cyan("Applying virtual services")
	applyOrDeleteResource(false, "virtualservices")
}

func applyOrDeleteResource(shouldDelete bool, resourcePath string) error {
	files, err := os.ReadDir(filepath.Join(".", resourcePath))
	if err != nil {
		color.Red("failed reading directory:", err.Error())
	}
	if shouldDelete {
		return deleteResource(resourcePath, files)
	}
	return applyResource(resourcePath, files)
}

func applyResource(path string, files []fs.DirEntry) error {
	for _, file := range files {
		err := simpleRun(func(cmd *exec.Cmd) {}, "kubectl", "apply", "-f", filepath.Join(".", path, file.Name()))
		if err != nil {
			color.Red(fmt.Sprintf("failed to apply: %s, %v\n", file.Name(), err))
		}

	}
	fmt.Println()
	return nil
}
func deleteResource(path string, files []fs.DirEntry) error {
	for _, file := range files {
		err := simpleRun(func(cmd *exec.Cmd) {}, "kubectl", "delete", "-f", filepath.Join(".", path, file.Name()))
		if err != nil {
			color.Red(fmt.Sprintf("failed to delete: %s, %v\n", file.Name(), err))
		}

	}
	fmt.Println()
	return nil
}

// TODO: AmIClean will eventually check to see if other resources exist that might clash
func AmIClean() (cleaninstall bool) {
	color.Cyan("TODO: we need to determine what it means to be clean")
	return true
}

func simpleRun(extraFunc func(*exec.Cmd), toRun ...string) error {
	cmd := exec.CommandContext(context.TODO(), toRun[0], toRun[1:]...)

	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout

	extraFunc(cmd)

	err := cmd.Run()
	return err

}

type Demo mg.Namespace

func (d Demo) CheckRoutetable() {
	d.CheckResource("routetables")
}

// CheckResource parses the resources we are showing off
func (Demo) CheckResource(resource string) {
	cmd := exec.CommandContext(context.TODO(), "kubectl", "get", resource, "--all-namespaces")
	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red(err.Error())
		return
	}

	for idx, line := range strings.Split(strings.TrimRight(string(output), "\n"), "\n") {
		if idx == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		fmt.Println()
		color.Cyan(fmt.Sprintf("Namespace: %s, Name: %s", fields[0], fields[1]))
		cmd := exec.CommandContext(context.TODO(), "kubectl", "get", fmt.Sprintf("%s/%s", resource, fields[1]), "-n", fields[0], "-oyaml")
		rTableOut, err := cmd.CombinedOutput()
		if err != nil {
			color.Red(err.Error())
			return
		}
		color.White(string(rTableOut))

	}
}

func (Demo) NoPortalController() {
	cmd := exec.CommandContext(context.TODO(), "kubectl", "scale", "--replicas=0", "-n", "gloo-portal", "deployment", "gloo-portal-controller")
	output, err := cmd.CombinedOutput()
	if err != nil {
		color.Red(err.Error())
		return
	}
	color.Green(string(output))
}

func (d Demo) NewProposed() {
	d.NoPortalController()
	Resources{}.applyProposedApproach()
}
