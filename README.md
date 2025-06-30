# My-Process-Manager
A simple process manager for Windows and Linux

<img src="https://github.com/doorbash/my-process-manager/blob/master/screenshot.jpg?raw=true" />

## How to Build

### Front-end
```
npm install --legacy-peer-deps
npm run build
```
### App
**Windows:**
  1. Set CGO_ENABLED=1
  2. Install https://jmeubank.github.io/tdm-gcc/download/
  3. ```wails build -ldflags="-s -w" -s```
  4. You might need to add the My Process Manager folder to Windows Security Exception folders

**Linux:**
  1. Run ```wails doctor``` and install required packages
  2. ```wails build -tags webkit2_41 -ldflags="-s -w" -s```

## License
MIT
