package config

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gandrille/go-commons/env"
	"github.com/gandrille/go-commons/misc"
	"github.com/gandrille/go-commons/result"
)

type configurationBuilder interface {
	kind() string
	build(e entry) (*Configuration, error)
}

func builders() []configurationBuilder {
	var builders []configurationBuilder

	builders = append(builders, unisonBuilder{})
	builders = append(builders, programmeBuilder{})
	builders = append(builders, rsyncBuilder{})

	return builders
}

/* ======================== */

type unisonBuilder struct{}

func (b unisonBuilder) kind() string {
	return "unison"
}

func (b unisonBuilder) build(e entry) (*Configuration, error) {
	if b.kind() != e.kind {
		return nil, fmt.Errorf("Can't build %s because entry has type %s", b.kind(), e.kind)
	}
	mandatory := []string{"profile"}
	optional := []string{}

	if err := e.checkParams(mandatory, optional); err != nil {
		return nil, err
	}

	profile := e.getParameterValue("profile")
	cmd := exec.Command("/usr/bin/unison", profile)
	runner := func() result.Result { return misc.RunCmd(cmd, e.key) }

	return &Configuration{e.key, runner}, nil
}

/* =========================== */

type programmeBuilder struct{}

func (b programmeBuilder) kind() string {
	return "prog"
}

func (b programmeBuilder) build(e entry) (*Configuration, error) {
	if b.kind() != e.kind {
		return nil, fmt.Errorf("Can't build %s because entry has type %s", b.kind(), e.kind)
	}
	mandatory := []string{"exe"}
	optional := []string{"arg"}

	if err := e.checkParams(mandatory, optional); err != nil {
		return nil, err
	}

	exe := e.getParameterValue("exe")
	args := e.getParameterSlice("arg")

	cmd := exec.Command(exe, args...)
	runner := func() result.Result { return misc.RunCmd(cmd, e.key) }

	return &Configuration{e.key, runner}, nil
}

/* =========================== */

type rsyncBuilder struct{}

func (b rsyncBuilder) kind() string {
	return "rsync"
}

func (b rsyncBuilder) build(e entry) (*Configuration, error) {
	if b.kind() != e.kind {
		return nil, fmt.Errorf("Can't build %s because entry has type %s", b.kind(), e.kind)
	}
	mandatory := []string{"src", "dst"}
	optional := []string{"host", "cmp", "flag", "exclude", "srcPrefix"}

	if err := e.checkParams(mandatory, optional); err != nil {
		return nil, err
	}

	host := e.getParameterValue("host")
	cmp := e.getParameterSlice("cmp")

	src := e.getParameterSlice("src")
	dst := e.getParameterValue("dst")
	flag := e.getParameterSlice("flag")
	exclude := e.getParameterSlice("exclude")
	srcPrefix := e.getParameterValue("srcPrefix")

	cmd, errCmd := rsyncCmd(host, cmp, src, dst, flag, exclude, srcPrefix)
	if errCmd != nil {
		return nil, errCmd
	}

	runner := func() result.Result { return misc.RunCmd(cmd, e.key) }
	return &Configuration{e.key, runner}, nil
}

func rsyncCmd(host string, cmp, src []string, dst string, flag, exclude []string, srcPrefix string) (*exec.Cmd, error) {
	// parameters checking
	if host != "" && host != env.Hostname() {
		return nil, errors.New("This rsync execution must be executed on " + host)
	}
	for _, mpoint := range cmp {
		mounted, err := misc.IsMounted(mpoint)
		if err != nil {
			return nil, errors.New("Error while checking mount point " + mpoint + " state: " + err.Error())
		}
		if !mounted {
			return nil, errors.New("Mount point " + mpoint + " is not mounted")
		}
	}
	if len(src) == 0 {
		return nil, errors.New("'src' must be provided at least once")
	}
	if dst == "" {
		return nil, errors.New("'dst' flag can't be empty")
	}

	// arguments computation
	args := flag
	for _, ex := range exclude {
		args = append(args, "--exclude", ex)
	}
	sources := buildSources(src, srcPrefix)
	args = append(args, sources...)
	args = append(args, dst)

	// building command (it is NOT executed !)
	cmd := exec.Command("/usr/bin/rsync", args...)

	return cmd, nil
}

func buildSources(src []string, srcPrefix string) []string {
	var srcList []string
	if srcPrefix == "" {
		srcList = src
	} else {
		prefix := strings.TrimSuffix(srcPrefix, "/")
		for _, s := range src {
			cleanSrc := strings.TrimPrefix(s, "/")
			srcList = append(srcList, prefix+"/"+cleanSrc)
		}
	}
	return srcList
}
