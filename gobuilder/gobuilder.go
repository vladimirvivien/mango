package gobuilder

import (
	"os"
	"runtime"

	"github.com/constabulary/gb"
)

type Builder struct {
	SrcRoot  string
	Packages []string

	ShouldForceBuild bool
	ShouldInstall    bool
	ShouldBeVerbose  bool
	ShouldDetectRace bool

	GOOS   string
	GOARCH string

	ParallelBuilds int

	BuildTags []string
	GCFlags   []string
	LDFlags   []string

	ctx gb.Context
}

func New() *Builder {
	return &Builder{
		SrcRoot:          os.Getenv("GOPATH"),
		ShouldInstall:    false,
		ShouldBeVerbose:  false,
		ShouldDetectRace: false,

		GOOS:   runtime.GOOS,
		GOARCH: runtime.GOARCH,

		ParallelBuilds: runtime.NumCPU(),
	}
}

func (b *Builder) Build() {
	project := gb.NewProject(b.SrcRoot)
	ctx, err := gb.NewContext(
		project,
		gb.GcToolchain(),
		gb.Gcflags(b.GCFlags...),
		gb.Ldflags(b.LDFlags...),
		gb.Tags(b.BuildTags...),
		func(c *gb.Context) error {
			// TODO assert supported platforms for race detector
			if b.ShouldDetectRace {
				return gb.WithRace(c)
			}
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	var pkgs []*gb.Package
	for _, pkgPath := range b.Packages {
		pkg, err := ctx.ResolvePackage(pkgPath)
		if err != nil {
			panic(err)
		}
		pkgs := append(pkgs, pkg)
		action, err := gb.BuildPackages(pkgs...)
		if err := gb.ExecuteConcurrent(action, b.ParallelBuilds, nil); err != nil {
			panic(err)
		}
	}
}
