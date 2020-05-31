package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
    "bytes"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type Property struct {
	Properties map[string]InnerPorperty `yaml:"properties"`
	Type       string                   `yaml:"type"`
}

type InnerPorperty struct {
	Type string `yaml:"type"`
}

type CRD struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Annotations struct {
			ControllerGenKubebuilderIoVersion string `yaml:"controller-gen.kubebuilder.io/version"`
		} `yaml:"annotations"`
		CreationTimestamp interface{} `yaml:"creationTimestamp"`
		Name              string      `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		AdditionalPrinterColumns []struct {
			JSONPath string `yaml:"JSONPath"`
			Name     string `yaml:"name"`
			Type     string `yaml:"type"`
		} `yaml:"additionalPrinterColumns"`
		Group string `yaml:"group"`
		Names struct {
			Categories []string `yaml:"categories"`
			Kind       string   `yaml:"kind"`
			ListKind   string   `yaml:"listKind"`
			Plural     string   `yaml:"plural"`
			Singular   string   `yaml:"singular"`
		} `yaml:"names"`
		Scope        string `yaml:"scope"`
		Subresources struct {
			Status struct {
			} `yaml:"status"`
		} `yaml:"subresources"`
		Validation struct {
			OpenAPIV3Schema struct {
				Description string `yaml:"description"`
				Properties  struct {
					APIVersion struct {
						Description string `yaml:"description"`
						Type        string `yaml:"type"`
					} `yaml:"apiVersion"`
					Kind struct {
						Description string `yaml:"description"`
						Type        string `yaml:"type"`
					} `yaml:"kind"`
					Metadata struct {
						Type string `yaml:"type"`
					} `yaml:"metadata"`
					Spec struct {
						Description string              `yaml:"description"`
						Properties  map[string]Property `yaml:"properties"`
						Required    []string            `yaml:"required"`
						Type        string              `yaml:"type"`
					} `yaml:"spec"`
				} `yaml:"properties"`
				Required []string `yaml:"required"`
				Type     string   `yaml:"type"`
			} `yaml:"openAPIV3Schema"`
		} `yaml:"validation"`
		Version  string `yaml:"version"`
		Versions []struct {
			Name    string `yaml:"name"`
			Served  bool   `yaml:"served"`
			Storage bool   `yaml:"storage"`
		} `yaml:"versions"`
	} `yaml:"spec"`
}

type CR struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name        string `yaml:"name"`
		Annotations struct {
		} `yaml:"annotations"`
	} `yaml:"metadata"`
	Spec map[string]interface{} `yaml:"spec"`
}

func main() {

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file",
			Usage: "CRD yaml file",
		},
	}

	app.Action = func(c *cli.Context) error {
		var data []byte
		fileInfo, _ := os.Stdin.Stat()

		if fileInfo.Mode()&os.ModeCharDevice == 0 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(os.Stdin)
			data = buf.Bytes()
		} else {
			data, _ = ioutil.ReadFile(c.String("file"))
		}

		o := CRD{}

		err := yaml.Unmarshal([]byte(data), &o)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		cr := CR{
			APIVersion: o.Spec.Group + "/" + o.Spec.Version,
			Kind:       o.Spec.Names.Kind,
			Spec:       make(map[string]interface{}),
		}

		for k, v := range o.Spec.Validation.OpenAPIV3Schema.Properties.Spec.Properties {
			if v.Type == "object" {
				m := make(map[string]string)
				for k1, v1 := range v.Properties {
					m[k1] = v1.Type
				}
				cr.Spec[k] = m
			} else {
				cr.Spec[k] = v.Type
			}
		}

		d, err := yaml.Marshal(&cr)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Printf(string(d))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
