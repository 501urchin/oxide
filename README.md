# taurine
taurine is a cli tool to easily provision and deploy to a vps 

# commands
- taurine init - initializes a project by creating a taurine.toml file
- taurine deploy - deploys to a vps
- taurine dev - demo app locally
- taurine rollback [version] - rollsback to previous version
- taurine logs - get server logs
- taurine health - performs healthcheck



# notes
- if user chooses exe as deploy method - build exe locally and scp to vps and make systemd file
- for docker, build image loalaly to tar file, send compressd file to vps via scp and laod it there