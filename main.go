package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"net"
	"strings"
	"runtime"
	"text/template"

	"github.com/fatih/color"
)

var fmtSucc = color.New(color.Bold, color.FgGreen)
var fmtInfo = color.New(color.FgYellow)
var fmtErr = color.New(color.Bold, color.FgRed)

var dotfilesSourcePath = "mydotfiles"
var params map[string]interface{}

func main() {
	fmtSucc.Println("Deploying your dotfiles…")

	// find all relevant files
	fmtInfo.Println("Searching deployable files…")
	files := getDeployableFiles(dotfilesSourcePath)
	fmtInfo.Println(fmt.Sprintf("Found %d files.", len(files)))

	params = map[string]interface{}{}
	if len(files) > 0 {
		// TODO read settings from json
		// TODO add invariants/system variables to settings
		hostname, err := os.Hostname()
		if err == nil {
			params["hostname"] = hostname
		}
		// operating system name
		params["os"] = runtime.GOOS
		// ip addresses
		addrs, err := net.InterfaceAddrs()
		if err != nil {
		    panic(err)
		}
		ipList := make([]string, len(addrs))
		for i, addr := range addrs {
		    addrParts := strings.Split(addr.String(), "/")
		    ipList[i] = addrParts[0]
		}
		params["ip"] = ipList
 	}

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	// loop the files to deploy
	for _, filename := range files {
		fmtInfo.Println(fmt.Sprintf("Processing file %s…", filename))
		outFilename := path.Join(usr.HomeDir, filename)
		// TODO backup files to be overridden
		applyTemplate(filename, outFilename, params)

		fmtInfo.Println(fmt.Sprintf("Deployed file %s to %s", filename, outFilename))
	}
	fmtSucc.Println("Deployment finished.")
}

func getDeployableFiles(root string) []string {
	files := []string{}
	err := filepath.Walk(root, func(fpath string, info os.FileInfo, err error) error {
		if fpath == root || info.IsDir() {
			return nil
		}
		if !strings.HasPrefix(path.Base(fpath), "._") {
			files = append(files, strings.TrimPrefix(fpath, root))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func applyTemplate(templateFilename, outputFilename string, params map[string]interface{}) {
	// prepare template
	templ := template.New("dotfile")

	templ = templ.Funcs(template.FuncMap{
		"include_file": func(filename string) string {
			fileContent, err := ioutil.ReadFile(path.Join(dotfilesSourcePath, filename))
			if err != nil {
				panic(err)
			}
			return string(fileContent)
		},
		"has_ip": func(ip string) bool {
			ipList, ok := params["ip"].([]string)
			if !ok {
				return false
			}
			for _, interfaceIp := range ipList {
				if ip == interfaceIp {
					return true
				}
			}
			return false
		},
	})

	templ, err := templ.ParseFiles(path.Join(dotfilesSourcePath, templateFilename))
	if err != nil {
		panic(err)
	}

	f, err := os.Create(outputFilename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)
	err = templ.ExecuteTemplate(w, path.Base(templateFilename), params)
	if err != nil {
		panic(err)
	}
	w.Flush()
	return
}

func cp(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}
