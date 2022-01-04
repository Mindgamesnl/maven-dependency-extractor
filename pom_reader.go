package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Dependency struct {
	Text       string `xml:",chardata"`
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope"`
	Type       string `xml:"type"`
	Exclusions struct {
		Text      string `xml:",chardata"`
		Exclusion struct {
			Text       string `xml:",chardata"`
			GroupId    string `xml:"groupId"`
			ArtifactId string `xml:"artifactId"`
		} `xml:"exclusion"`
	} `xml:"exclusions"`
}

type Pom struct {
	XMLName        xml.Name `xml:"project"`
	Text           string   `xml:",chardata"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	ModelVersion   string   `xml:"modelVersion"`
	GroupId        string   `xml:"groupId"`
	ArtifactId     string   `xml:"artifactId"`
	Version        string   `xml:"version"`
	Packaging      string   `xml:"packaging"`
	Name           string   `xml:"name"`
	Description    string   `xml:"description"`
	Properties     struct {
		Text                       string `xml:",chardata"`
		ProjectBuildSourceEncoding string `xml:"project.build.sourceEncoding"`
	} `xml:"properties"`
	URL          string `xml:"url"`
	Dependencies struct {
		Text       string       `xml:",chardata"`
		Dependency []Dependency `xml:"dependency"`
	} `xml:"dependencies"`
	Build struct {
		Text        string `xml:",chardata"`
		DefaultGoal string `xml:"defaultGoal"`
		Plugins     struct {
			Text   string `xml:",chardata"`
			Plugin []struct {
				Text          string `xml:",chardata"`
				GroupId       string `xml:"groupId"`
				ArtifactId    string `xml:"artifactId"`
				Version       string `xml:"version"`
				Configuration struct {
					Text                        string `xml:",chardata"`
					Source                      string `xml:"source"`
					Target                      string `xml:"target"`
					ClasspathDependencyExcludes struct {
						Text                       string   `xml:",chardata"`
						ClasspathDependencyExclude []string `xml:"classpathDependencyExclude"`
					} `xml:"classpathDependencyExcludes"`
					TrimStackTrace string `xml:"trimStackTrace"`
				} `xml:"configuration"`
				Executions struct {
					Text      string `xml:",chardata"`
					Execution struct {
						Text  string `xml:",chardata"`
						Phase string `xml:"phase"`
						Goals struct {
							Text string `xml:",chardata"`
							Goal string `xml:"goal"`
						} `xml:"goals"`
						Configuration struct {
							Text        string `xml:",chardata"`
							Relocations struct {
								Text       string `xml:",chardata"`
								Relocation []struct {
									Text          string `xml:",chardata"`
									Pattern       string `xml:"pattern"`
									ShadedPattern string `xml:"shadedPattern"`
								} `xml:"relocation"`
							} `xml:"relocations"`
						} `xml:"configuration"`
					} `xml:"execution"`
				} `xml:"executions"`
				Dependencies struct {
					Text       string `xml:",chardata"`
					Dependency []struct {
						Text       string `xml:",chardata"`
						GroupId    string `xml:"groupId"`
						ArtifactId string `xml:"artifactId"`
						Version    string `xml:"version"`
					} `xml:"dependency"`
				} `xml:"dependencies"`
			} `xml:"plugin"`
		} `xml:"plugins"`
		Resources struct {
			Text     string `xml:",chardata"`
			Resource struct {
				Text      string `xml:",chardata"`
				Directory string `xml:"directory"`
				Filtering string `xml:"filtering"`
			} `xml:"resource"`
		} `xml:"resources"`
	} `xml:"build"`
	Profiles struct {
		Text    string `xml:",chardata"`
		Profile struct {
			Text       string `xml:",chardata"`
			ID         string `xml:"id"`
			Activation struct {
				Text     string `xml:",chardata"`
				Property struct {
					Text string `xml:",chardata"`
					Name string `xml:"name"`
				} `xml:"property"`
			} `xml:"activation"`
			Build struct {
				Text    string `xml:",chardata"`
				Plugins struct {
					Text   string `xml:",chardata"`
					Plugin struct {
						Text       string `xml:",chardata"`
						GroupId    string `xml:"groupId"`
						ArtifactId string `xml:"artifactId"`
						Version    string `xml:"version"`
						Executions struct {
							Text      string `xml:",chardata"`
							Execution struct {
								Text  string `xml:",chardata"`
								Phase string `xml:"phase"`
								Goals struct {
									Text string `xml:",chardata"`
									Goal string `xml:"goal"`
								} `xml:"goals"`
								Configuration struct {
									Text             string `xml:",chardata"`
									Executable       string `xml:"executable"`
									WorkingDirectory string `xml:"workingDirectory"`
								} `xml:"configuration"`
							} `xml:"execution"`
						} `xml:"executions"`
					} `xml:"plugin"`
				} `xml:"plugins"`
			} `xml:"build"`
		} `xml:"profile"`
	} `xml:"profiles"`
	Repositories struct {
		Text       string `xml:",chardata"`
		Repository []struct {
			Text string `xml:",chardata"`
			ID   string `xml:"id"`
			URL  string `xml:"url"`
		} `xml:"repository"`
	} `xml:"repositories"`
}

func (d Dependency) localJarPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var jarPath = homeDir + "/.m2/repository/" + strings.ReplaceAll(d.GroupId, ".", "/") + "/" + d.ArtifactId + "/" + d.Version

	// find files
	files, err := ioutil.ReadDir(jarPath)
	if err != nil {
		log.Fatal("Couldn't get content of " + jarPath)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".jar") {
			return jarPath + "/" + file.Name()
		}
	}

	log.Println("WARNING! couldn't find " + d.ArtifactId + " locally!")

	return ""
}

func parsePom(input []byte) Pom {
	var pom Pom
	err := xml.Unmarshal(input, &pom)
	if err != nil {
		fmt.Println("ERROR: That isn't a valid pom file")
		os.Exit(1)
	}

	return pom
}
