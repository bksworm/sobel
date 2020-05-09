# Sobel Operator Edge Detection

This is a command line utility which implements edge detection using the [sobel filter](https://en.wikipedia.org/wiki/Sobel_operator). 

# Installation 
Given that this is a go project, installation is simple. 
The only tool you'll need to run this implementation is [go](https://golang.org). 

If go is installed, run: 
```
 $ git clone https://github.com/bksworm/sobel && cd sobel
 $ go build 
``` 
# Usage 
Usage is relatively simple, there are just a few options.
This tool supports both png and jpeg image input formats, and file type is determined automatically. 

To run detection on a **specific file**, use: 

`$ ./examples/edgedetect -f <file.[png, jpg]>`

You can optionally specify an **output file**:

`$ ./examples/edgedetect -f <file.[png, jpg]> -o <output.[png, jpg]>`

Default output is `sobel.jpg` or `sobel.png`

*Note: this is not a conversion program. 
The output file will be of the same format as the inut file, so name files accordingly.* 

# Usage as a go library 
The internal package `sobel` can be used in any standard go program. It exposes a single function.
Here is an example:
```go 
package main 

import (
  "image"
  "os"
  "github.com/bksworm/sobel" //package which implements the filter
 )
 
 var edge image.Image
 
 func main() {
    f, err := os.Open("example.jpg")
    if err != nil { panic(err) }
    defer f.Close()
    
    img, _, err := image.Decode(f)
    if err != nil { panic(err) }
    
    edge = sobel.FilterMath(img) //converts "img" to grayscale and runs edge detect. Returns an image.Image with changes.
    //do something with detected image...
 }
```
Function FilterMath() is pure go filter implementation. If you like things to go 8-9 times faster you may use FilterSimd(). This one is based on [Simd library](https://ermig1979.github.io/Simd/help/group__sobel__filter.html#gace953da81ab3f334ec6435d92ac52c05).  But if you don't need it, you may clear the mess I have here :)

There are a few another implementations of the filter. You can use them in benchmark tests to get an idea about go code performance and memory management tricks.

Happy coding for everyone!