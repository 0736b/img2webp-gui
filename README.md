# img2webp-gui

Converted .webp files will be in `output` folder in the same directory of executable.

![test](https://github.com/user-attachments/assets/2a04051f-1508-4ad9-9c56-30c1aa2fb7b3)

### Build
```
go build -a -ldflags="-s -w -H windowsgui -extldflags '-O2'" .
```
