# Add this at the end of your .bashrc file
gvm use go1.22 > /dev/null 2>&1 || true


# Make sure Go bin directory is always in the PATH
export PATH=$PATH:$(go env GOPATH)/bin 