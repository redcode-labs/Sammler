package main
import (

    "github.com/fatih/color"
    "path/filepath"
    "github.com/akamensky/argparse"
    "fmt"
    "os"
    "regexp"
    "io/ioutil"
    "io"
    "strings"
)
var red = color.New(color.FgRed).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()
var cyan = color.New(color.FgBlue).SprintFunc()
var bold = color.New(color.Bold).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var magenta = color.New(color.FgMagenta).SprintFunc()
var blue = color.New(color.FgBlue).SprintFunc()

var found = 0

func read_file(file string) (string, error){
    fil, err := os.Open(file)
    //defer file.Close()
    if err != nil {
        return "", err
    }
    b, err := ioutil.ReadAll(fil)
    if err != nil {
        return "", err
    }
    return string(b), nil
}

func write_file(filename string, data string){
    file, err := os.Create(filename)
    exit_on_error("FILE CREATION ERROR", err)
    defer file.Close()

    _, err = io.WriteString(file, data)
    exit_on_error("FILE WRITE ERROR", err)
}

func exit_on_error(message string, err error){
    if err != nil{
        fmt.Println("%s %s: %s", red("[X]"), bold(message), err.Error())
        os.Exit(0)
    }
}

func print_good(msg string){
    fmt.Printf("%s :: %s \n", green(bold("[+]")), msg)
}

func print_info(msg string){
    fmt.Printf("[*] :: %s\n", msg)
}

func print_error(msg string){
    fmt.Printf("%s :: %s \n", red(bold("[x]")), msg)
}

func data_extract(source string) (map[string][]string, error) {
    src := ""
    if _, err := os.Stat(source); os.IsNotExist(err) {
        src = source
    } else {
        src, err = read_file(source)
        if err != nil {
            return map[string][]string{}, err
        }
        print_info("Loaded file: "+bold(source))
    }
    regexes := map[string]string{
        "mail"   : "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
        "ip"     : `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`,
        "mac"    : `^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`,
        "date"   : `\d{4}-\d{2}-\d{2}`,
        "domain" : `^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`,
        "phone"  : `^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`,
        "ccn"    : `^(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11})$`,
        "time"   : `^([0-9]|0[0-9]|1[0-9]|2[0-3]):([0-9]|[0-5][0-9])$`,
	"crypto" : `^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,39}$`,
    }
    results := map[string][]string{
        "mail"   : []string{},
        "ip"     : []string{},
        "mac"     : []string{},
        "date"   : []string{},
        "domain" : []string{},
        "phone" : []string{},
        "ccn" : []string{},
        "time" : []string{},
	"crypto" : []string{},
    }
    for regex_name, regex := range(regexes){
        r := regexp.MustCompile(regex)
        matches := r.FindAllString(src, -1)
        for _, m := range(matches){
            //b := results[regex_name]
            results[regex_name] = append(results[regex_name], m)
            found += 1
        }
    }
    return results, nil
}

func analyze(source string, remove bool){
    hits, _ := data_extract(source)
    contents, _ := read_file(source)
    if found != 0 {
        print_good(fmt.Sprintf("Found %d interesting strings\n", found))
        for id, matches := range(hits){
            fmt.Println(blue("~~~~~~~~~~~~~~~~~~~~ ["+id+"]"))
            fmt.Println(strings.Join(matches, "\n*"))
        }
    } else{
        print_error("No results found")
    }
    if remove{
        for _, matches := range(hits){
            for _, m := range(matches){
                contents = strings.Replace(contents, m, "", -1)
            }
        }
        write_file(source, contents)
    }
}

func get_all_files() ([]string, error) {
    files := []string{}
    err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        files = append(files, path)
        return nil
    })
    if err != nil{
        return []string{}, err
    } else{
        return files, nil
    }
}

func main(){
    parser := argparse.NewParser("sammler", "")//, usage_prologue)
    //var out *string = parser.String("o", "out", &argparse.Options{Required: false, Default:"default", Help: "Save found strings to a file"})
    var file *[]string = parser.List("f", "file", &argparse.Options{Required: false, Help: "Files/strings to analyze"})
    var all *bool = parser.Flag("a", "all", &argparse.Options{Required: false, Help: "Analyze all files from current directory"})
    var remove *bool = parser.Flag("r", "remove", &argparse.Options{Required: false, Help: "Remove found data from the analyzed file"})
	commandline_args := os.Args
	err := parser.Parse(commandline_args)
    exit_on_error("Parser error", err)

    to_scan := *file
    if *all{
        to_scan, err = get_all_files()
        exit_on_error("Dirwalk error",err)
    }
    for _, f := range(to_scan){
        analyze(f, *remove)
    }
}
