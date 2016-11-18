package main

import(
   "fmt"
   "golang.org/x/sys/windows/registry"
)

func main(){
   disable_execution()
}

func disable_execution(){
   k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\PowerShell\1\ShellIds\Microsoft.PowerShell`, registry.SET_VALUE)
   if err!= nil{
      fmt.Println(err)
   }
   defer k.Close()
   err = k.SetStringValue("ExecutionPolicy", "restricted")
   if err!= nil{
      fmt.Println(err)
   }
}

