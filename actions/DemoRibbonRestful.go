package actions

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chunhui2001/go-starter/core/utils"
	"github.com/gin-gonic/gin"

	_ "github.com/chunhui2001/go-starter/core/commons"
	. "github.com/fogleman/fauxgl"
	"github.com/fogleman/ribbon/pdb"
	"github.com/fogleman/ribbon/ribbon"
	"github.com/nfnt/resize"
)

const (
	size  = 2048
	scale = 4
)

func downloadAndParse(structureID string) ([]*pdb.Model, error) {
	url := fmt.Sprintf(
		"https://files.rcsb.org/download/%s.pdb.gz",
		strings.ToUpper(structureID))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return pdb.NewReader(r).ReadAll()
}

func timed(name string) func() {
	if len(name) > 0 {
		fmt.Printf("%s ... ", name)
	}
	start := time.Now()
	return func() {
		fmt.Println(time.Since(start))
	}
}

func RibbonDiagramsRouter(ctx *gin.Context) {

	structureID := ctx.Query("channel")

	if structureID == "" {
		structureID = "4hhb"
	}

	var done func()

	done = timed("downloading pdb file")
	models, err := downloadAndParse(structureID)
	if err != nil {
		log.Fatal(err)
	}
	model := models[0]
	done()

	fmt.Printf("atoms       = %d\n", len(model.Atoms))
	fmt.Printf("residues    = %d\n", len(model.Residues))
	fmt.Printf("chains      = %d\n", len(model.Chains))
	fmt.Printf("helixes     = %d\n", len(model.Helixes))
	fmt.Printf("strands     = %d\n", len(model.Strands))
	fmt.Printf("het-atoms   = %d\n", len(model.HetAtoms))
	fmt.Printf("connections = %d\n", len(model.Connections))

	done = timed("generating triangle mesh")
	mesh := ribbon.ModelMesh(model)
	done()

	fmt.Printf("triangles   = %d\n", len(mesh.Triangles))

	done = timed("transforming mesh")
	m := mesh.BiUnitCube()
	done()

	done = timed("finding ideal camera position")
	c := ribbon.PositionCamera(model, m)
	done()

	if tmpFileStl, err := ioutil.TempFile(os.TempDir(), "myname.*"+structureID+".stl"); err != nil {
		panic(err)
	} else {
		// defer os.Remove(tmpFileStl.Name())
		done = timed("writing mesh to disk, " + tmpFileStl.Name())
		mesh.SaveSTL(tmpFileStl.Name())
		done()
	}

	// render
	done = timed("rendering image")
	context := NewContext(int(size*scale*c.Aspect), size*scale)
	context.ClearColorBufferWith(HexColor("1D181F"))
	matrix := LookAt(c.Eye, c.Center, c.Up).Perspective(c.Fovy, c.Aspect, 1, 100)
	light := c.Eye.Sub(c.Center).Normalize()
	shader := NewPhongShader(matrix, light, c.Eye)
	shader.AmbientColor = Gray(0.3)
	shader.DiffuseColor = Gray(0.9)
	context.Shader = shader
	context.DrawTriangles(mesh.Triangles)
	done()

	// save image
	done = timed("downsampling image")
	image := context.Image()
	image = resize.Resize(uint(size*c.Aspect), size, image, resize.Bilinear)
	done()

	if tmpFilePng, err := ioutil.TempFile(os.TempDir(), "myname.*"+structureID+".png"); err != nil {
		panic(err)
	} else {

		defer os.Remove(tmpFilePng.Name())

		done = timed("writing image to disk, " + tmpFilePng.Name())
		SavePNG(tmpFilePng.Name(), image)
		done()

		if b, err := utils.ReadFile(tmpFilePng.Name()); err != nil {
			panic(err)
		} else {

			ctx.Header("Content-Type", "image/png")
			// Browser download or preview
			ctx.Header("Content-Disposition", "inline")
			ctx.Header("Content-Transfer-Encoding", "binary")
			ctx.Header("Cache-Control", "no-cache")
			ctx.Header("Content-Length", fmt.Sprintf("%d", len(b)))

			ctx.Writer.Write(b)
		}

	}
}
