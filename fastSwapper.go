package main

import ("fmt"
        "os" 
        "path/filepath"
        // "log"
        )


func main() {
 
  fmt.Println(getDirsInDir("~/Documents/go_testing/fastSwapper/testFolder"))

  fmt.Println("Hello World!")
}


func getDirsInDir(dir string) [2]string{
  entries, err := os.ReadDir(filepath.FromSlash(dir))
  if err != nil{
    //log.Fatal(err)
    fmt.Println("error occured")
  }
  
  for _, e := range entries{
    fmt.Println(e.Name())
  }
  return [2]string{"a", "z"}
}
