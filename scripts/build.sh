# Build PKGs
mkdir -p ./pkg

GOOS=darwin GOARCH=amd64 go build -o ./pkg/wrapper-darwin-amd64 ./cmd/wrapper/main.go
GOOS=darwin GOARCH=arm64 go build -o ./pkg/wrapper-darwin-arm64 ./cmd/wrapper/main.go
GOOS=linux GOARCH=386 go build -o ./pkg/wrapper-linux-386 ./cmd/wrapper/main.go
GOOS=linux GOARCH=amd64 go build -o ./pkg/wrapper-linux-amd64 ./cmd/wrapper/main.go
GOOS=linux GOARCH=arm go build -o ./pkg/wrapper-linux-arm ./cmd/wrapper/main.go
GOOS=linux GOARCH=arm64 go build -o ./pkg/wrapper-linux-arm64 ./cmd/wrapper/main.go
GOOS=windows GOARCH=386 go build -o ./pkg/wrapper-windows-386.exe ./cmd/wrapper/main.go
GOOS=windows GOARCH=amd64 go build -o ./pkg/wrapper-windows-amd64.exe ./cmd/wrapper/main.go
GOOS=windows GOARCH=arm go build -o ./pkg/wrapper-windows-arm.exe ./cmd/wrapper/main.go
GOOS=windows GOARCH=arm64 go build -o ./pkg/wrapper-windows-arm64.exe ./cmd/wrapper/main.go

# Build Apps

# Windows
mkdir -p ./build/windows/386
mkdir -p ./build/windows/amd64
mkdir -p ./build/windows/arm
mkdir -p ./build/windows/arm64

GOOS=windows GOARCH=386 go build -o ./build/windows/386/exwrap.exe ./cmd/exwrap/
GOOS=windows GOARCH=amd64 go build -o ./build/windows/amd64/exwrap.exe ./cmd/exwrap/
GOOS=windows GOARCH=arm go build -o ./build/windows/arm/exwrap.exe ./cmd/exwrap/
GOOS=windows GOARCH=arm64 go build -o ./build/windows/arm64/exwrap.exe ./cmd/exwrap/

# Linux
mkdir -p ./build/linux/386
mkdir -p ./build/linux/amd64
mkdir -p ./build/linux/arm
mkdir -p ./build/linux/arm64

GOOS=linux GOARCH=386 go build -o ./build/linux/386/exwrap ./cmd/exwrap/
GOOS=linux GOARCH=amd64 go build -o ./build/linux/amd64/exwrap ./cmd/exwrap/
GOOS=linux GOARCH=arm go build -o ./build/linux/arm/exwrap ./cmd/exwrap/
GOOS=linux GOARCH=arm64 go build -o ./build/linux/arm64/exwrap ./cmd/exwrap/

# OSX
mkdir -p ./build/darwin/amd64
mkdir -p ./build/darwin/arm64

GOOS=darwin GOARCH=amd64 go build -o ./build/darwin/amd64/exwrap ./cmd/exwrap/
GOOS=darwin GOARCH=arm64 go build -o ./build/darwin/arm64/exwrap ./cmd/exwrap/

# Copy the PKGs

cp -r ./pkg ./build/windows/386/pkg
cp -r ./pkg ./build/windows/amd64/pkg
cp -r ./pkg ./build/windows/arm/pkg
cp -r ./pkg ./build/windows/arm64/pkg
cp -r ./pkg ./build/linux/386/pkg
cp -r ./pkg ./build/linux/amd64/pkg
cp -r ./pkg ./build/linux/arm/pkg
cp -r ./pkg ./build/linux/arm64/pkg
cp -r ./pkg ./build/darwin/amd64/pkg
cp -r ./pkg ./build/darwin/arm64/pkg

# Fix executables
find ./build -name "*-386" -exec chmod +x {} \;
find ./build -name "*-amd64" -exec chmod +x {} \;
find ./build -name "*-arm" -exec chmod +x {} \;
find ./build -name "*-arm64" -exec chmod +x {} \;
find ./build -name "*exwrap" -exec chmod +x {} \;

