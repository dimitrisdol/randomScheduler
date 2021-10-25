package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type podReq struct {
	Name     string
	Category string
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage:\n\t$%s <category>\n", os.Args[0])
	}
	//f, err := os.CreateTemp("./manifests", "pod-*.yml")
	f, err := os.CreateTemp("./manifests", fmt.Sprintf("pod-%s-*.yml", strings.ToLower(os.Args[1])))
	if err != nil {
		log.Fatalf("Error creating temporary file in %s: %v\n", os.TempDir(), err)
	}
	req := podReq{Name: strings.Split(filepath.Base(f.Name()), ".")[0], Category: os.Args[1]}
	tmpl, err := template.New("pod").Parse(podTemplate)
	if err != nil {
		log.Fatalf("Error parsing pod template: %v\n", err)
	}
	if err = tmpl.Execute(f, &req); err != nil {
		log.Fatalf("Error executing pod template: %v\n", err)
	}
	fmt.Printf("Generated pod manifest: %q\n", f.Name())
}

const podTemplate = `---
apiVersion: v1
kind: Pod
metadata:
  name: {{ .Name }}
  namespace: default
  labels:
    category: cat{{ .Category }}
spec:
  schedulerName: default-scheduler
  containers:
  - name: {{ .Name }}
    image: k8s.gcr.io/pause:3.2
`
