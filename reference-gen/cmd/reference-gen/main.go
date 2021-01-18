package main

import (
	"github.com/oauth2-proxy/tools/reference-gen/pkg/generator"
	flag "github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var (
	packageName   = flag.String("package", "", "api directory (or import path), for the package for which references should be generated")
	requiredTypes = flag.StringSlice("types", []string{}, "types from the package for which references should be generated")
	templateDir   = flag.String("template-dir", "", "path to output templates dir, if unset uses default templates")
	headerFile    = flag.String("header-file", "", "file including header text to prepend to generated data")
	outputFile    = flag.String("out-file", "", "path to output file to save the result")
)

func main() {
	klog.InitFlags(nil)
	flag.Set("logtostderr", "true")
	flag.Parse()

	gen, err := generator.NewGenerator(*packageName, *requiredTypes, *headerFile, *outputFile, *templateDir)
	if err != nil {
		klog.Fatalf("error constructing generator: %v", err)
	}

	klog.Infof("Running generator on package %q", *packageName)
	if err := gen.Run(); err != nil {
		klog.Fatalf("error running generator: %v", err)
	}
}
