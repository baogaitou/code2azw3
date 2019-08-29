package main

import (
	"bytes"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var dn = flag.String("dir", "", "Please input dir name")
var en = flag.String("name", "", "Please input ebook name")

type page struct {
	Fullpath   string
	FileName   string
	Section    map[int]string
	Structs    map[int]string
	Source     string
	OutputFile string
}

type chapter struct {
	ID          int
	Name        string
	ChapterPath string
}

type ebook struct {
	BookName    string
	BookIssuer  string
	Chapters    []chapter
	Fullpath    string
	NcxFile     string
	OpfFile     string
	TOCFile     string
	WelcomeFile string
}

type file struct {
	Name string
}

func (f *file) trimSlashes() {
	if string(f.Name[len(f.Name)-1]) == "/" {
		f.Name = f.Name[:len(f.Name)-1]
	}
}

func main() {
	flag.Parse()

	eb := ebook{}
	eb.BookName = string(*en)
	eb.BookIssuer = "Amazon.cn"
	eb.NcxFile = eb.BookName + "/" + eb.BookName + ".ncx"
	eb.OpfFile = eb.BookName + "/" + eb.BookName + ".opf"
	eb.TOCFile = eb.BookName + "/toc.html"
	eb.WelcomeFile = eb.BookName + "/welcome.html"

	sourceDir := file{Name: string(*dn)}
	sourceDir.trimSlashes()

	var s []string
	s, _ = getAllFile(sourceDir.Name, s)

	eb.Chapters = make([]chapter, len(s))

	template := loadTemplate("templates/ebook.html")
	ncxTemplate := loadTemplate("templates/ebook.ncx")
	opfTemplate := loadTemplate("templates/ebook.opf")
	tocTemplate := loadTemplate("templates/ebook.toc.html")
	welcomeTemplate := loadTemplate("templates/ebook.welcome.html")

	for i, f := range s {
		sourceFile := file{Name: f}

		chptr := chapter{}

		outBuff := new(bytes.Buffer)
		ncxBuff := new(bytes.Buffer)
		opfBuff := new(bytes.Buffer)
		tocBuff := new(bytes.Buffer)
		welcomeBuff := new(bytes.Buffer)

		data := make(map[string]interface{})
		data["Title"] = eb.BookName

		rt := sourceFile.Parse()

		rt.FileName = strings.Replace(rt.Fullpath, sourceDir.Name+"/", "", -1)

		hFileName := packageFileName(rt.FileName)
		rt.OutputFile = fmt.Sprintf("%s/char-%d-%s.html", eb.BookName, i, hFileName)

		// 创建文件夹
		checkPath(eb.BookName)

		data["Filename"] = rt.FileName
		data["Body"] = rt.Source
		data["Sections"] = rt.Section
		data["Structs"] = rt.Structs

		chptr.ID = i
		chptr.Name = rt.FileName
		chptr.ChapterPath = fmt.Sprintf("char-%d-%s.html", i, hFileName)

		eb.Chapters[i] = chptr

		ncxTemplate.Execute(ncxBuff, eb)
		eb.Save(ncxBuff, eb.NcxFile)

		opfTemplate.Execute(opfBuff, eb)
		eb.Save(opfBuff, eb.OpfFile)

		tocTemplate.Execute(tocBuff, eb)
		eb.Save(tocBuff, eb.TOCFile)

		welcomeTemplate.Execute(welcomeBuff, eb)
		eb.Save(welcomeBuff, eb.WelcomeFile)

		template.Execute(outBuff, data)
		rt.Save(outBuff)
	}

	// // 生成封面图
	// bgColor := color.NRGBA{0, 0, 0, 255}
	// dist := imaging.New(568, 800, bgColor)

	// fontSource := "fonts/FiraCode-Regular.ttf"

	// fontSourceBytes, err := ioutil.ReadFile(fontSource)
	// if err != nil {
	// 	fmt.Println("ioutil.ReadFile(fontSource) failed.")
	// }

	// trueTypeFont, err := freetype.ParseFont(fontSourceBytes)
	// if err != nil {
	// 	fmt.Println("freetype.ParseFont(fontSourceBytes) failed.")
	// }

	// fontColor := color.NRGBA{255, 255, 255, 255}
	// fg := image.NewUniform(fontColor)

	// fc := freetype.NewContext()
	// fc.SetDPI(72 * 2)
	// fc.SetFont(trueTypeFont)
	// fc.SetFontSize(float64(12 * 2))
	// fc.SetClip(dist.Bounds())
	// fc.SetDst(dist)
	// fc.SetSrc(fg)

	// pt := freetype.Pt(10, 200)
	// _, err = fc.DrawString(eb.BookName, pt)
	// if err != nil {
	// 	log.Println("fc.DrawString failed.")
	// }
	// dist = imaging.Resize(dist, 560, 800, imaging.Lanczos)
	// imaging.Save(dist, "cover.png")

	// sourceFile := file{Name: string(*fn)}
	// if sourceFile.Name != "" {
	// 	outBuff := new(bytes.Buffer)
	// 	template := loadTemplate("template.html")

	// 	data := make(map[string]string)
	// 	data["Title"] = "Code"

	// 	rt := sourceFile.Parse()
	// 	rt.OutputFile = "output.html"

	// 	data["Body"] = rt.Source

	// 	template.Execute(outBuff, data)

	// 	rt.Save(outBuff)
	// }
}

func getAllFile(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if string(fi.Name()[0]) != "." {
			if fi.IsDir() {
				fullDir := pathname + "/" + fi.Name()
				s, err = getAllFile(fullDir, s)
				if err != nil {
					fmt.Println("read dir fail:", err)
					return s, err
				}
			} else {
				fullName := pathname + "/" + fi.Name()
				if getFileContentType(fullName) == "text/plain; charset=utf-8" {
					s = append(s, fullName)
				}
			}
		}
	}
	return s, nil
}

