package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
    "bytes"
    "time"

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
			ControllerGenKubebuilderIoVersion           string `yaml:"controller-gen.kubebuilder.io/version"`
			KubectlKubernetesIoLastAppliedConfiguration string `yaml:"kubectl.kubernetes.io/last-applied-configuration"`
		} `yaml:"annotations"`
		CreationTimestamp time.Time `yaml:"creationTimestamp"`
		Generation        int       `yaml:"generation"`
		Name              string    `yaml:"name"`
		ResourceVersion   string    `yaml:"resourceVersion"`
		SelfLink          string    `yaml:"selfLink"`
		UID               string    `yaml:"uid"`
	} `yaml:"metadata"`
	Spec struct {
		Conversion struct {
			Strategy string `yaml:"strategy"`
		} `yaml:"conversion"`
		Group string `yaml:"group"`
		Names struct {
			Categories []string `yaml:"categories"`
			Kind       string   `yaml:"kind"`
			ListKind   string   `yaml:"listKind"`
			Plural     string   `yaml:"plural"`
			Singular   string   `yaml:"singular"`
		} `yaml:"names"`
		PreserveUnknownFields bool   `yaml:"preserveUnknownFields"`
		Scope                 string `yaml:"scope"`
		Versions              []struct {
			AdditionalPrinterColumns []struct {
				JSONPath string `yaml:"jsonPath"`
				Name     string `yaml:"name"`
				Type     string `yaml:"type"`
			} `yaml:"additionalPrinterColumns"`
			Name   string `yaml:"name"`
			Schema struct {
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
							Description string `yaml:"description"`
							Properties map[string]Property`yaml:"properties"`
							Required []string `yaml:"required"`
							Type     string   `yaml:"type"`
						} `yaml:"spec"`
						Status struct {
							Description string `yaml:"description"`
							Properties  struct {
								AtProvider struct {
									Description string `yaml:"description"`
									Properties  struct {
										IsDefault struct {
											Description string `yaml:"description"`
											Type        string `yaml:"type"`
										} `yaml:"isDefault"`
										OwnerID struct {
											Description string `yaml:"description"`
											Type        string `yaml:"type"`
										} `yaml:"ownerID"`
										Tags struct {
											Description string `yaml:"description"`
											Items       struct {
												Description string `yaml:"description"`
												Properties  struct {
													Key struct {
														Description string `yaml:"description"`
														Type        string `yaml:"type"`
													} `yaml:"key"`
													Value struct {
														Description string `yaml:"description"`
														Type        string `yaml:"type"`
													} `yaml:"value"`
												} `yaml:"properties"`
												Required []string `yaml:"required"`
												Type     string   `yaml:"type"`
											} `yaml:"items"`
											Type string `yaml:"type"`
										} `yaml:"tags"`
										VpcID struct {
											Type string `yaml:"type"`
										} `yaml:"vpcId"`
										VpcState struct {
											Description string   `yaml:"description"`
											Enum        []string `yaml:"enum"`
											Type        string   `yaml:"type"`
										} `yaml:"vpcState"`
									} `yaml:"properties"`
									Type string `yaml:"type"`
								} `yaml:"atProvider"`
								BindingPhase struct {
									Description string   `yaml:"description"`
									Enum        []string `yaml:"enum"`
									Type        string   `yaml:"type"`
								} `yaml:"bindingPhase"`
								Conditions struct {
									Description string `yaml:"description"`
									Items       struct {
										Description string `yaml:"description"`
										Properties  struct {
											LastTransitionTime struct {
												Description string `yaml:"description"`
												Format      string `yaml:"format"`
												Type        string `yaml:"type"`
											} `yaml:"lastTransitionTime"`
											Message struct {
												Description string `yaml:"description"`
												Type        string `yaml:"type"`
											} `yaml:"message"`
											Reason struct {
												Description string `yaml:"description"`
												Type        string `yaml:"type"`
											} `yaml:"reason"`
											Status struct {
												Description string `yaml:"description"`
												Type        string `yaml:"type"`
											} `yaml:"status"`
											Type struct {
												Description string `yaml:"description"`
												Type        string `yaml:"type"`
											} `yaml:"type"`
										} `yaml:"properties"`
										Required []string `yaml:"required"`
										Type     string   `yaml:"type"`
									} `yaml:"items"`
									Type string `yaml:"type"`
								} `yaml:"conditions"`
							} `yaml:"properties"`
							Required []string `yaml:"required"`
							Type     string   `yaml:"type"`
						} `yaml:"status"`
					} `yaml:"properties"`
					Required []string `yaml:"required"`
					Type     string   `yaml:"type"`
				} `yaml:"openAPIV3Schema"`
			} `yaml:"schema"`
			Served       bool `yaml:"served"`
			Storage      bool `yaml:"storage"`
			Subresources struct {
				Status struct {
				} `yaml:"status"`
			} `yaml:"subresources"`
		} `yaml:"versions"`
	} `yaml:"spec"`
	Status struct {
		AcceptedNames struct {
			Categories []string `yaml:"categories"`
			Kind       string   `yaml:"kind"`
			ListKind   string   `yaml:"listKind"`
			Plural     string   `yaml:"plural"`
			Singular   string   `yaml:"singular"`
		} `yaml:"acceptedNames"`
		Conditions []struct {
			LastTransitionTime time.Time `yaml:"lastTransitionTime"`
			Message            string    `yaml:"message"`
			Reason             string    `yaml:"reason"`
			Status             string    `yaml:"status"`
			Type               string    `yaml:"type"`
		} `yaml:"conditions"`
		StoredVersions []string `yaml:"storedVersions"`
	} `yaml:"status"`
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
			APIVersion: o.Spec.Group + "/" + o.Spec.Versions[0].Name,
			Kind:       o.Spec.Names.Kind,
			Spec:       make(map[string]interface{}),
		}

		for k, v := range o.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties.Spec.Properties {
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
