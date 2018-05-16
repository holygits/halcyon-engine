
I often find "getting started" sections are missing or gloss over important details for first timers which are taken for granted by project maintainers.

Here is a guide to using containerd and developing a simple daemon.

## Clone & Build
```
go get https://github.com/containerd/containerd
cd $GOPATH/src/github.com/containerd/containerd
make
sudo make install
```

## Install systemd service
```
cp containerd.service /etc/systemd/system/
sudo systemctl enable containerd.service
sudo systemctl start containerd.service
```

# Install RunC
sudo pacman -S runc