func getFileContentType(fn string) (t string) {
	out, _ := os.Open(fn)
	defer out.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return ""
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType
}

func (f file) Parse() (oPage page) {
	oPage.Fullpath = f.Name
	sourceFile, err := os.Open(f.Name)
	defer sourceFile.Close()
	if err != nil {
		fmt.Println(err)
	}

	rr, err := ioutil.ReadAll(sourceFile)
	oPage.Source = string(rr)

	if err != nil {
		fmt.Println("ioutil.ReadALl Err")
	}

	// Get func list
	reg := regexp.MustCompile(`(?m)(^func .*)$`)
	funcs := reg.FindAllString(oPage.Source, -1)

	oPage.Section = make(map[int]string, len(funcs))

	for i, fc := range funcs {
		fc = strings.Replace(fc, "{", "", -1)
		fc = strings.Replace(fc, "func ", "", -1)
		oPage.Section[i] = fc
	}

	// Get type list
	reg2 := regexp.MustCompile(`(?m)(^type .*)$`)
	structs := reg2.FindAllString(oPage.Source, -1)

	oPage.Structs = make(map[int]string, len(structs))

	for i, fc := range structs {
		fc = strings.Replace(fc, "type ", "", -1)
		fc = strings.Replace(fc, "struct {", "", -1)
		oPage.Structs[i] = fc
	}

	oPage.Source = html.EscapeString(oPage.Source)
	oPage.Source = strings.Replace(oPage.Source, "  ", "&nbsp;&nbsp;", -1)
	oPage.Source = strings.Replace(oPage.Source, "	", "&nbsp;&nbsp;&nbsp;&nbsp;", -1)
	oPage.Source = strings.Replace(oPage.Source, "\n\n", "<br /><br />", -1)
	oPage.Source = strings.Replace(oPage.Source, "\n", "<br />", -1)

	return
}

func (p page) Save(o *bytes.Buffer) {
	f2, err := os.OpenFile(p.OutputFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	defer f2.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f2.Write([]byte(o.String()))
	}
}

func (p ebook) Save(o *bytes.Buffer, outputFile string) {
	f2, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	defer f2.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f2.Write([]byte(o.String()))
	}
}

func loadTemplate(tName string) *template.Template {
	templateFile, _ := os.Open(tName)
	defer templateFile.Close()

	templateByte, _ := ioutil.ReadAll(templateFile)
	templateString := string(templateByte)

	// Parse
	t, _ := template.New("Page").Parse(templateString)

	return t
}

func packageFileName(fullpath string) (fn string) {
	r := strings.Split(fullpath, "/")
	for _, a := range r {
		fn = fn + "_" + a
	}
	fn = strings.Replace(fn[1:len(fn)], ".", "_", -1)
	return strings.ToLower(fn)
}

func checkPath(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		}
	}
}
