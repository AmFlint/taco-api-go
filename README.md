# Golang REST Api for Taco Application (Trello Like)

## How to

### Install

First, you will need to install (Golang)[https://golang.org/doc/install]. Then, you will have to create the directory to clone this project into.
***Note*** that you have to clone this repository inside your $GOPATH/src/ folder (By default on mac/linux -> ~/go/src)
```bash
# Create right directory to clone this project, replace $GOPATH with ~/go if you followed default installation
mkdir -p $GOPATH/src/github.com/AmFlint
# Then, download it, replace $GOPATH
git clone https://github.com/AmFlint/taco-api-go.git $GOPATH/src/github.com/AmFlint/taco-api-go
```

Once the project is installed, please install (Glide)[https://github.com/Masterminds/glide] (one of many Golang dependency manager).
Then, run the following commands:
```bash
# Install dependencies
glide install
```

### Run the application

In order to run the application, you have multiple choices:

- For development purposes, you will probably want to use

```bash
go run main.go
```

- For deployment purposes, use the command:
```bash
# Build the application -> Binaries to execute under the file 'taco'
go build main.go -o taco

# And then, just run the binaries
./taco
```

### Write tests
