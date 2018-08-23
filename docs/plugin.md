## TODO
- [ ] package manager
- [ ] update packages

### Plugin architecture
I don't know why but vg fix the plugin issue
https://github.com/GetStream/vg
https://github.com/golang/go/issues/20481
https://github.com/hashicorp/go-plugin
https://stackoverflow.com/questions/42388090/go-1-8-plugin-use-custom-interface
https://stackoverflow.com/questions/42218472/how-do-go-plugin-dependencies-work/42220856#42220856


## package manager
```sh
# install and build plugin
tada install package-name

# link plugin
tada link foo-bar.so

```
* read installed packages from a storage (file or kv store?)

### how to install packages?
install .so or install .go and compile them?

download source code from github and compile them

remove that plugin repo before install or maybe we can pull(fetch + merge) that repo

handle failed
