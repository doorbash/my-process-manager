# My-Process-Manager
A simple process manager for Windows

<img src="https://github.com/doorbash/my-process-manager/blob/master/screenshot.jpg?raw=true" />

## Build

### Front End
```
npm install
npm run build
```
### App
```
wails build -tags webkit2_41 -ldflags="-s -w" -race -s -upx -upxflags "-9" -skipbindings
```

## License
MIT
