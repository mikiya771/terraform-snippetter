package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	cs "github.com/mikiya771/terraform-snippetter/internal/create_skelton"
)

func main() {
	os.Exit(func() int {
		if err := gen(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		return 0
	}())
}

func gen() error {
	ctx := context.Background()
	log.Println("ensuring terraform is installed")

	tmpDir, err := ioutil.TempDir("", "tfinstall")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	execPath, err := tfinstall.Find(ctx, tfinstall.LookPath(), tfinstall.LatestVersion(tmpDir, false))
	if err != nil {
		return err
	}

	log.Println("running terraform init")

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	tf, err := tfexec.NewTerraform(cwd, execPath)
	if err != nil {
		return err
	}

	err = tf.Init(ctx, tfexec.Upgrade(true))
	if err != nil {
		return err
	}

	log.Println("creating schemas/data dir")
	err = os.MkdirAll("data", 0755)
	if err != nil {
		return err
	}
	//fs := http.Dir("data")

	// TODO upstream change to have tfexec write to file directly instead of unmarshal/remarshal
	log.Println("running terraform providers schema")
	ps, err := tf.ProvidersSchema(ctx)
	if err != nil {
		return err
	}

	log.Println("creating schemas file")
	schemasFile, err := os.Create(filepath.Join("data", "schemas.json"))
	if err != nil {
		return err
	}
	defer schemasFile.Close()

	log.Println("writing schemas to file")
	err = json.NewEncoder(schemasFile).Encode(ps)
	if err != nil {
		return err
	}

	for k, v := range ps.Schemas {
		rtf, dtf := cs.CreateSnippets(v)
		fs := strings.Split(k, "/")
		fn := fs[len(fs)-1]
		fileName := fmt.Sprintf("tf_%s.snippets", fn)
		snip, err := os.Create(filepath.Join("snippet", fileName))
		if err != nil {
			return err
		}
		for _, st := range rtf {
			for _, l := range st {
				fmt.Fprintln(snip, l)
			}
			fmt.Fprintln(snip, "")
		}
		for _, st := range dtf {
			for _, l := range st {
				fmt.Fprintln(snip, l)
			}
			fmt.Fprintln(snip, "")
		}
		snip.Close()
	}

	return nil
}
